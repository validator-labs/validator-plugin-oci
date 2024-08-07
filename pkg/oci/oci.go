// Package oci defines the OCI registry rule service and implements the reconcile function for OCI registry rules.
package oci

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/google/go-containerregistry/pkg/name"

	vapi "github.com/validator-labs/validator/api/v1alpha1"
	"github.com/validator-labs/validator/pkg/types"
	"github.com/validator-labs/validator/pkg/util"

	"github.com/validator-labs/validator-plugin-oci/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-oci/pkg/constants"
	ocic "github.com/validator-labs/validator-plugin-oci/pkg/ociclient"
)

// RuleService defines the service for OCI registry rules.
type RuleService struct {
	log       logr.Logger
	ociClient *ocic.Client
}

// Option is a functional option for configuring a RuleService.
type Option func(*RuleService)

// NewRuleService creates a new OCI registry rule service.
func NewRuleService(log logr.Logger, opts ...Option) *RuleService {
	s := &RuleService{
		log: log,
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// WithOCIClient sets the OCI client for the OCI registry rule service.
func WithOCIClient(client *ocic.Client) Option {
	return func(s *RuleService) {
		s.ociClient = client
	}
}

// ReconcileOciRegistryRule reconciles an OCI registry rule from the OCIValidator config.
func (s *RuleService) ReconcileOciRegistryRule(rule v1alpha1.OciRegistryRule) (*types.ValidationRuleResult, error) {
	l := s.log.V(0).WithValues("rule", rule.Name(), "host", rule.Host)
	vr := BuildValidationResult(rule)

	var err error
	ctx := context.Background()

	errMsg := "failed to validate repositories in registry"
	if len(rule.Artifacts) == 0 {
		details, errs := s.validateRepos(ctx, rule, vr)

		if len(errs) > 0 {
			l.Error(errs[len(errs)-1], errMsg)
			err = errors.New(errMsg)
		}
		s.updateResult(vr, errs, errMsg, details...)

		return vr, err
	}

	details := make([]string, 0)
	errs := make([]error, 0)
	errMsg = "failed to validate artifact in registry"

	for _, artifact := range rule.Artifacts {
		ref, err := generateRef(rule.Host, artifact.Ref, vr)
		if err != nil {
			details = append(details, fmt.Sprintf("failed to generate reference for artifact %s/%s", rule.Host, artifact.Ref))
			errs = append(errs, err)
			l.Error(err, errMsg, "artifact", artifact)
			continue
		}

		skipLayerValidation := s.shouldSkipLayerValidation(rule, artifact)

		d, e := s.validateReference(ctx, ref, skipLayerValidation, rule.SignatureVerification)
		if len(e) > 0 {
			l.Error(e[len(e)-1], errMsg, "artifact", artifact)
		}
		details = append(details, d...)
		errs = append(errs, e...)
	}

	if len(errs) > 0 {
		err = errors.New(errMsg)
	}
	s.updateResult(vr, errs, errMsg, details...)

	return vr, err
}

func (s *RuleService) shouldSkipLayerValidation(rule v1alpha1.OciRegistryRule, artifact v1alpha1.Artifact) bool {
	if artifact.SkipLayerValidation == nil {
		return rule.SkipLayerValidation
	}
	return *artifact.SkipLayerValidation
}

// validateArtifact validates a single artifact within an OCI registry. Used when either a path to a repo or artifact(s) are specified in an OciRegistryRule.
func (s *RuleService) validateReference(ctx context.Context, ref name.Reference, skipLayerValidation bool, sv v1alpha1.SignatureVerification) ([]string, []error) {
	details := make([]string, 0)
	errs := make([]error, 0)

	// validate artifact existence by issuing a HEAD request
	if _, err := s.ociClient.Head(ref); err != nil {
		details = append(details, fmt.Sprintf("failed to get descriptor for artifact %s", ref.Name()))
		errs = append(errs, err)
		return details, errs
	}

	// download image without storing it on disk
	img, err := s.ociClient.PullImage(ref)
	if err != nil {
		details = append(details, fmt.Sprintf("failed to download artifact %s", ref.Name()))
		errs = append(errs, err)
		return details, errs
	}

	// validate image
	if err := s.ociClient.ValidateImage(img, skipLayerValidation); err != nil {
		details = append(details, fmt.Sprintf("failed validation for artifact %s", ref.Name()))
		errs = append(errs, err)
		return details, errs
	}
	details = append(details, fmt.Sprintf("pulled and validated artifact %s", ref.Name()))

	// verify image signature (optional)
	if sv.SecretName != "" {
		verifyDetails, verifyErrs := s.ociClient.VerifySignature(ctx, ref)
		if len(verifyDetails) > 0 {
			details = append(details, verifyDetails...)
		}
		if len(verifyErrs) > 0 {
			errs = append(errs, verifyErrs...)
		}
	}

	return details, errs
}

// validateRepos validates repos within an OCI registry. Used when no specific artifacts are provided in an OciRegistryRule.
func (s *RuleService) validateRepos(ctx context.Context, rule v1alpha1.OciRegistryRule, vr *types.ValidationRuleResult) ([]string, []error) {
	details := make([]string, 0)
	errs := make([]error, 0)

	reg, err := name.NewRegistry(rule.Host)
	if err != nil {
		details = append(details, fmt.Sprintf("failed to get registry %s", rule.Host))
		errs = append(errs, err)
		return details, errs
	}

	repoList, err := s.ociClient.Catalog(ctx, reg)
	if err != nil {
		details = append(details, fmt.Sprintf("failed to list repositories in registry %s", rule.Host))
		errs = append(errs, err)
		return details, errs
	}

	if len(repoList) == 0 {
		details = append(details, fmt.Sprintf("no repositories found in registry %s", rule.Host))
		return details, errs
	}

	var repo name.Repository
	var ref name.Reference
	foundArtifact := false

	for _, curRepo := range repoList {
		repo, err = name.NewRepository(rule.Host + "/" + curRepo)
		if err != nil {
			errs = append(errs, err)
			details = append(details, fmt.Sprintf("failed to get repository %s/%s", rule.Host, curRepo))
			continue
		}

		tags, err := s.ociClient.List(repo)
		if err != nil {
			errs = append(errs, err)
			details = append(details, fmt.Sprintf("failed to get tags for repository %s/%s", rule.Host, curRepo))
			continue
		}
		if len(tags) == 0 {
			details = append(details, fmt.Sprintf("no tags found in repository %s/%s", rule.Host, curRepo))
			continue
		}

		tag := tags[0]
		ref, err = generateRef(rule.Host, fmt.Sprintf("%s:%s", curRepo, tag), vr)
		if err != nil {
			errs = append(errs, err)
			details = append(details, fmt.Sprintf("failed to generate reference for artifact %s/%s:%s", rule.Host, curRepo, tag))
			continue
		}

		foundArtifact = true
		break
	}
	if !foundArtifact {
		return details, errs
	}

	return s.validateReference(ctx, ref, false, rule.SignatureVerification)
}

func (s *RuleService) updateResult(vr *types.ValidationRuleResult, errs []error, errMsg string, details ...string) {
	if len(errs) > 0 {
		vr.State = util.Ptr(vapi.ValidationFailed)
		vr.Condition.Message = errMsg
		for _, err := range errs {
			vr.Condition.Failures = append(vr.Condition.Failures, err.Error())
		}
	}
	vr.Condition.Details = append(vr.Condition.Details, details...)
}

// BuildValidationResult builds a default ValidationResult for a given validation type.
func BuildValidationResult(rule v1alpha1.OciRegistryRule) *types.ValidationRuleResult {
	state := vapi.ValidationSucceeded
	latestCondition := vapi.DefaultValidationCondition()
	latestCondition.Details = make([]string, 0)
	latestCondition.Failures = make([]string, 0)
	latestCondition.Message = fmt.Sprintf("All %s checks passed", constants.OciRegistry)
	latestCondition.ValidationRule = rule.Name()
	latestCondition.ValidationType = constants.OciRegistry
	return &types.ValidationRuleResult{Condition: &latestCondition, State: &state}
}

// generateRef generates a name.Reference for a given OCI registry and artifact.
func generateRef(registry, artifact string, vr *types.ValidationRuleResult) (name.Reference, error) {
	if strings.Contains(artifact, "@sha256:") {
		return name.NewDigest(fmt.Sprintf("%s/%s", registry, artifact))
	}

	if !strings.Contains(artifact, ":") {
		vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf("artifact %s does not contain a tag or digest, defaulting to \"latest\" tag", artifact))
	}

	return name.NewTag(fmt.Sprintf("%s/%s", registry, artifact))
}
