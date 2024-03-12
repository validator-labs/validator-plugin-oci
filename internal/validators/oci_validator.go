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
	"time"

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
	soci "github.com/spectrocloud-labs/validator-plugin-oci/internal/verifier"
	vapi "github.com/spectrocloud-labs/validator/api/v1alpha1"
	"github.com/spectrocloud-labs/validator/pkg/types"
	"github.com/spectrocloud-labs/validator/pkg/util"
)

const (
	verificationTimeout = 60 * time.Second
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
func (s *OciRuleService) ReconcileOciRegistryRule(rule v1alpha1.OciRegistryRule, username, password string, pubKeys [][]byte) (*types.ValidationRuleResult, error) {
	l := s.log.V(0).WithValues("rule", rule.Name(), "host", rule.Host)
	vr := buildValidationResult(rule)
	errs := make([]error, 0)
	details := make([]string, 0)

	opts, err := setupTransportOpts([]remote.Option{}, rule.CaCert)
	if err != nil {
		errs = append(errs, err)
		errMsg := "failed to setup http client transport"
		l.Error(err, errMsg, "caCert", rule.CaCert)
		s.updateResult(vr, errs, errMsg, rule.Name())
		return vr, err
	}

	ctx := context.Background()
	opts, err = setupAuthOpts(ctx, opts, rule.Host, username, password)
	if err != nil {
		errs = append(errs, err)
		errMsg := "failed to setup authentication"
		l.Error(err, errMsg, "auth", rule.Auth)
		s.updateResult(vr, errs, errMsg, rule.Name())
		return vr, err
	}

	if rule.SignatureVerification.SecretName != "" && len(pubKeys) == 0 {
		details = append(details, "no public keys provided for signature verification")
	}

	if len(rule.Artifacts) == 0 {
		errMsg := "failed to validate repositories in registry"
		d, e := validateRepos(ctx, rule.Host, opts, pubKeys, vr)
		details = append(details, d...)
		errs = append(errs, e...)

		if len(e) > 0 {
			l.Error(e[len(e)-1], errMsg)
			s.updateResult(vr, errs, errMsg, details...)
			return vr, errors.New(errMsg)
		}

		return vr, nil
	}

	errMsg := "failed to validate artifact in registry"
	for _, artifact := range rule.Artifacts {
		ref, err := generateRef(rule.Host, artifact.Ref, vr)
		if err != nil {
			details = append(details, fmt.Sprintf("failed to generate reference for artifact %s/%s", rule.Host, artifact.Ref))
			errs = append(errs, err)
			l.Error(err, errMsg, "artifact", artifact)
			continue
		}

		d, e := validateReference(ctx, ref, artifact.LayerValidation, pubKeys, opts)
		if len(e) > 0 {
			l.Error(e[len(e)-1], errMsg, "artifact", artifact)
		}
		details = append(details, d...)
		errs = append(errs, e...)
	}
	s.updateResult(vr, errs, errMsg, details...)

	if len(errs) > 0 {
		return vr, errors.New(errMsg)
	}
	return vr, nil
}

func generateRef(registry, artifact string, vr *types.ValidationRuleResult) (name.Reference, error) {
	if strings.Contains(artifact, "@sha256:") {
		return name.NewDigest(fmt.Sprintf("%s/%s", registry, artifact))
	}

	if !strings.Contains(artifact, ":") {
		vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf("artifact %s does not contain a tag or digest, defaulting to \"latest\" tag", artifact))
	}

	return name.NewTag(fmt.Sprintf("%s/%s", registry, artifact))
}

// validateArtifact validates a single artifact within a registry. This function is to be used when a path to a repo or an individual artifact is provided
func validateReference(ctx context.Context, ref name.Reference, fullLayerValidation bool, pubKeys [][]byte, opts []remote.Option) (details []string, errs []error) {
	details = make([]string, 0)
	errs = make([]error, 0)

	// validate artifact existence by issuing a HEAD request
	_, err := remote.Head(ref, opts...)
	if err != nil {
		details = append(details, fmt.Sprintf("failed to get descriptor for artifact %s", ref.Name()))
		errs = append(errs, err)
		return
	}

	// download image without storing it on disk
	img, err := remote.Image(ref, opts...)
	if err != nil {
		details = append(details, fmt.Sprintf("failed to download artifact %s", ref.Name()))
		errs = append(errs, err)
		return
	}

	var validateOpts []validate.Option
	if !fullLayerValidation {
		validateOpts = append(validateOpts, validate.Fast)
	}

	err = validate.Image(img, validateOpts...)
	if err != nil {
		details = append(details, fmt.Sprintf("failed validation for artifact %s", ref.Name()))
		errs = append(errs, err)
		return
	}

	if pubKeys != nil {
		details, errs = verifySignature(ctx, ref, pubKeys)
		if len(errs) > 0 {
			details = append(details, fmt.Sprintf("failed to verify signature for artifact %s", ref.Name()))
			return
		}
	}

	return
}

