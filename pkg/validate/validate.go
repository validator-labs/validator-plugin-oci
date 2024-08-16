// Package validate defines a Validate function that evaluates an OciValidatorSpec and returns a ValidationResponse.
package validate

import (
	"os"

	"github.com/go-logr/logr"

	"github.com/validator-labs/validator/pkg/types"

	"github.com/validator-labs/validator-plugin-oci/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-oci/pkg/auth"
	"github.com/validator-labs/validator-plugin-oci/pkg/constants"
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

		if rule.Auth.ECR != nil {
			err := configureEcrEnvVars(rule, log)
			if err != nil {
				log.Error(err, "failed to configure ECR environment variables", "ruleName", rule.Name())
			}
			resp.AddResult(vrr, err)
			continue
		}

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

func configureEcrEnvVars(rule v1alpha1.OciRegistryRule, log logr.Logger) error {
	if err := os.Setenv(constants.AwsAccessKey, rule.Auth.ECR.AccessKeyID); err != nil {
		return err
	}
	log.Info("Set environment variable", "key", constants.AwsAccessKey, "ruleName", rule.Name())

	if err := os.Setenv(constants.AwsSecretAccessKey, rule.Auth.ECR.SecretAccessKey); err != nil {
		return err
	}
	log.Info("Set environment variable", "key", constants.AwsSecretAccessKey, "ruleName", rule.Name())

	if rule.Auth.ECR.SessionToken != "" {
		if err := os.Setenv(constants.AwsSessionToken, rule.Auth.ECR.SessionToken); err != nil {
			return err
		}
		log.Info("Set environment variable", "key", constants.AwsSessionToken, "ruleName", rule.Name())
	}

	return nil
}
