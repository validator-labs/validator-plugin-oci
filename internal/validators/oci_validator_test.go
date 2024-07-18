package validators

import (
	"context"
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-logr/logr"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/registry"
	"github.com/google/go-containerregistry/pkg/v1/random"
	"github.com/stretchr/testify/assert"
	"github.com/validator-labs/validator-plugin-oci/api/v1alpha1"
	"github.com/validator-labs/validator/pkg/oci"
	"github.com/validator-labs/validator/pkg/types"
)

const (
	registryName    = "registry"
	validArtifact   = "test/oci-image"
	invalidArtifact = "test/oci-image/does-not-exist"
)

var (
	vrr = buildValidationResult(v1alpha1.OciRegistryRule{})
)

func TestGenerateRef(t *testing.T) {
	type testCase struct {
		registry             string
		artifact             string
		validationRuleResult *types.ValidationRuleResult
		expectedRefName      string
		expectErr            bool
	}

	testCases := []testCase{
		{
			registry:             registryName,
			artifact:             "artifact@sha256:ddbac6e7732bf90a4e674a01bf279ce27ea8691530b8d209e6fe77477e0fa406",
			validationRuleResult: vrr,
			expectedRefName:      "registry/artifact@sha256:ddbac6e7732bf90a4e674a01bf279ce27ea8691530b8d209e6fe77477e0fa406",
			expectErr:            false,
		},
		{
			registry:             registryName,
			artifact:             "artifact:v1.0.0",
			validationRuleResult: vrr,
			expectedRefName:      "registry/artifact:v1.0.0",
			expectErr:            false,
		},
		{
			registry:             registryName,
			artifact:             "artifact",
			validationRuleResult: vrr,
			expectedRefName:      "registry/artifact:latest",
			expectErr:            false,
		},
		{
			registry:             registryName,
			artifact:             "invalidArtifact",
			validationRuleResult: vrr,
			expectedRefName:      "",
			expectErr:            true,
		},
	}

	for _, tc := range testCases {
		ref, err := generateRef(tc.registry, tc.artifact, tc.validationRuleResult)

		if tc.expectErr {
			assert.NotNil(t, err)
		} else {
			assert.Contains(t, ref.Name(), tc.expectedRefName)
			assert.NoError(t, err)
		}
	}
}

func TestValidateReference(t *testing.T) {
	s := httptest.NewServer(registry.New())
	defer s.Close()

	url, err := uploadArtifact(s, validArtifact)
	if err != nil {
		t.Fatal(err)
	}

	validRef, err := name.ParseReference(fmt.Sprintf("%s/%s", url.Host, validArtifact))
	if err != nil {
		t.Fatal(err)
	}

	invalidRef, err := name.ParseReference(fmt.Sprintf("%s/%s", url.Host, invalidArtifact))
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		ref             name.Reference
		layerValidation bool
		pubKeys         [][]byte
		expectedDetail  string
		expectErr       bool
	}

	testCases := []testCase{
		{
			ref:             validRef,
			layerValidation: false,
			pubKeys:         nil,
			expectedDetail:  "",
			expectErr:       false,
		},
		{
			ref:             validRef,
			layerValidation: true,
			pubKeys:         nil,
			expectedDetail:  "",
			expectErr:       false,
		},
		{
			ref:             invalidRef,
			layerValidation: false,
			pubKeys:         nil,
			expectedDetail:  "failed to get descriptor for artifact",
			expectErr:       true,
		},
		{
			ref:             validRef,
			layerValidation: true,
			pubKeys:         [][]byte{[]byte("invalid-pub-key-1"), []byte("invalid-pub-key-2")},
			expectedDetail:  "failed to create verifier with public key",
			expectErr:       true,
		},
		{
			ref:             validRef,
			layerValidation: true,
			pubKeys:         [][]byte{},
			expectedDetail:  "no matching signatures were found",
			expectErr:       true,
		},
		{
			ref:             validRef,
			layerValidation: true,
			pubKeys: [][]byte{
				[]byte(`-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEKPuCo9AmJCpqGWhefjbhkFcr1GA3
iNa765seE3jYC3MGUe5h52393Dhy7B5bXGsg6EfPpNYamlAEWjxCpHF3Lg==
-----END PUBLIC KEY-----`),
			},
			expectedDetail: "no matching signatures were found",
			expectErr:      true,
		},
	}

	for _, tc := range testCases {
		ociClient, err := oci.NewOCIClient(
			oci.WithAnonymousAuth(),
			oci.WithVerificationPublicKeys(tc.pubKeys),
		)
		if err != nil {
			t.Error(err)
		}
		ociService := NewOciRuleService(logr.Logger{}, WithOCIClient(ociClient))

		ctx := context.Background()
		details, errs := ociService.validateReference(ctx, tc.ref, tc.layerValidation, v1alpha1.SignatureVerification{
			SecretName: "secret",
		})

		if tc.expectedDetail == "" {
			assert.Empty(t, details)
		} else {
			assert.Contains(t, details[len(details)-1], tc.expectedDetail)
		}
		assert.Equal(t, tc.expectErr, len(errs) > 0)
	}
}

