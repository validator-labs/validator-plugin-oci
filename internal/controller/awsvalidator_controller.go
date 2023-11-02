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
	"time"

	"github.com/go-logr/logr"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ktypes "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/spectrocloud-labs/validator-plugin-aws/api/v1alpha1"
	"github.com/spectrocloud-labs/validator-plugin-aws/internal/constants"
	aws_utils "github.com/spectrocloud-labs/validator-plugin-aws/internal/utils/aws"
	"github.com/spectrocloud-labs/validator-plugin-aws/internal/validators/iam"
	"github.com/spectrocloud-labs/validator-plugin-aws/internal/validators/servicequota"
	"github.com/spectrocloud-labs/validator-plugin-aws/internal/validators/tag"
	vapi "github.com/spectrocloud-labs/validator/api/v1alpha1"
	"github.com/spectrocloud-labs/validator/pkg/types"
	vres "github.com/spectrocloud-labs/validator/pkg/validationresult"
)

// AwsValidatorReconciler reconciles a AwsValidator object
type AwsValidatorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=validation.spectrocloud.labs,resources=awsvalidators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=validation.spectrocloud.labs,resources=awsvalidators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=validation.spectrocloud.labs,resources=awsvalidators/finalizers,verbs=update

// Reconcile reconciles each rule found in each AWSValidator in the cluster and creates ValidationResults accordingly
func (r *AwsValidatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.V(0).Info("Reconciling AwsValidator", "name", req.Name, "namespace", req.Namespace)

	validator := &v1alpha1.AwsValidator{}
	if err := r.Get(ctx, req.NamespacedName, validator); err != nil {
		// Ignore not-found errors, since they can't be fixed by an immediate requeue
		if apierrs.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		r.Log.Error(err, "failed to fetch AwsValidator")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get the active validator's validation result
	vr := &vapi.ValidationResult{}
	nn := ktypes.NamespacedName{
		Name:      fmt.Sprintf("validator-plugin-aws-%s", validator.Name),
		Namespace: req.Namespace,
	}
	if err := r.Get(ctx, nn, vr); err == nil {
		res, err := vres.HandleExistingValidationResult(nn, vr, r.Log)
		if res != nil {
			return *res, err
		}
	} else {
		if !apierrs.IsNotFound(err) {
			r.Log.V(0).Error(err, "unexpected error getting ValidationResult", "name", nn.Name, "namespace", nn.Namespace)
		}
		res, err := vres.HandleNewValidationResult(r.Client, constants.PluginCode, nn, vr, r.Log)
		if res != nil {
			return *res, err
		}
	}

	failed := &types.MonotonicBool{}

	// IAM rules
	awsApi, err := aws_utils.NewAwsApi(validator.Spec.DefaultRegion)
	if err != nil {
		r.Log.V(0).Error(err, "failed to get AWS client")
	} else {
		iamRuleService := iam.NewIAMRuleService(r.Log, awsApi.IAM)

		for _, rule := range validator.Spec.IamRoleRules {
			validationResult, err := iamRuleService.ReconcileIAMRoleRule(rule)
			if err != nil {
				r.Log.V(0).Error(err, "failed to reconcile IAM role rule")
			}
			vres.SafeUpdateValidationResult(r.Client, nn, validationResult, failed, err, r.Log)
		}
		for _, rule := range validator.Spec.IamUserRules {
			validationResult, err := iamRuleService.ReconcileIAMUserRule(rule)
			if err != nil {
				r.Log.V(0).Error(err, "failed to reconcile IAM user rule")
			}
			vres.SafeUpdateValidationResult(r.Client, nn, validationResult, failed, err, r.Log)
		}
		for _, rule := range validator.Spec.IamGroupRules {
			validationResult, err := iamRuleService.ReconcileIAMGroupRule(rule)
			if err != nil {
				r.Log.V(0).Error(err, "failed to reconcile IAM group rule")
			}
			vres.SafeUpdateValidationResult(r.Client, nn, validationResult, failed, err, r.Log)
		}
		for _, rule := range validator.Spec.IamPolicyRules {
			validationResult, err := iamRuleService.ReconcileIAMPolicyRule(rule)
			if err != nil {
				r.Log.V(0).Error(err, "failed to reconcile IAM policy rule")
			}
			vres.SafeUpdateValidationResult(r.Client, nn, validationResult, failed, err, r.Log)
		}
	}

	// Service Quota rules
	for _, rule := range validator.Spec.ServiceQuotaRules {
		awsApi, err := aws_utils.NewAwsApi(rule.Region)
		if err != nil {
			r.Log.V(0).Error(err, "failed to reconcile Service Quota rule")
			continue
		}
		svcQuotaService := servicequota.NewServiceQuotaRuleService(
			r.Log,
			awsApi.EC2,
			awsApi.EFS,
			awsApi.ELB,
			awsApi.ELBV2,
			awsApi.SQ,
		)
		validationResult, err := svcQuotaService.ReconcileServiceQuotaRule(rule)
		if err != nil {
			r.Log.V(0).Error(err, "failed to reconcile Service Quota rule")
		}
		vres.SafeUpdateValidationResult(r.Client, nn, validationResult, failed, err, r.Log)
	}

	// Tag rules
	for _, rule := range validator.Spec.TagRules {
		awsApi, err := aws_utils.NewAwsApi(rule.Region)
		if err != nil {
			r.Log.V(0).Error(err, "failed to reconcile Tag rule")
			continue
		}
		tagRuleService := tag.NewTagRuleService(r.Log, awsApi.EC2)
		validationResult, err := tagRuleService.ReconcileTagRule(rule)
		if err != nil {
			r.Log.V(0).Error(err, "failed to reconcile Tag rule")
		}
		vres.SafeUpdateValidationResult(r.Client, nn, validationResult, failed, err, r.Log)
	}

	r.Log.V(0).Info("Requeuing for re-validation in two minutes.", "name", req.Name, "namespace", req.Namespace)
	return ctrl.Result{RequeueAfter: time.Second * 120}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AwsValidatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.AwsValidator{}).
		Complete(r)
}
