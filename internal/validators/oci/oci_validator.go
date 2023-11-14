package oci

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/http"

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
		s.log.V(0).Error(err, "failed to create registry client")
		return vr, err
	}

	httpClient, err := newHTTPClient(rule.Cert)
	if err != nil {
		s.log.V(0).Error(err, "failed to create http client", "cert", rule.Cert)
		return vr, err
	}

	// Setup credentials if username and password are provided
	if rule.Auth.Username != "" && rule.Auth.Password != "" {
		reg.Client = &auth.Client{
			Client: httpClient,
			Cache:  auth.DefaultCache,
			Credential: auth.StaticCredential(rule.Host, auth.Credential{
				Username: rule.Auth.Username,
				Password: rule.Auth.Password,
			}),
		}
	}

	ctx := context.Background()
	if len(rule.RepositoryPaths) == 0 {
		err := s.validateRepo(ctx, reg, rule.Host, "", vr)
		if err != nil {
			return vr, err
		}
	} else {
		for _, repo := range rule.RepositoryPaths {
			err := s.validateRepo(ctx, reg, rule.Host, repo, vr)
			if err != nil {
				return vr, err
			}
		}
	}

	return vr, nil
}

func (s *OciRuleService) validateRepo(ctx context.Context, reg *remote.Registry, host string, repoPath string, vr *types.ValidationResult) error {
	// Get chosen repositories
	repoList := make([]string, 0)
	if repoPath == "" {
		err := reg.Repositories(ctx, "", func(repos []string) error {
			for _, repo := range repos {
				repoList = append(repoList, repo)
			}
			return nil
		})
		if err != nil {
			s.log.V(0).Error(err, "failed to list repositories in registry", "host", host)
			return err
		}

		if len(repoList) == 0 {
			// no repositories in registry, not possible to run any further validations
			vr.Condition.Details = append(vr.Condition.Details, "no repositories found in registry")
			return nil
		}
	} else {
		repoList = append(repoList, repoPath)
	}

	var repo registry.Repository
	var err error
	for _, curRepo := range repoList {
		// Try to get repo from regisry
		if repo, err = reg.Repository(ctx, curRepo); err == nil {
			break
		}
	}

	if repo == nil || err != nil {
		s.log.V(0).Error(err, "unable to authenticate to a repository", "host", host, "path", repoPath)
		return err
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
		s.log.V(0).Error(err, "failed to pull tags from repository", "host", host, "path", repoPath)
		return err
	}

	return nil
}

func newHTTPClient(cert string) (*http.Client, error) {
	httpClient := retry.DefaultClient

	// Add cert as trust
	if cert != "" {
		repoCert, err := base64.StdEncoding.DecodeString(cert)
		if err != nil {
			return nil, err
		}

		caCertPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		caCertPool.AppendCertsFromPEM(repoCert)

		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    caCertPool,
		}
		httpClient.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	}

	return httpClient, nil
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
