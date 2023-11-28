package validators

import (
	"errors"
	"fmt"
	"testing"

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
