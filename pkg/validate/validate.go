// Package validate defines a Validate function that evaluates an OciValidatorSpec and returns a ValidationResponse.
package validate

import (
	"github.com/go-logr/logr"

	"github.com/validator-labs/validator/pkg/types"

	"github.com/validator-labs/validator-plugin-oci/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-oci/pkg/auth"
	"github.com/validator-labs/validator-plugin-oci/pkg/oci"
	ocic "github.com/validator-labs/validator-plugin-oci/pkg/ociclient"
)

// Validate evaluates an OciValidatorSpec and returns a ValidationResponse.
func Validate(spec v1alpha1.OciValidatorSpec, auths [][]string, pubKeys [][][]byte, log logr.Logger) types.ValidationResponse {
	resp := types.ValidationResponse{
		ValidationRuleResults: make([]*types.ValidationRuleResult, 0, spec.ResultCount()),
		ValidationRuleErrors:  make([]error, 0, spec.ResultCount()),
	}

	// OCI Registry rules
	for i, rule := range spec.OciRegistryRules {
		vrr := oci.BuildValidationResult(rule)

		ociClient, err := ocic.NewOCIClient(
			ocic.WithBasicAuth(auths[i][0], auths[i][1]),
			ocic.WithMultiAuth(auth.GetKeychain(rule.Host)),
			ocic.WithTLSConfig(rule.InsecureSkipTLSVerify, rule.CaCert, ""),
			ocic.WithVerificationPublicKeys(pubKeys[i]),
		)
		if err != nil {
			log.Error(err, "failed to create OCI client", "ruleName", rule.Name)
			resp.AddResult(vrr, err)
			continue
		}

		svc := oci.NewRuleService(log, oci.WithOCIClient(ociClient))

		vrr, err = svc.ReconcileOciRegistryRule(rule)
		if err != nil {
			log.Error(err, "failed to reconcile OCI Registry rule", "ruleName", rule.Name)
		}
		resp.AddResult(vrr, err)
	}

	return resp
}
