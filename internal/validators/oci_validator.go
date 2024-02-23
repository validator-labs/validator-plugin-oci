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

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/go-logr/logr"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/validate"

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
func (s *OciRuleService) ReconcileOciRegistryRule(rule v1alpha1.OciRegistryRule, username, password string, pubKeys [][]byte) (*vapitypes.ValidationResult, error) {
	vr := buildValidationResult(rule)

	opts, err := setupTransportOpts([]remote.Option{}, rule.CaCert)
	if err != nil {
		errMsg := "failed to setup http client transport"
		s.log.V(0).Info(errMsg, "error", err.Error(), "rule", rule.Name(), "caCert", rule.CaCert)
		s.updateResult(vr, []error{err}, errMsg, rule.Name())
		return vr, err
	}

	opts, err = setupAuthOpts(opts, rule.Host, username, password)
	if err != nil {
		errMsg := "failed to setup authentication"
		s.log.V(0).Info(errMsg, "error", err.Error(), "rule", rule.Name(), "host", rule.Host, "auth", rule.Auth)
		s.updateResult(vr, []error{err}, errMsg, rule.Name())
		return vr, err
	}

	errs := make([]error, 0)
	details := make([]string, 0)
	ctx := context.Background()
	if len(rule.Artifacts) == 0 {
		errMsg := "failed to validate repositories in registry"
		d, err := validateRepos(ctx, rule.Host, opts, vr)
		if err != nil {
			errs = append(errs, err)
			s.log.V(0).Info(errMsg, "error", err.Error(), "rule", rule.Name(), "host", rule.Host)
		}
		if len(d) > 0 {
			details = append(details, d...)
		}

		s.updateResult(vr, errs, errMsg, rule.Name(), details...)
		return vr, err
	}

	errMsg := "failed to validate artifact in registry"
	for _, artifact := range rule.Artifacts {
		ref, err := generateRef(rule.Host, artifact.Ref, vr)
		if err != nil {
			details = append(details, fmt.Sprintf("failed to generate reference for artifact %s/%s", rule.Host, artifact.Ref))
			errs = append(errs, err)
			s.log.V(0).Info(errMsg, "error", err.Error(), "rule", rule.Name(), "host", rule.Host, "artifact", artifact)
			continue
		}

		detail, err := validateReference(ref, artifact.LayerValidation, opts)
		if err != nil {
			errs = append(errs, err)
			s.log.V(0).Info(errMsg, "error", err.Error(), "rule", rule.Name(), "host", rule.Host, "artifact", artifact)
		}
		if detail != "" {
			details = append(details, detail)
		}
	}
	s.updateResult(vr, errs, errMsg, rule.Name(), details...)

	if len(errs) > 0 {
		return vr, errs[0]
	}
	return vr, nil
}

func generateRef(registry, artifact string, vr *types.ValidationResult) (name.Reference, error) {
	if strings.Contains(artifact, "@sha256:") {
		return name.NewDigest(fmt.Sprintf("%s/%s", registry, artifact))
	}

	if !strings.Contains(artifact, ":") {
		vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf("artifact %s does not contain a tag or digest, defaulting to \"latest\" tag", artifact))
	}

	return name.NewTag(fmt.Sprintf("%s/%s", registry, artifact))
}

// validateArtifact validates a single artifact within a registry. This function is to be used when a path to a repo or an individual artifact is provided
func validateReference(ref name.Reference, fullLayerValidation bool, opts []remote.Option) (string, error) {
	// validate artifact existence by issuing a HEAD request
	_, err := remote.Head(ref, opts...)
	if err != nil {
		return fmt.Sprintf("failed to get descriptor for artifact %s", ref.Name()), err
	}

	// download image without storing it on disk
	img, err := remote.Image(ref, opts...)
	if err != nil {
		return fmt.Sprintf("failed to download artifact %s", ref.Name()), err
	}

	var validateOpts []validate.Option
	if !fullLayerValidation {
		validateOpts = append(validateOpts, validate.Fast)
	}

	err = validate.Image(img, validateOpts...)
	if err != nil {
		return fmt.Sprintf("failed validation for artifact %s", ref.Name()), err
	}

	// TODO: add signature verification here

	return "", nil
}

