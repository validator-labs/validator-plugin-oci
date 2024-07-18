/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package controller defines a controller for reconciling OciValidator objects.
package controller

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ktypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	vapi "github.com/validator-labs/validator/api/v1alpha1"
	"github.com/validator-labs/validator/pkg/types"
	"github.com/validator-labs/validator/pkg/util"
	vres "github.com/validator-labs/validator/pkg/validationresult"

	"github.com/validator-labs/validator-plugin-oci/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-oci/internal/constants"
	val "github.com/validator-labs/validator-plugin-oci/internal/validators"
	"github.com/validator-labs/validator-plugin-oci/pkg/oci"
)

// OciValidatorReconciler reconciles a OciValidator object
type OciValidatorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=validation.spectrocloud.labs,resources=ocivalidators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=validation.spectrocloud.labs,resources=ocivalidators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=validation.spectrocloud.labs,resources=ocivalidators/finalizers,verbs=update

// Reconcile reconciles each rule found in each OCIValidator in the cluster and creates ValidationResults accordingly
func (r *OciValidatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := r.Log.V(0).WithValues("name", req.Name, "namespace", req.Namespace)

	l.Info("Reconciling OciValidator")

	validator := &v1alpha1.OciValidator{}
	if err := r.Get(ctx, req.NamespacedName, validator); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get the active validator's validation result
	vr := &vapi.ValidationResult{}
	p, err := patch.NewHelper(vr, r.Client)
	if err != nil {
		l.Error(err, "failed to create patch helper")
		return ctrl.Result{}, err
	}
	nn := ktypes.NamespacedName{
		Name:      validationResultName(validator),
		Namespace: req.Namespace,
	}
	if err := r.Get(ctx, nn, vr); err == nil {
		vres.HandleExistingValidationResult(vr, r.Log)
	} else {
		if !apierrs.IsNotFound(err) {
			l.Error(err, "unexpected error getting ValidationResult")
		}
		if err := vres.HandleNewValidationResult(ctx, r.Client, p, buildValidationResult(validator), r.Log); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: time.Millisecond}, nil
	}

	// Always update the expected result count in case the validator's rules have changed
	vr.Spec.ExpectedResults = validator.Spec.ResultCount()

	resp := types.ValidationResponse{
		ValidationRuleResults: make([]*types.ValidationRuleResult, 0, vr.Spec.ExpectedResults),
		ValidationRuleErrors:  make([]error, 0, vr.Spec.ExpectedResults),
	}

	// OCI Registry rules
	for _, rule := range validator.Spec.OciRegistryRules {
		username, password, err := r.secretKeyAuth(req, rule)
		if err != nil {
			l.Error(err, "failed to get auth for rule: %s", rule.Name)
			resp.AddResult(nil, err)
			continue
		}

		pubKeys, err := r.signaturePubKeys(req, rule)
		if err != nil {
			l.Error(err, "failed to get signature verification public keys for rule: %s", rule.Name)
			resp.AddResult(nil, err)
			continue
		}

		ociClient, err := oci.NewOCIClient(
			oci.WithBasicAuth(username, password),
			oci.WithMultiAuth(),
			oci.WithTLSConfig(rule.InsecureSkipTLSVerify, rule.CaCert, ""),
			oci.WithVerificationPublicKeys(pubKeys),
		)
		if err != nil {
			l.Error(err, "failed to create OCI client for rule: %s", rule.Name)
			resp.AddResult(nil, err)
			continue
		}

		ociRuleService := val.NewOciRuleService(r.Log, val.WithOCIClient(ociClient))

		vrr, err := ociRuleService.ReconcileOciRegistryRule(rule)
		if err != nil {
			l.Error(err, "failed to reconcile OCI Registry rule")
		}
		resp.AddResult(vrr, err)
	}

	if err := vres.SafeUpdateValidationResult(ctx, p, vr, resp, r.Log); err != nil {
		return ctrl.Result{}, err
	}

	l.Info("Requeuing for re-validation in two minutes.")
	return ctrl.Result{RequeueAfter: time.Second * 120}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OciValidatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.OciValidator{}).
		Complete(r)
}

// secretKeyAuth retrieves the username and password from the secret referenced in the rule's auth field.
// Any additional key-value pairs in the secret are set as environment variables, to be picked up by auth keychains (e.g. ECR, Azure).
func (r *OciValidatorReconciler) secretKeyAuth(req ctrl.Request, rule v1alpha1.OciRegistryRule) (string, string, error) {
	if rule.Auth.SecretName == "" {
		return "", "", nil
	}

	authSecret := &corev1.Secret{}
	nn := ktypes.NamespacedName{Name: rule.Auth.SecretName, Namespace: req.Namespace}

	if err := r.Get(context.Background(), nn, authSecret); err != nil {
		return "", "", fmt.Errorf("failed to fetch auth secret %s/%s for rule %s: %w", nn.Name, nn.Namespace, rule.Name(), err)
	}

	if len(authSecret.Data) == 0 {
		return "", "", fmt.Errorf("secret %s/%s has no data", nn.Name, nn.Namespace)
	}

	var username, password string
	for k, v := range authSecret.Data {
		if k == "username" {
			username = string(v)
			continue
		}
		if k == "password" {
			password = string(v)
			continue
		}
		if err := os.Setenv(k, string(v)); err != nil {
			return username, password, err
		}
		r.Log.Info("Set environment variable", "key", k)
	}

	return username, password, nil
}

func (r *OciValidatorReconciler) signaturePubKeys(req ctrl.Request, rule v1alpha1.OciRegistryRule) ([][]byte, error) {
	if rule.SignatureVerification.SecretName == "" {
		return nil, nil
	}

	pubKeysSecret := &corev1.Secret{}
	nn := ktypes.NamespacedName{Name: rule.SignatureVerification.SecretName, Namespace: req.Namespace}

	if err := r.Get(context.Background(), nn, pubKeysSecret); err != nil {
		return nil, fmt.Errorf("failed to fetch public keys secret %s/%s for rule %s: %w",
			nn.Name, nn.Namespace, rule.Name(), err,
		)
	}

	pubKeys := make([][]byte, 0)
	for k, data := range pubKeysSecret.Data {
		if strings.HasSuffix(k, ".pub") {
			pubKeys = append(pubKeys, data)
		}
	}
	if len(pubKeys) == 0 {
		return nil, fmt.Errorf("no public keys found in secret %s/%s for rule: %s", nn.Name, nn.Namespace, rule.Name())
	}

	return pubKeys, nil
}

func buildValidationResult(validator *v1alpha1.OciValidator) *vapi.ValidationResult {
	return &vapi.ValidationResult{
		ObjectMeta: metav1.ObjectMeta{
			Name:      validationResultName(validator),
			Namespace: validator.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: validator.APIVersion,
					Kind:       validator.Kind,
					Name:       validator.Name,
					UID:        validator.UID,
					Controller: util.Ptr(true),
				},
			},
		},
		Spec: vapi.ValidationResultSpec{
			Plugin:          constants.PluginCode,
			ExpectedResults: validator.Spec.ResultCount(),
		},
	}
}

func validationResultName(validator *v1alpha1.OciValidator) string {
	return fmt.Sprintf("validator-plugin-oci-%s", validator.Name)
}
