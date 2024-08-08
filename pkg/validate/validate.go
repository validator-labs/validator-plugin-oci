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
func Validate(spec v1alpha1.OciValidatorSpec, authMap map[string][]string, pubKeyMap map[string][][]byte, log logr.Logger) types.ValidationResponse {
	resp := types.ValidationResponse{
		ValidationRuleResults: make([]*types.ValidationRuleResult, 0, spec.ResultCount()),
		ValidationRuleErrors:  make([]error, 0, spec.ResultCount()),
	}

	// OCI Registry rules
	for _, rule := range spec.OciRegistryRules {
		vrr := oci.BuildValidationResult(rule)

		opts := []ocic.Option{
			ocic.WithMultiAuth(auth.GetKeychain(rule.Host)),
			ocic.WithTLSConfig(rule.InsecureSkipTLSVerify, rule.CaCert, ""),
		}
		if pubKeys, ok := pubKeyMap[rule.Name()]; ok {
			opts = append(opts, ocic.WithVerificationPublicKeys(pubKeys))
		}

		if auths, ok := authMap[rule.Name()]; ok {
			if len(auths) == 2 {
				opts = append(opts, ocic.WithBasicAuth(auths[0], auths[1]))
			}
		}

		ociClient, err := ocic.NewOCIClient(opts...)
		if err != nil {
			log.Error(err, "failed to create OCI client", "ruleName", rule.Name())
			resp.AddResult(vrr, err)
			continue
		}

		svc := oci.NewRuleService(log, oci.WithOCIClient(ociClient))

		vrr, err = svc.ReconcileOciRegistryRule(rule)
		if err != nil {
			log.Error(err, "failed to reconcile OCI Registry rule", "ruleName", rule.Name())
		}
		resp.AddResult(vrr, err)
	}

	return resp
}
