package validators

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/spectrocloud-labs/validator-plugin-oci/api/v1alpha1"
	"github.com/spectrocloud-labs/validator/pkg/types"
	"github.com/stretchr/testify/assert"
)

const (
	validURL  = "745150053801.dkr.ecr.us-east-1.amazonaws.com"
	longURL   = "745150053801.dkr.ecr.us-east-1.amazonaws.com.invalid"
	shortURL  = "dkr.ecr.us-east-1.amazonaws.com"
	notEcrURL = "745150053801.dkr.notEcr.us-east-1.amazonaws.com"

	registry = "registry"

	caCert = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURGRENDQWZ5Z0F3SUJBZ0lSQU1IM2F5UWUvYUkwK1Y0OGE0QnlNYVF3RFFZSktvWklodmNOQVFFTEJRQXcKRkRFU01CQUdBMVVFQXhNSmFHRnlZbTl5TFdOaE1CNFhEVEl6TURneE9UQXdNRGMwTjFvWERUSTBNRGd4T0RBdwpNRGMwTjFvd0ZERVNNQkFHQTFVRUF4TUphR0Z5WW05eUxXTmhNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DCkFROEFNSUlCQ2dLQ0FRRUFtSTlaOUp4Zmg1SlZ0REMyS1U3WS85K0srSzFGanI1T3ZOUmk1RE1iR3NpZ2t6N2wKR054ZkgwSWdJRFdPZ1Q5L3YvNTJ5N1NZcnNrYWJYRVR1TEs3ajlaTXdXck9ZZm1mckcva1VMK3FlTThPYjZZdQorSUhNV3E4Z3VOdzJ2UW9yK214eW1JRUFTc3ZsTDBzd25vSXVQWk1GbFg5NEpWNUJtR3BtVjFrNmZaSVh2b05nClVUaHFoSE4vUFVIVDNibkxYaGlTdFNCZjBIMFR1U3BLMitEVXpvOFVRdlNvaStyV0k5SXRRRENZemtrWjg0bjIKeEp6WCtHSXlvYjNsdGdXU3ZSYmRURU9VK1pmYm0xVTRMV1U4YjdhVWRZSVdwM1EzSEVZK2F1WG1SbmlRSld2aQpQVUJrNTBUQnVPNFFJSWx0VGtHS3VTM0svR2s2SU0ra2FibUY0d0lEQVFBQm8yRXdYekFPQmdOVkhROEJBZjhFCkJBTUNBcVF3SFFZRFZSMGxCQll3RkFZSUt3WUJCUVVIQXdFR0NDc0dBUVVGQndNQ01BOEdBMVVkRXdFQi93UUYKTUFNQkFmOHdIUVlEVlIwT0JCWUVGRk4vYkhTS256ZE9IZ0k4d2ttNlpPbnV0eTRxTUEwR0NTcUdTSWIzRFFFQgpDd1VBQTRJQkFRQXRxNk9vRDI2NWF4Y2x3QVg3ZzdTdEtiZFNkeVNNcC9GbEJZOEJTS0QzdUxDWUtJZmRMdnJJClhKa0Z6MUFXa3hLb1dDbyt2RFl2cEUybE42WXAvakRQZUhZd1c3WG1HQTZJZDRVZ2FtdzV2NHhVZXg5Wis0V1IKbzdqNnV1NkVYK0xOdkQzREFSOFk4aEN3S1NDV3JNWURGbWV3Wmh6N05kY1VBcEp5M3phWTZWeHMvS3dlTGxicwpwbHh2TjlIWCtocVZobC8rWkFtbFZOOVZmZkhHblpsZm5tZW5Tb3RSbjJnR3Rmc0VrV3dhR3UvOUNPbTNQZlhTCjNTY0NGZTNNSjBZbjYvcG1iQkFVVnRtRjFUOTNsT2FYZ3VIek1pWEhJdyt4NUhadnhidkRQbmZ0Z0tnQWpWWU0KRmY0ODlRb28yalVuRVNmK2JRZFczcnpjMUFaMndwbmgKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="
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
			URL:            validURL,
			expectedRegion: "us-east-1",
			expectedErr:    nil,
		},
		{
			URL:            longURL,
			expectedRegion: "",
			expectedErr:    errors.New(fmt.Sprintf("Invalid ECR URL %s", longURL)),
		},
		{
			URL:            longURL,
			expectedRegion: "",
			expectedErr:    errors.New(fmt.Sprintf("Invalid ECR URL %s", longURL)),
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
			registry:         registry,
			artifact:         "artifact@sha256:ddbac6e7732bf90a4e674a01bf279ce27ea8691530b8d209e6fe77477e0fa406",
			validationResult: vr,
			expectedRefName:  "registry/artifact@sha256:ddbac6e7732bf90a4e674a01bf279ce27ea8691530b8d209e6fe77477e0fa406",
			expectErr:        false,
		},
		{
			registry:         registry,
			artifact:         "artifact:v1.0.0",
			validationResult: vr,
			expectedRefName:  "registry/artifact:v1.0.0",
			expectErr:        false,
		},
		{
			registry:         registry,
			artifact:         "artifact",
			validationResult: vr,
			expectedRefName:  "registry/artifact:latest",
			expectErr:        false,
		},
		{
			registry:         registry,
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
	// create table test for the following function
	// func setupTransportOpts(opts []remote.Option, caCert string) ([]remote.Option, error)

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
			expectedOpts: []remote.Option{},
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
