package validators

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-logr/logr"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"

	"github.com/spectrocloud-labs/validator-plugin-oci/api/v1alpha1"
	"github.com/spectrocloud-labs/validator-plugin-oci/internal/constants"
	vapi "github.com/spectrocloud-labs/validator/api/v1alpha1"
	"github.com/spectrocloud-labs/validator/pkg/types"
	vapitypes "github.com/spectrocloud-labs/validator/pkg/types"
	"github.com/spectrocloud-labs/validator/pkg/util/ptr"
)

type OciRuleService struct {
	log logr.Logger
}

func NewOciRuleService(log logr.Logger) *OciRuleService {
	return &OciRuleService{
		log: log,
	}
}

// ReconcileOciRegistryRule reconciles an OCI registry rule from the OCIValidator config
func (s *OciRuleService) ReconcileOciRegistryRule(rule v1alpha1.OciRegistryRule) (*vapitypes.ValidationResult, error) {

	// Build the default ValidationResult for this rule
	vr := buildValidationResult(rule)

	// Create a new registry instance
	reg, err := remote.NewRegistry(rule.Host)
	if err != nil {
		errMsg := "failed to create registry client"
		s.log.V(0).Info(errMsg, "error", err.Error(), "rule", rule.Name())
		s.updateResult(vr, []error{err}, errMsg, rule.Name())
		return vr, nil
	}

	httpClient, err := newHTTPClient(rule.CaCert)
	if err != nil {
		errMsg := "failed to create http client"
		s.log.V(0).Info(errMsg, "error", err.Error(), "rule", rule.Name(), "caCert", rule.CaCert)
		s.updateResult(vr, []error{err}, errMsg, rule.Name())
		return vr, nil
	}

	// Setup credentials if username and password are provided
	if rule.Auth.Basic.Username != "" && rule.Auth.Basic.Password != "" {
		reg.Client = &auth.Client{
			Client: httpClient,
			Cache:  auth.DefaultCache,
			Credential: auth.StaticCredential(rule.Host, auth.Credential{
				Username: rule.Auth.Basic.Username,
				Password: rule.Auth.Basic.Password,
			}),
		}
	}

	errs := make([]error, 0)
	details := make([]string, 0)
	ctx := context.Background()
	if len(rule.Artifacts) == 0 {
		errMsg := "failed to validate repos in registry"
		detail, err := s.validateRepos(ctx, reg, rule.Host, vr)
		if err != nil {
			errs = append(errs, err)
			s.log.V(0).Info(errMsg, "error", err.Error(), "rule", rule.Name(), "host", rule.Host)
		}
		if detail != "" {
			details = append(details, detail)
		}

		s.updateResult(vr, errs, errMsg, rule.Name(), details...)
		return vr, nil
	}

	errMsg := "failed to validate artifact in registry"
	for _, artifact := range rule.Artifacts {
		detail, err := s.validateArtifact(ctx, reg, rule.Host, artifact, vr)
		if err != nil {
			errs = append(errs, err)
			s.log.V(0).Info(errMsg, "error", err.Error(), "rule", rule.Name(), "host", rule.Host, "artifact", artifact)
		}
		if detail != "" {
			details = append(details, detail)
		}
	}
	s.updateResult(vr, errs, errMsg, rule.Name(), details...)
	return vr, nil
}

// validateRepo validates repos within a registry. This function is to be used when no particular repo or artifact is specified
func (s *OciRuleService) validateRepos(ctx context.Context, reg *remote.Registry, host string, vr *types.ValidationResult) (string, error) {
	// Get chosen repositories
	repoList := make([]string, 0)
	err := reg.Repositories(ctx, "", func(repos []string) error {
		for _, repo := range repos {
			repoList = append(repoList, repo)
		}
		return nil
	})
	if err != nil {
		return fmt.Sprintf("failed to list repositories in registry"), err
	}

	if len(repoList) == 0 {
		return fmt.Sprintf("no repositories found in registry"), nil
	}

	var repo registry.Repository
	for _, curRepo := range repoList {
		// Try to get repo from registry
		if repo, err = reg.Repository(ctx, curRepo); err == nil {
			break
		}
	}
	if repo == nil || err != nil {
		return fmt.Sprintf("failed to authenticate to a repository"), err
	}

	// Get tags on artifacts in repository
	var tagList []string
	err = repo.Tags(ctx, "", func(tags []string) error {
		for _, tag := range tags {
			tagList = append(tagList, tag)
		}
		return nil
	})
	if err != nil {
		return fmt.Sprintf("failed to pull tags from repository"), err
	}

	return "", nil
}

// validateArtifact validates a single artifact within a registry. This function is to be used when a path to a repo or an individual artifact is provided
func (s *OciRuleService) validateArtifact(ctx context.Context, reg *remote.Registry, host string, repoPath string, vr *types.ValidationResult) (string, error) {
	path, ref, err := parseArtifact(repoPath)
	if err != nil {
		return fmt.Sprintf("failed to get artifact path and reference"), err
	}

	// Try to get repo from registry
	repo, err := reg.Repository(ctx, path)
	if err != nil {
		return fmt.Sprintf("failed to authenticate to a repository"), err
	}

	if ref == "" {
		// Get all tags on artifacts in repository
		var tagList []string
		err = repo.Tags(ctx, "", func(tags []string) error {
			for _, tag := range tags {
				tagList = append(tagList, tag)
			}
			return nil
		})
		if err != nil {
			return fmt.Sprintf("failed to pull tags from repository"), err
		}
		return "", nil
	}

	// Get reference of artifact
	_, _, err = repo.FetchReference(ctx, ref)
	if err != nil {
		return fmt.Sprintf("failed to fetch reference to artifact"), err
	}

	return "", nil
}

func newHTTPClient(caCert string) (*http.Client, error) {
	httpClient := retry.DefaultClient

	// Add cert as trust
	if caCert != "" {
		cert, err := base64.StdEncoding.DecodeString(caCert)
		if err != nil {
			return nil, err
		}

		caCertPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		caCertPool.AppendCertsFromPEM(cert)

		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    caCertPool,
		}
		httpClient.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	}

	return httpClient, nil
}

func parseArtifact(artifactPath string) (string, string, error) {
	path := ""
	ref := ""
	parts := strings.Split(artifactPath, "@")
	if len(parts) > 2 {
		return "", "", errors.New("invalid artifact path")
	}

	path = parts[0]
	if len(parts) > 1 {
		ref = parts[1]
	}

	return path, ref, nil
}

// buildValidationResult builds a default ValidationResult for a given validation type
func buildValidationResult(rule v1alpha1.OciRegistryRule) *types.ValidationResult {
	state := vapi.ValidationSucceeded
	latestCondition := vapi.DefaultValidationCondition()
	latestCondition.Details = make([]string, 0)
	latestCondition.Failures = make([]string, 0)
	latestCondition.Message = fmt.Sprintf("All %s checks passed", constants.OciRegistry)
	latestCondition.ValidationRule = rule.Name()
	latestCondition.ValidationType = constants.OciRegistry
	return &types.ValidationResult{Condition: &latestCondition, State: &state}
}

func (s *OciRuleService) updateResult(vr *types.ValidationResult, errs []error, errMsg, ruleName string, details ...string) {
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
