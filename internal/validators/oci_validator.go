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
	vr := buildValidationResult(rule)

	opts, err := setupTransportOpts([]remote.Option{}, rule.CaCert)
	if err != nil {
		errMsg := "failed to setup http client transport"
		s.log.V(0).Info(errMsg, "error", err.Error(), "rule", rule.Name(), "caCert", rule.CaCert)
		s.updateResult(vr, []error{err}, errMsg, rule.Name())
		return vr, nil
	}

	opts, err = setupAuthOpts(opts, rule.Host, rule.Auth)
	if err != nil {
		errMsg := "failed to setup authentication"
		s.log.V(0).Info(errMsg, "error", err.Error(), "rule", rule.Name(), "host", rule.Host, "auth", rule.Auth)
		s.updateResult(vr, []error{err}, errMsg, rule.Name())
		return vr, nil
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
		return vr, nil
	}

	errMsg := "failed to validate artifact in registry"
	for _, artifact := range rule.Artifacts {
		detail, err := validateArtifact(ctx, rule.Host, artifact, opts, vr)
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

// validateArtifact validates a single artifact within a registry. This function is to be used when a path to a repo or an individual artifact is provided
func validateArtifact(ctx context.Context, host string, artifact v1alpha1.Artifact, opts []remote.Option, vr *types.ValidationResult) (string, error) {
	ref, err := generateRef(host, artifact.Ref, vr)
	if err != nil {
		return fmt.Sprintf("failed to generate reference for artifact %v/%v", host, artifact), err
	}

	if artifact.Download {
		// download image without storing it on disk
		_, err = remote.Image(ref, opts...)
		if err != nil {
			return fmt.Sprintf("failed to download artifact %v", ref.Name()), err
		}
	}

	return "", nil
}

// validateRepos validates repos within a registry. This function is to be used when no particular artifact in a registry is provided
func validateRepos(ctx context.Context, host string, opts []remote.Option, vr *types.ValidationResult) ([]string, error) {
	details := make([]string, 0)

	reg, err := name.NewRegistry(host)
	if err != nil {
		details = append(details, fmt.Sprintf("failed to get registry %v", host))
		return details, err
	}

	repoList, err := remote.Catalog(context.Background(), reg, opts...)
	if err != nil {
		details = append(details, fmt.Sprintf("failed to list repositories in registry %v", host))
		return details, err
	}

	if len(repoList) == 0 {
		details = append(details, fmt.Sprintf("no repositories found in registry %v", host))
		return details, nil
	}

	var repo name.Repository
	var ref name.Reference
	for _, curRepo := range repoList {
		repo, err = name.NewRepository(host + "/" + curRepo)
		if err != nil {
			details = append(details, fmt.Sprintf("failed to get repository %v/%v", host, curRepo))
			continue
		}

		tags, err := remote.List(repo, opts...)
		if err != nil {
			details = append(details, fmt.Sprintf("failed to get tags for repository %v/%v", host, curRepo))
			continue
		}

		if len(tags) == 0 {
			details = append(details, fmt.Sprintf("no tags found in repository %v/%v", host, curRepo))
			continue
		}

		tag := tags[0]
		ref, err = generateRef(host, fmt.Sprintf("%v:%v", curRepo, tag), vr)
		if err != nil {
			details = append(details, fmt.Sprintf("failed to generate reference for artifact %v/%v:%v", host, curRepo, tag))
			continue
		}
		break
	}
	if err != nil {
		return details, err
	}

	// download image without storing it on disk
	_, err = remote.Image(ref, opts...)
	if err != nil {
		details = append(details, fmt.Sprintf("failed to download artifact %v", ref.Name()))
		return details, err
	}

	return []string{}, nil
}

func generateRef(registry, artifact string, vr *types.ValidationResult) (name.Reference, error) {
	if strings.Contains(artifact, "@sha256:") {
		return name.NewDigest(fmt.Sprintf("%s/%s", registry, artifact))
	}

	if !strings.Contains(artifact, ":") {
		vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf("artifact %v does not contain a tag or digest, defaulting to \"latest\" tag", artifact))
	}

	return name.NewTag(fmt.Sprintf("%s/%s", registry, artifact))
}

func getEcrLoginToken(username, password, region string) (string, error) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(username, password, "")),
	)
	if err != nil {
		return "", fmt.Errorf("failed to load configuration, %v", err)
	}

	client := ecr.NewFromConfig(cfg)

	output, err := client.GetAuthorizationToken(context.TODO(), &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return "", fmt.Errorf("failed to get ECR authorization token: %v", err)
	}

	for _, data := range output.AuthorizationData {
		return *data.AuthorizationToken, nil
	}

	return "", fmt.Errorf("no authorization data available")
}

func parseEcrRegion(url string) (string, error) {
	parts := strings.Split(url, ".")
	if len(parts) != 6 || parts[2] != "ecr" {
		return "", errors.New(fmt.Sprintf("Invalid ECR URL %v", url))
	}

	region := parts[3]
	return region, nil
}

func setupAuthOpts(opts []remote.Option, registryName string, authentication v1alpha1.Auth) ([]remote.Option, error) {
	var auth authn.Authenticator

	if strings.Contains(registryName, "amazonaws.com") {
		region, err := parseEcrRegion(registryName)
		if err != nil {
			return nil, err
		}

		token, err := getEcrLoginToken(authentication.Username, authentication.Password, region)
		if err != nil {
			return nil, err
		}

		decodedToken, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			return nil, err
		}

		parts := strings.SplitN(string(decodedToken), ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed ECR login token: %v", token)
		}

		authentication.Username = parts[0]
		authentication.Password = parts[1]
	}

	if authentication.Username != "" && authentication.Password != "" {
		auth = &authn.Basic{Username: authentication.Username, Password: authentication.Password}
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
