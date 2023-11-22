package validators

/*
import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArtifact(t *testing.T) {
	type testCase struct {
		repoPath     string
		expectedPath string
		expectedRef  string
		expectedErr  error
	}

	testCases := []testCase{
		{
			repoPath:     "repo/artifact",
			expectedPath: "repo/artifact",
			expectedRef:  "",
			expectedErr:  nil,
		},
		{
			repoPath:     "repo/artifact@v1.1.1",
			expectedPath: "repo/artifact",
			expectedRef:  "v1.1.1",
			expectedErr:  nil,
		},
		{
			repoPath:     "repo/artifact@sha256:65ae8fd8713ede1977d26991821ba7eb3beb48ec575b31947568f30dbdd36863",
			expectedPath: "repo/artifact",
			expectedRef:  "sha256:65ae8fd8713ede1977d26991821ba7eb3beb48ec575b31947568f30dbdd36863",
			expectedErr:  nil,
		},
		{
			repoPath:     "repo/artifact@v1.1.1@invalid",
			expectedPath: "",
			expectedRef:  "",
			expectedErr:  errors.New("invalid artifact path"),
		},
	}

	for _, tc := range testCases {
		path, ref, err := parseArtifact(tc.repoPath)

		assert.Equal(t, tc.expectedPath, path)
		assert.Equal(t, tc.expectedRef, ref)
		if tc.expectedErr != nil {
			assert.EqualError(t, err, tc.expectedErr.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}
*/
