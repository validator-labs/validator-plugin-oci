package validators

import (
	"context"
	"fmt"
	"net"
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
	"github.com/validator-labs/validator/pkg/types"

	"github.com/validator-labs/validator-plugin-oci/pkg/oci"
)

const (
	registryName    = "registry"
	validArtifact   = "test/oci-image"
	invalidArtifact = "test/oci-image/does-not-exist"
)

var (
	vrr = BuildValidationResult(v1alpha1.OciRegistryRule{})
)

func TestGenerateRef(t *testing.T) {
	testCases := []struct {
		name                 string
		registry             string
		artifact             string
		validationRuleResult *types.ValidationRuleResult
		expectedRefName      string
		expectErr            bool
	}{
		{
			name:                 "Pass: valid artifact with SHA",
			registry:             registryName,
			artifact:             "artifact@sha256:ddbac6e7732bf90a4e674a01bf279ce27ea8691530b8d209e6fe77477e0fa406",
			validationRuleResult: vrr,
			expectedRefName:      "registry/artifact@sha256:ddbac6e7732bf90a4e674a01bf279ce27ea8691530b8d209e6fe77477e0fa406",
			expectErr:            false,
		},
		{
			name:                 "Pass: valid artifact with semver tag",
			registry:             registryName,
			artifact:             "artifact:v1.0.0",
			validationRuleResult: vrr,
			expectedRefName:      "registry/artifact:v1.0.0",
			expectErr:            false,
		},
		{
			name:                 "Pass: valid artifact with latest tag",
			registry:             registryName,
			artifact:             "artifact",
			validationRuleResult: vrr,
			expectedRefName:      "registry/artifact:latest",
			expectErr:            false,
		},
		{
			name:                 "Fail: invalid artifact",
			registry:             registryName,
			artifact:             "invalidArtifact",
			validationRuleResult: vrr,
			expectedRefName:      "",
			expectErr:            true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ref, err := generateRef(tc.registry, tc.artifact, tc.validationRuleResult)

			if tc.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Contains(t, ref.Name(), tc.expectedRefName)
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateReference(t *testing.T) {
	s := httptest.NewServer(registry.New())
	defer s.Close()
	port := s.Listener.Addr().(*net.TCPAddr).Port

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

	testCases := []struct {
		name            string
		ref             name.Reference
		layerValidation bool
		pubKeys         [][]byte
		sv              v1alpha1.SignatureVerification
		expectedDetail  string
		expectErr       bool
	}{
		{
			name:            "Pass: valid ref, no layer validation",
			ref:             validRef,
			layerValidation: false,
			pubKeys:         nil,
			expectedDetail:  fmt.Sprintf("pulled and validated artifact 127.0.0.1:%d/test/oci-image:latest", port),
			expectErr:       false,
		},
		{
			name:            "Pass: valid ref, layer validation",
			ref:             validRef,
			layerValidation: true,
			pubKeys:         nil,
			expectedDetail:  fmt.Sprintf("pulled and validated artifact 127.0.0.1:%d/test/oci-image:latest", port),
			expectErr:       false,
		},
		{
			name:            "Fail: invalid ref",
			ref:             invalidRef,
			layerValidation: false,
			pubKeys:         nil,
			expectedDetail:  "failed to get descriptor for artifact",
			expectErr:       true,
		},
		{
			name:            "Fail: valid ref, signature verification enabled with invalid keys",
			ref:             validRef,
			layerValidation: true,
			pubKeys:         [][]byte{[]byte("invalid-pub-key-1"), []byte("invalid-pub-key-2")},
			sv: v1alpha1.SignatureVerification{
				SecretName: "secret",
			},
			expectedDetail: "failed to create verifier with public key",
			expectErr:      true,
		},
		{
			name:            "Fail: valid ref, signature verification enabled with no keys",
			ref:             validRef,
			layerValidation: true,
			pubKeys:         [][]byte{},
			sv: v1alpha1.SignatureVerification{
				SecretName: "secret",
			},
			expectedDetail: fmt.Sprintf("pulled and validated artifact 127.0.0.1:%d/test/oci-image:latest", port),
			expectErr:      true,
		},
		{
			name:            "Fail: valid ref, signature verification enabled with invalid public key",
			ref:             validRef,
			layerValidation: true,
			pubKeys: [][]byte{
				[]byte(`-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEKPuCo9AmJCpqGWhefjbhkFcr1GA3
iNa765seE3jYC3MGUe5h52393Dhy7B5bXGsg6EfPpNYamlAEWjxCpHF3Lg==
-----END PUBLIC KEY-----`),
			},
			sv: v1alpha1.SignatureVerification{
				SecretName: "secret",
			},
			expectedDetail: "no matching signatures were found",
			expectErr:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ociClient, err := oci.NewOCIClient(
				oci.WithAnonymousAuth(),
				oci.WithVerificationPublicKeys(tc.pubKeys),
			)
			if err != nil {
				t.Error(err)
			}
			ociService := NewOciRuleService(logr.Logger{}, WithOCIClient(ociClient))

			ctx := context.Background()
			details, errs := ociService.validateReference(ctx, tc.ref, tc.layerValidation, tc.sv)

			if tc.expectedDetail == "" {
				assert.Empty(t, details)
			} else {
				assert.Contains(t, details[len(details)-1], tc.expectedDetail)
			}
			assert.Equal(t, tc.expectErr, len(errs) > 0)
		})
	}
}

func TestValidateRepos(t *testing.T) {
	s1 := httptest.NewServer(registry.New())
	defer s1.Close()
	port := s1.Listener.Addr().(*net.TCPAddr).Port

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

	testCases := []struct {
		name           string
		host           string
		expectedDetail string
		expectErr      bool
	}{
		{
			name:           "Pass: valid registry",
			host:           urlWithArtifact.Host,
			expectedDetail: fmt.Sprintf("pulled and validated artifact 127.0.0.1:%d/test/oci-image:latest", port),
			expectErr:      false,
		},
		{
			name:           "Fail: invalid registry: no repositories",
			host:           urlNoArtifact.Host,
			expectedDetail: "no repositories found in registry",
			expectErr:      false,
		},
		{
			name:           "Fail: invalid registry",
			host:           "invalidHost",
			expectedDetail: "failed to list repositories in registry",
			expectErr:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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
		})
	}
}

func TestReconcileOciRegistryRule(t *testing.T) {
	s1 := httptest.NewServer(registry.New())
	defer s1.Close()

	url, err := uploadArtifact(s1, validArtifact)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name      string
		rule      v1alpha1.OciRegistryRule
		expectErr bool
	}{
		{
			name: "Pass: valid host, no artifacts",
			rule: v1alpha1.OciRegistryRule{
				Host: url.Host,
			},
			expectErr: false,
		},
		{
			name: "Pass: valid host with artifacts",
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
			name: "Fail: valid host, invalid artifact",
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
		t.Run(tc.name, func(t *testing.T) {
			ociClient, err := oci.NewOCIClient(oci.WithAnonymousAuth())
			if err != nil {
				t.Error(err)
			}
			ociService := NewOciRuleService(logr.Logger{}, WithOCIClient(ociClient))

			_, err = ociService.ReconcileOciRegistryRule(tc.rule)

			if tc.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
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
