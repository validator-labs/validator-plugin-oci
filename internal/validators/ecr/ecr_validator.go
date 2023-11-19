package ecr

import (
	"fmt"

	"github.com/go-logr/logr"

	"github.com/spectrocloud-labs/validator-plugin-oci/api/v1alpha1"
	"github.com/spectrocloud-labs/validator-plugin-oci/internal/constants"
	vapi "github.com/spectrocloud-labs/validator/api/v1alpha1"
	"github.com/spectrocloud-labs/validator/pkg/types"
	vapitypes "github.com/spectrocloud-labs/validator/pkg/types"
	"github.com/spectrocloud-labs/validator/pkg/util/ptr"
)

type EcrRuleService struct {
	log logr.Logger
}

func NewEcrRuleService(log logr.Logger) *EcrRuleService {
	return &EcrRuleService{
		log: log,
	}
}

// ReconcileEcrRegistryRule reconciles an ECR registry rule from the OCIValidator config
func (s *EcrRuleService) ReconcileEcrRegistryRule(rule v1alpha1.EcrRegistryRule) (*vapitypes.ValidationResult, error) {

	// Build the default ValidationResult for this rule
	vr := buildValidationResult(rule)

	return vr, nil
}

// buildValidationResult builds a default ValidationResult for a given validation type
func buildValidationResult(rule v1alpha1.EcrRegistryRule) *types.ValidationResult {
	state := vapi.ValidationSucceeded
	latestCondition := vapi.DefaultValidationCondition()
	latestCondition.Details = make([]string, 0)
	latestCondition.Failures = make([]string, 0)
	latestCondition.Message = fmt.Sprintf("All %s checks passed", constants.EcrRegistry)
	latestCondition.ValidationRule = rule.Name()
	latestCondition.ValidationType = constants.EcrRegistry
	return &types.ValidationResult{Condition: &latestCondition, State: &state}
}

func (s *EcrRuleService) updateResult(vr *types.ValidationResult, errs []error, errMsg, ruleName string, details ...string) {
	if len(errs) > 0 {
		vr.State = ptr.Ptr(vapi.ValidationFailed)
		vr.Condition.Message = errMsg
		for _, err := range errs {
			vr.Condition.Failures = append(vr.Condition.Failures, err.Error())
		}
	}
	for _, detail := range details {
		vr.Condition.Details = append(vr.Condition.Details, detail)
	}
}
