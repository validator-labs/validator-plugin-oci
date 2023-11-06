package ecr

import (
	"fmt"

	"github.com/go-logr/logr"

	"github.com/spectrocloud-labs/validator-plugin-oci/api/v1alpha1"
	"github.com/spectrocloud-labs/validator-plugin-oci/internal/constants"
	vapi "github.com/spectrocloud-labs/validator/api/v1alpha1"
	"github.com/spectrocloud-labs/validator/pkg/types"
	vapitypes "github.com/spectrocloud-labs/validator/pkg/types"
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

	s.log.Info("ecr registry rule validation is unimplemented")
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