func TestValidateRepos(t *testing.T) {
	s1 := httptest.NewServer(registry.New())
	defer s1.Close()

	urlWithArtifact, err := uploadArtifact(s1, validArtifact)
	if err != nil {
		t.Fatal(err)
	}

	s2 := httptest.NewServer(registry.New())
	defer s1.Close()
	urlNoArtifact, err := url.Parse(s2.URL)
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		host           string
		expectedDetail string
		expectErr      bool
	}

	testCases := []testCase{
		{
			host:           urlWithArtifact.Host,
			expectedDetail: "",
			expectErr:      false,
		},
		{
			host:           urlNoArtifact.Host,
			expectedDetail: "no repositories found in registry",
			expectErr:      false,
		},
		{
			host:           "invalidHost",
			expectedDetail: "failed to list repositories in registry",
			expectErr:      true,
		},
	}

	for _, tc := range testCases {
		ociClient, err := oci.NewOCIClient(oci.WithAnonymousAuth())
		if err != nil {
			t.Error(err)
		}
		ociService := NewOciRuleService(logr.Logger{}, WithOCIClient(ociClient))

		rule := v1alpha1.OciRegistryRule{
			Host:                  tc.host,
			SignatureVerification: v1alpha1.SignatureVerification{},
		}
		details, errs := ociService.validateRepos(context.Background(), rule, &types.ValidationRuleResult{})

		if tc.expectedDetail == "" {
			assert.Empty(t, details)
		} else {
			assert.Len(t, details, 1)
			assert.Contains(t, details[0], tc.expectedDetail)
		}

		if !tc.expectErr {
			assert.Empty(t, errs)
		} else {
			assert.Len(t, errs, 1)
		}
	}
}

func TestReconcileOciRegistryRule(t *testing.T) {
	s1 := httptest.NewServer(registry.New())
	defer s1.Close()

	url, err := uploadArtifact(s1, validArtifact)
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		rule      v1alpha1.OciRegistryRule
		expectErr bool
	}

	testCases := []testCase{
		{
			rule: v1alpha1.OciRegistryRule{
				Host: url.Host,
			},
			expectErr: false,
		},
		{
			rule: v1alpha1.OciRegistryRule{
				Host: url.Host,
				Artifacts: []v1alpha1.Artifact{
					{
						Ref: validArtifact,
					},
				},
			},
			expectErr: false,
		},
		{
			rule: v1alpha1.OciRegistryRule{
				Host: url.Host,
				Artifacts: []v1alpha1.Artifact{
					{
						Ref: invalidArtifact,
					},
				},
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		l := logr.New(nil)
		s := NewOciRuleService(l)
		_, err := s.ReconcileOciRegistryRule(tc.rule)

		if tc.expectErr {
			assert.NotNil(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

// uploadArtifact uploads a random image to the registry at the given path.
func uploadArtifact(s *httptest.Server, path string) (*url.URL, error) {
	u, err := url.Parse(s.URL)
	if err != nil {
		return nil, err
	}

	img, err := random.Image(1024, 5)
	if err != nil {
		return nil, err
	}

	if err := crane.Push(img, fmt.Sprintf("%s/%s", u.Host, path)); err != nil {
		return nil, err
	}
	return u, nil
}