// validateRepos validates repos within a registry. This function is to be used when no particular artifact in a registry is provided
func validateRepos(ctx context.Context, host string, opts []remote.Option, pubKeys [][]byte, vr *types.ValidationRuleResult) (details []string, errs []error) {
	details = make([]string, 0)
	errs = make([]error, 0)

	reg, err := name.NewRegistry(host)
	if err != nil {
		details = append(details, fmt.Sprintf("failed to get registry %s", host))
		errs = append(errs, err)
		return
	}

	repoList, err := remote.Catalog(ctx, reg, opts...)
	if err != nil {
		details = append(details, fmt.Sprintf("failed to list repositories in registry %s", host))
		errs = append(errs, err)
		return
	}

	if len(repoList) == 0 {
		details = append(details, fmt.Sprintf("no repositories found in registry %s", host))
		return
	}

	var repo name.Repository
	var ref name.Reference
	foundArtifact := false

	for _, curRepo := range repoList {
		repo, err = name.NewRepository(host + "/" + curRepo)
		if err != nil {
			errs = append(errs, err)
			details = append(details, fmt.Sprintf("failed to get repository %s/%s", host, curRepo))
			continue
		}

		tags, err := remote.List(repo, opts...)
		if err != nil {
			errs = append(errs, err)
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
			errs = append(errs, err)
			details = append(details, fmt.Sprintf("failed to generate reference for artifact %s/%s:%s", host, curRepo, tag))
			continue
		}

		foundArtifact = true
		break
	}
	if !foundArtifact {
		return
	}

	details, errs = validateReference(ctx, ref, true, pubKeys, opts)
	return
}

func getEcrLoginToken(ctx context.Context, username, password, region string) (string, error) {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(username, password, "")),
	)
	if err != nil {
		return "", fmt.Errorf("failed to load configuration, %s", err)
	}

	client := ecr.NewFromConfig(cfg)

	output, err := client.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{})
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
		return "", fmt.Errorf("invalid ECR URL %s", url)
	}

	region := parts[3]
	return region, nil
}

func setupAuthOpts(ctx context.Context, opts []remote.Option, registryName, username, password string) ([]remote.Option, error) {
	var auth authn.Authenticator

	if strings.Contains(registryName, "amazonaws.com") {
		region, err := parseEcrRegion(registryName)
		if err != nil {
			return nil, err
		}

		token, err := getEcrLoginToken(ctx, username, password, region)
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

// verifySignature verifies the authenticity of the given image reference URL using the provided public keys.
func verifySignature(ctx context.Context, ref name.Reference, pubKeys [][]byte, opt ...remote.Option) (details []string, errs []error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, verificationTimeout)
	defer cancel()

	defaultCosignOciOpts := []soci.Options{
		soci.WithRemoteOptions(opt...),
	}

	details = make([]string, 0)
	errs = make([]error, 0)

	for _, key := range pubKeys {
		verifier, err := soci.NewCosignVerifier(ctxTimeout, append(defaultCosignOciOpts, soci.WithPublicKey(key))...)
		if err != nil {
			details = append(details, fmt.Sprintf("failed to create verifier with public key %s", key))
			errs = append(errs, err)
			return
		}

		hasValidSignature, err := verifier.Verify(ctxTimeout, ref)
		if err != nil {
			details = append(details, fmt.Sprintf("failed to verify signature of %s with public key %s", ref, key))
			errs = append(errs, err)
			continue
		}

		if hasValidSignature {
			details = nil
			errs = nil
			return
		}
	}

	details = append(details, fmt.Sprintf("no matching signatures were found for '%s'", ref))
	return
}

// buildValidationResult builds a default ValidationResult for a given validation type
func buildValidationResult(rule v1alpha1.OciRegistryRule) *types.ValidationRuleResult {
	state := vapi.ValidationSucceeded
	latestCondition := vapi.DefaultValidationCondition()
	latestCondition.Details = make([]string, 0)
	latestCondition.Failures = make([]string, 0)
	latestCondition.Message = fmt.Sprintf("All %s checks passed", constants.OciRegistry)
	latestCondition.ValidationRule = rule.Name()
	latestCondition.ValidationType = constants.OciRegistry
	return &types.ValidationRuleResult{Condition: &latestCondition, State: &state}
}

func (s *OciRuleService) updateResult(vr *types.ValidationRuleResult, errs []error, errMsg string, details ...string) {
	if len(errs) > 0 {
		vr.State = util.Ptr(vapi.ValidationFailed)
		vr.Condition.Message = errMsg
		for _, err := range errs {
			vr.Condition.Failures = append(vr.Condition.Failures, err.Error())
		}
	}
	vr.Condition.Details = append(vr.Condition.Details, details...)
}
