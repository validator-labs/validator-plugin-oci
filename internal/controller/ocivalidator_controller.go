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
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ktypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	vapi "github.com/validator-labs/validator/api/v1alpha1"
	vres "github.com/validator-labs/validator/pkg/validationresult"

	"github.com/validator-labs/validator-plugin-oci/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-oci/pkg/constants"
	"github.com/validator-labs/validator-plugin-oci/pkg/validate"
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
		Name:      vres.Name(validator),
		Namespace: req.Namespace,
	}
	if err := r.Get(ctx, nn, vr); err == nil {
		vres.HandleExisting(vr, r.Log)
	} else {
		if !apierrs.IsNotFound(err) {
			l.Error(err, "unexpected error getting ValidationResult")
		}
		if err := vres.HandleNew(ctx, r.Client, p, vres.Build(validator), r.Log); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: time.Millisecond}, nil
	}

	// Always update the expected result count in case the validator's rules have changed
	vr.Spec.ExpectedResults = validator.Spec.ResultCount()

	// Fetch OCI registry basic auth secrets and signature verification public keys
	auths := make(map[string][]string)
	allPubKeys := make(map[string][][]byte)

	for _, rule := range validator.Spec.OciRegistryRules {
		username, password, err := r.auth(req, rule)
		if err != nil {
			l.Error(err, "failed to get secret auth", "ruleName", rule.Name())
			return ctrl.Result{}, err
		}
		if username != "" || password != "" {
			auths[rule.Name()] = []string{username, password}
		}

		pubKeys, err := r.signaturePubKeys(req, rule)
		if err != nil {
			l.Error(err, "failed to get signature verification public keys", "ruleName", rule.Name())
			return ctrl.Result{}, err
		}
		if pubKeys != nil {
			allPubKeys[rule.Name()] = pubKeys
		}
	}

	// Validate the rules
	resp := validate.Validate(validator.Spec, auths, allPubKeys, r.Log)

	// Patch the ValidationResult with the latest ValidationRuleResults
	if err := vres.SafeUpdate(ctx, p, vr, resp, r.Log); err != nil {
		return ctrl.Result{}, err
	}

	// Capturing reconciliation frequency from annotations .
	// Defaulting to 120 seconds if annotation is not found.
	var frequency time.Duration
	if secondsString, ok := validator.Annotations[constants.ReconciliationFrequencyAnnotation]; ok {
		seconds, err := strconv.Atoi(secondsString)
		if err != nil {
			l.Info("failed to convert frequency annotation: defaulting to 120 seconds")
			frequency = time.Second * 120
		} else {
			frequency = time.Second * time.Duration(seconds)
		}
	} else {
		frequency = time.Second * 120
	}

	l.Info("Requeuing for re-validation.")
	return ctrl.Result{RequeueAfter: frequency}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OciValidatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.OciValidator{}).
		Complete(r)
}

// auth retrieves the username and password from the rule's auth field.
// If a secretName is provided in the rule's auth field, then the secret is fetched and the username and password are retrieved from the secret.
// Any additional key-value pairs in the secret are set as environment variables, to be picked up by auth keychains (e.g. ECR, Azure).
func (r *OciValidatorReconciler) auth(req ctrl.Request, rule v1alpha1.OciRegistryRule) (string, string, error) {
	if rule.Auth.SecretName != nil {
		return r.secretKeyAuth(req, rule)
	}

	if rule.Auth.Basic != nil {
		return rule.Auth.Basic.Username, rule.Auth.Basic.Password, nil
	}

	return "", "", nil
}

func (r *OciValidatorReconciler) secretKeyAuth(req ctrl.Request, rule v1alpha1.OciRegistryRule) (string, string, error) {
	if *rule.Auth.SecretName == "" {
		return "", "", nil
	}

	authSecret := &corev1.Secret{}
	nn := ktypes.NamespacedName{Name: *rule.Auth.SecretName, Namespace: req.Namespace}

	if err := r.Get(context.Background(), nn, authSecret); err != nil {
		return "", "", fmt.Errorf("failed to fetch auth secret %s/%s for rule %s: %w", nn.Namespace, nn.Name, rule.Name(), err)
	}

	if len(authSecret.Data) == 0 {
		return "", "", fmt.Errorf("secret %s/%s has no data", nn.Namespace, nn.Name)
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

// signaturePubKeys retrieves the public keys that are used for signature verification.
// If a secretName is provided in the SignatureVerification field, then the secret is fetched and the pub keys are retrieved from the secret.
// Otherwise, the public keys are retrieved from the inline PublicKeys field if provided.
func (r *OciValidatorReconciler) signaturePubKeys(req ctrl.Request, rule v1alpha1.OciRegistryRule) ([][]byte, error) {
	if rule.SignatureVerification.SecretName != "" {
		return r.signaturePubKeysSecret(req, rule)
	}

	if len(rule.SignatureVerification.PublicKeys) > 0 {
		return r.signaturePubKeysInline(rule.SignatureVerification.PublicKeys), nil
	}

	return nil, nil
}
func (r *OciValidatorReconciler) signaturePubKeysSecret(req ctrl.Request, rule v1alpha1.OciRegistryRule) ([][]byte, error) {
	pubKeysSecret := &corev1.Secret{}
	nn := ktypes.NamespacedName{Name: rule.SignatureVerification.SecretName, Namespace: req.Namespace}

	if err := r.Get(context.Background(), nn, pubKeysSecret); err != nil {
		return nil, fmt.Errorf("failed to fetch public keys secret %s/%s for rule %s: %w",
			nn.Namespace, nn.Name, rule.Name(), err,
		)
	}

	pubKeys := make([][]byte, 0)
	for k, data := range pubKeysSecret.Data {
		if strings.HasSuffix(k, ".pub") {
			pubKeys = append(pubKeys, data)
		}
	}
	if len(pubKeys) == 0 {
		return nil, fmt.Errorf("no public keys found in secret %s/%s for rule: %s", nn.Namespace, nn.Name, rule.Name())
	}

	return pubKeys, nil
}

func (r *OciValidatorReconciler) signaturePubKeysInline(pKeys []string) [][]byte {
	pubKeys := make([][]byte, len(pKeys))

	for i, key := range pKeys {
		pubKeys[i] = []byte(key)
	}

	return pubKeys
}
