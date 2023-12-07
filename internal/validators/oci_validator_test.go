package validators

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/registry"
	"github.com/google/go-containerregistry/pkg/v1/random"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/spectrocloud-labs/validator-plugin-oci/api/v1alpha1"
	"github.com/spectrocloud-labs/validator/pkg/types"
	"github.com/stretchr/testify/assert"
)

const (
	ecrRegistry     = "745150053801.dkr.ecr.us-east-1.amazonaws.com"
	longURL         = "745150053801.dkr.ecr.us-east-1.amazonaws.com.invalid"
	shortURL        = "dkr.ecr.us-east-1.amazonaws.com"
	notEcrURL       = "745150053801.dkr.notEcr.us-east-1.amazonaws.com"
	registryName    = "registry"
	caCert          = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURGRENDQWZ5Z0F3SUJBZ0lSQU1IM2F5UWUvYUkwK1Y0OGE0QnlNYVF3RFFZSktvWklodmNOQVFFTEJRQXcKRkRFU01CQUdBMVVFQXhNSmFHRnlZbTl5TFdOaE1CNFhEVEl6TURneE9UQXdNRGMwTjFvWERUSTBNRGd4T0RBdwpNRGMwTjFvd0ZERVNNQkFHQTFVRUF4TUphR0Z5WW05eUxXTmhNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DCkFROEFNSUlCQ2dLQ0FRRUFtSTlaOUp4Zmg1SlZ0REMyS1U3WS85K0srSzFGanI1T3ZOUmk1RE1iR3NpZ2t6N2wKR054ZkgwSWdJRFdPZ1Q5L3YvNTJ5N1NZcnNrYWJYRVR1TEs3ajlaTXdXck9ZZm1mckcva1VMK3FlTThPYjZZdQorSUhNV3E4Z3VOdzJ2UW9yK214eW1JRUFTc3ZsTDBzd25vSXVQWk1GbFg5NEpWNUJtR3BtVjFrNmZaSVh2b05nClVUaHFoSE4vUFVIVDNibkxYaGlTdFNCZjBIMFR1U3BLMitEVXpvOFVRdlNvaStyV0k5SXRRRENZemtrWjg0bjIKeEp6WCtHSXlvYjNsdGdXU3ZSYmRURU9VK1pmYm0xVTRMV1U4YjdhVWRZSVdwM1EzSEVZK2F1WG1SbmlRSld2aQpQVUJrNTBUQnVPNFFJSWx0VGtHS3VTM0svR2s2SU0ra2FibUY0d0lEQVFBQm8yRXdYekFPQmdOVkhROEJBZjhFCkJBTUNBcVF3SFFZRFZSMGxCQll3RkFZSUt3WUJCUVVIQXdFR0NDc0dBUVVGQndNQ01BOEdBMVVkRXdFQi93UUYKTUFNQkFmOHdIUVlEVlIwT0JCWUVGRk4vYkhTS256ZE9IZ0k4d2ttNlpPbnV0eTRxTUEwR0NTcUdTSWIzRFFFQgpDd1VBQTRJQkFRQXRxNk9vRDI2NWF4Y2x3QVg3ZzdTdEtiZFNkeVNNcC9GbEJZOEJTS0QzdUxDWUtJZmRMdnJJClhKa0Z6MUFXa3hLb1dDbyt2RFl2cEUybE42WXAvakRQZUhZd1c3WG1HQTZJZDRVZ2FtdzV2NHhVZXg5Wis0V1IKbzdqNnV1NkVYK0xOdkQzREFSOFk4aEN3S1NDV3JNWURGbWV3Wmh6N05kY1VBcEp5M3phWTZWeHMvS3dlTGxicwpwbHh2TjlIWCtocVZobC8rWkFtbFZOOVZmZkhHblpsZm5tZW5Tb3RSbjJnR3Rmc0VrV3dhR3UvOUNPbTNQZlhTCjNTY0NGZTNNSjBZbjYvcG1iQkFVVnRtRjFUOTNsT2FYZ3VIek1pWEhJdyt4NUhadnhidkRQbmZ0Z0tnQWpWWU0KRmY0ODlRb28yalVuRVNmK2JRZFczcnpjMUFaMndwbmgKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="
	username        = "user"
	password        = "pa$$w0rd"
	validArtifact   = "test/oci-image"
	invalidArtifact = "test/oci-image/does-not-exist"
)

var (
	vr = buildValidationResult(v1alpha1.OciRegistryRule{})
)