// validateRepos validates repos within a registry. This function is to be used when no particular artifact in a registry is provided
func validateRepos(ctx context.Context, host string, opts []remote.Option, vr *types.ValidationResult) ([]string, error) {
	details := make([]string, 0)

	reg, err := name.NewRegistry(host)
	if err != nil {
		details = append(details, fmt.Sprintf("failed to get registry %s", host))
		return details, err
	}

	repoList, err := remote.Catalog(context.Background(), reg, opts...)
	if err != nil {
		details = append(details, fmt.Sprintf("failed to list repositories in registry %s", host))
		return details, err
	}

	if len(repoList) == 0 {
		details = append(details, fmt.Sprintf("no repositories found in registry %s", host))
		return details, nil
	}

	var repo name.Repository
	var ref name.Reference
	for _, curRepo := range repoList {
		repo, err = name.NewRepository(host + "/" + curRepo)
		if err != nil {
			details = append(details, fmt.Sprintf("failed to get repository %s/%s", host, curRepo))
			continue
		}

		tags, err := remote.List(repo, opts...)
		if err != nil {
			details = append(details, fmt.Sprintf("failed to get tags for repository %s/%s", host, curRepo))
			continue
		}

		if len(tags) == 0 {
			details = append(details, fmt.Sprintf("no tags found in repository %s/%s", host, curRepo))
			continue
		}

		tag := tags[0]
		ref, err = generateRef(host, fmt.Sprintf("%s:%s", curRepo, tag), vr)
		if err != nil {
			details = append(details, fmt.Sprintf("failed to generate reference for artifact %s/%s:%s", host, curRepo, tag))
			continue
		}
		break
	}
	if err != nil {
		return details, err
	}

	detail, err := validateReference(ref, true, opts)
	if err != nil {
		details = append(details, detail)
		return details, err
	}

	return []string{}, nil
}

func getEcrLoginToken(username, password, region string) (string, error) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(username, password, "")),
	)
	if err != nil {
		return "", fmt.Errorf("failed to load configuration, %s", err)
	}

	client := ecr.NewFromConfig(cfg)

	output, err := client.GetAuthorizationToken(context.Background(), &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return "", fmt.Errorf("failed to get ECR authorization token: %s", err)
	}

	for _, data := range output.AuthorizationData {
		return *data.AuthorizationToken, nil
	}

	return "", fmt.Errorf("no authorization data available")
}

func parseEcrRegion(url string) (string, error) {
	parts := strings.Split(url, ".")
	if len(parts) != 6 || parts[2] != "ecr" {
		return "", errors.New(fmt.Sprintf("Invalid ECR URL %s", url))
	}

	region := parts[3]
	return region, nil
}

func setupAuthOpts(opts []remote.Option, registryName, username, password string) ([]remote.Option, error) {
	var auth authn.Authenticator

	if strings.Contains(registryName, "amazonaws.com") {
		region, err := parseEcrRegion(registryName)
		if err != nil {
			return nil, err
		}

		token, err := getEcrLoginToken(username, password, region)
		if err != nil {
			return nil, err
		}

		decodedToken, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			return nil, err
		}

		parts := strings.SplitN(string(decodedToken), ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed ECR login token: %s", token)
		}

		username = parts[0]
		password = parts[1]
	}

	if username != "" && password != "" {
		auth = &authn.Basic{Username: username, Password: password}
	} else {
		auth = authn.Anonymous
	}

	opts = append(opts, remote.WithAuth(auth))
	return opts, nil
}

func setupTransportOpts(opts []remote.Option, caCert string) ([]remote.Option, error) {
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
		opts = append(opts, remote.WithTransport(&http.Transport{TLSClientConfig: tlsConfig}))
	}
	return opts, nil
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
