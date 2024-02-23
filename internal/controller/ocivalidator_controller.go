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

package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ktypes "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/spectrocloud-labs/validator-plugin-oci/api/v1alpha1"
	"github.com/spectrocloud-labs/validator-plugin-oci/internal/constants"
	val "github.com/spectrocloud-labs/validator-plugin-oci/internal/validators"
	vapi "github.com/spectrocloud-labs/validator/api/v1alpha1"
	"github.com/spectrocloud-labs/validator/pkg/util/ptr"
	vres "github.com/spectrocloud-labs/validator/pkg/validationresult"
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
	r.Log.V(0).Info("Reconciling OciValidator", "name", req.Name, "namespace", req.Namespace)

	validator := &v1alpha1.OciValidator{}
	if err := r.Get(ctx, req.NamespacedName, validator); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get the active validator's validation result
	vr := &vapi.ValidationResult{}
	nn := ktypes.NamespacedName{
		Name:      validationResultName(validator),
		Namespace: req.Namespace,
	}
	if err := r.Get(ctx, nn, vr); err == nil {
		vres.HandleExistingValidationResult(nn, vr, r.Log)
	} else {
		if !apierrs.IsNotFound(err) {
			r.Log.V(0).Error(err, "unexpected error getting ValidationResult", "name", nn.Name, "namespace", nn.Namespace)
		}
		if err := vres.HandleNewValidationResult(r.Client, buildValidationResult(validator), r.Log); err != nil {
			return ctrl.Result{}, err
		}
	}

	ociRuleService := val.NewOciRuleService(r.Log)

	// OCI Registry rules
	for _, rule := range validator.Spec.OciRegistryRules {
		username, password := r.secretKeyAuth(req, rule)
		pubKeys := r.signaturePubKeys(req, rule)

		validationResult, err := ociRuleService.ReconcileOciRegistryRule(rule, username, password, pubKeys)
		if err != nil {
			r.Log.V(0).Error(err, "failed to reconcile OCI Registry rule")
		}
		vres.SafeUpdateValidationResult(r.Client, nn, validationResult, validator.Spec.ResultCount(), err, r.Log)
	}

	r.Log.V(0).Info("Requeuing for re-validation in two minutes.", "name", req.Name, "namespace", req.Namespace)
	return ctrl.Result{RequeueAfter: time.Second * 120}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OciValidatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.OciValidator{}).
		Complete(r)
}

func (r *OciValidatorReconciler) secretKeyAuth(req ctrl.Request, rule v1alpha1.OciRegistryRule) (string, string) {
	if rule.Auth.SecretName == "" {
		return "", ""
	}

	authSecret := &corev1.Secret{}
	secretName := rule.Auth.SecretName
	nn := ktypes.NamespacedName{Name: secretName, Namespace: req.Namespace}

	if err := r.Get(context.Background(), nn, authSecret); err != nil {
		if apierrs.IsNotFound(err) {
			// no secrets found, set creds to empty string
			r.Log.V(0).Error(err, fmt.Sprintf("Auth secret %s not found for rule %s", secretName, rule.Name()))
			return "", ""
		} else {
			r.Log.V(0).Error(err, fmt.Sprintf("Failed to fetch auth secret %s for rule %s", secretName, rule.Name()))
			return "", ""
		}
	}

	errMalformedSecret := fmt.Errorf("malformed secret %s/%s", authSecret.Namespace, authSecret.Name)
	username, ok := authSecret.Data["username"]
	if !ok {
		r.Log.V(0).Error(errMalformedSecret, "Auth secret missing username, defaulting to empty username", "name", secretName, "namespace", req.Namespace)
	}

	password, ok := authSecret.Data["password"]
	if !ok {
		r.Log.V(0).Error(errMalformedSecret, "Auth secret missing password, defaulting to empty password", "name", secretName, "namespace", req.Namespace)
	}

	return string(username), string(password)
}

func (r *OciValidatorReconciler) signaturePubKeys(req ctrl.Request, rule v1alpha1.OciRegistryRule) [][]byte {
	if rule.SignatureVerification.SecretName == "" {
		return nil
	}

	pubKeysSecret := &corev1.Secret{}
	secretName := rule.SignatureVerification.SecretName
	nn := ktypes.NamespacedName{Name: secretName, Namespace: req.Namespace}

	// make a slice of byte slices
	pubKeys := make([][]byte, 0)

	if err := r.Get(context.Background(), nn, pubKeysSecret); err != nil {
		if apierrs.IsNotFound(err) {
			// no secrets found, set creds to empty string
			r.Log.V(0).Error(err, fmt.Sprintf("Public Keys secret %s not found for rule %s", secretName, rule.Name()))
			return pubKeys
		} else {
			r.Log.V(0).Error(err, fmt.Sprintf("Failed to fetch Public Keys secret %s for rule %s", secretName, rule.Name()))
			return pubKeys
		}
	}

	for k, data := range pubKeysSecret.Data {
		// search fro public keys in the secret
		if strings.HasSuffix(k, ".pub") {
			pubKeys = append(pubKeys, data)
		}
	}

	return pubKeys
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
					Controller: ptr.Ptr(true),
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