func TestParseEcrRegion(t *testing.T) {
	type testCase struct {
		URL            string
		expectedRegion string
		expectedErr    error
	}

	testCases := []testCase{
		{
			URL:            ecrRegistry,
			expectedRegion: "us-east-1",
			expectedErr:    nil,
		},
		{
			URL:            longURL,
			expectedRegion: "",
			expectedErr:    errors.New(fmt.Sprintf("Invalid ECR URL %s", longURL)),
		},
		{
			URL:            shortURL,
			expectedRegion: "",
			expectedErr:    errors.New(fmt.Sprintf("Invalid ECR URL %s", shortURL)),
		},
		{
			URL:            notEcrURL,
			expectedRegion: "",
			expectedErr:    errors.New(fmt.Sprintf("Invalid ECR URL %s", notEcrURL)),
		},
	}

	for _, tc := range testCases {
		region, err := parseEcrRegion(tc.URL)

		assert.Equal(t, tc.expectedRegion, region)
		if tc.expectedErr != nil {
			assert.EqualError(t, err, tc.expectedErr.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestGenerateRef(t *testing.T) {
	type testCase struct {
		registry         string
		artifact         string
		validationResult *types.ValidationResult
		expectedRefName  string
		expectErr        bool
	}

	testCases := []testCase{
		{
			registry:         registryName,
			artifact:         "artifact@sha256:ddbac6e7732bf90a4e674a01bf279ce27ea8691530b8d209e6fe77477e0fa406",
			validationResult: vr,
			expectedRefName:  "registry/artifact@sha256:ddbac6e7732bf90a4e674a01bf279ce27ea8691530b8d209e6fe77477e0fa406",
			expectErr:        false,
		},
		{
			registry:         registryName,
			artifact:         "artifact:v1.0.0",
			validationResult: vr,
			expectedRefName:  "registry/artifact:v1.0.0",
			expectErr:        false,
		},
		{
			registry:         registryName,
			artifact:         "artifact",
			validationResult: vr,
			expectedRefName:  "registry/artifact:latest",
			expectErr:        false,
		},
		{
			registry:         registryName,
			artifact:         "invalidArtifact",
			validationResult: vr,
			expectedRefName:  "",
			expectErr:        true,
		},
	}

	for _, tc := range testCases {
		ref, err := generateRef(tc.registry, tc.artifact, tc.validationResult)

		if tc.expectErr {
			assert.NotNil(t, err)
		} else {
			assert.Contains(t, ref.Name(), tc.expectedRefName)
			assert.NoError(t, err)
		}
	}
}

func TestSetupTransportOpts(t *testing.T) {
	type testCase struct {
		inputOpts    []remote.Option
		caCert       string
		expectedOpts []remote.Option
		expectErr    bool
	}

	testCases := []testCase{
		{
			inputOpts:    []remote.Option{},
			caCert:       "",
			expectedOpts: []remote.Option{},
			expectErr:    false,
		},
		{
			inputOpts:    []remote.Option{},
			caCert:       caCert,
			expectedOpts: []remote.Option{remote.WithTransport(&http.Transport{})},
			expectErr:    false,
		},
		{
			inputOpts:    []remote.Option{},
			caCert:       "invalid cert",
			expectedOpts: nil,
			expectErr:    true,
		},
	}

	for _, tc := range testCases {
		opts, err := setupTransportOpts(tc.inputOpts, tc.caCert)

		if tc.expectErr {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, len(tc.expectedOpts), len(opts))
			assert.NoError(t, err)
		}
	}
}

func TestSetupAuthOpts(t *testing.T) {
	type testCase struct {
		inputOpts    []remote.Option
		registryName string
		username     string
		password     string
		expectedOpts []remote.Option
		expectErr    bool
	}

	testCases := []testCase{
		{
			inputOpts:    []remote.Option{},
			registryName: registryName,
			username:     "",
			password:     "",
			expectedOpts: []remote.Option{remote.WithAuth(authn.Anonymous)},
			expectErr:    false,
		},
		{
			inputOpts:    []remote.Option{},
			registryName: registryName,
			username:     username,
			password:     password,
			expectedOpts: []remote.Option{remote.WithAuth(&authn.Basic{Username: username, Password: password})},
			expectErr:    false,
		},
		{
			inputOpts:    []remote.Option{},
			registryName: ecrRegistry,
			username:     "",
			password:     "",
			expectedOpts: nil,
			expectErr:    true,
		},
	}

	for _, tc := range testCases {
		opts, err := setupAuthOpts(tc.inputOpts, tc.registryName, tc.username, tc.password)

		if tc.expectErr {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, len(tc.expectedOpts), len(opts))
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
		ref            name.Reference
		download       bool
		expectedDetail string
		expectErr      bool
	}

	testCases := []testCase{
		{
			ref:            validRef,
			download:       false,
			expectedDetail: "",
			expectErr:      false,
		},
		{
			ref:            validRef,
			download:       true,
			expectedDetail: "",
			expectErr:      false,
		},
		{
			ref:            invalidRef,
			download:       false,
			expectedDetail: "failed to get descriptor for artifact",
			expectErr:      true,
		},
	}

	for _, tc := range testCases {
		detail, err := validateReference(tc.ref, tc.download, []remote.Option{remote.WithAuth(authn.Anonymous)})

		if tc.expectErr {
			assert.NotNil(t, err)
			assert.Contains(t, detail, tc.expectedDetail)
		} else {
			assert.Empty(t, detail)
			assert.NoError(t, err)
		}
	}
}

func TestValidateRepos(t *testing.T) {
	s1 := httptest.NewServer(registry.New())
	defer s1.Close()

	hostWithArtifact, err := uploadArtifact(s1, validArtifact)
	if err != nil {
		t.Fatal(err)
	}

	s2 := httptest.NewServer(registry.New())
	defer s1.Close()
	hostNoArtifact, err := url.Parse(s2.URL)
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
			host:           hostWithArtifact.Host,
			expectedDetail: "",
			expectErr:      false,
		},
		{
			host:           hostNoArtifact.Host,
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
		details, err := validateRepos(context.Background(), tc.host, []remote.Option{remote.WithAuth(authn.Anonymous)}, &types.ValidationResult{})

		if tc.expectedDetail == "" {
			assert.Empty(t, details)
		} else {
			assert.Len(t, details, 1)
			assert.Contains(t, details[0], tc.expectedDetail)
		}

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
