package registryauth

import (
	"github.com/go-logr/logr"

	"github.com/spectrocloud-labs/validator-plugin-oci/api/v1alpha1"
	vapi "github.com/spectrocloud-labs/validator/api/v1alpha1"
	vapitypes "github.com/spectrocloud-labs/validator/pkg/types"
)

type RegistryAuthRuleService struct {
	log logr.Logger
	// TODO: add anything needed from oras api
}

func NewRegistryAuthRuleService(log logr.Logger) *RegistryAuthRuleService {
	return &RegistryAuthRuleService{
		log: log,
	}
}

// ReconcileRegistryAuthRule reconciles an AWS service quota validation rule from the AWSValidator config
func (s *RegistryAuthRuleService) ReconcileRegistryAuthRule(rule v1alpha1.RegistryAuthRule) (*vapitypes.ValidationResult, error) {

	// Build the default latest condition for this tag rule
	state := vapi.ValidationSucceeded
	latestCondition := vapi.DefaultValidationCondition()
	validationResult := &vapitypes.ValidationResult{Condition: &latestCondition, State: &state}

	// TODO: implement registry auth validation rules
	s.log.Info("registry auth rule validation is unimplemented")

	return validationResult, nil
}
