package validators

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	validURL  = "745150053801.dkr.ecr.us-east-1.amazonaws.com"
	longURL   = "745150053801.dkr.ecr.us-east-1.amazonaws.com.invalid"
	shortURL  = "dkr.ecr.us-east-1.amazonaws.com"
	notEcrURL = "745150053801.dkr.notEcr.us-east-1.amazonaws.com"
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
			expectedErr:    errors.New(fmt.Sprintf("Invalid ecr url %v", longURL)),
		},
		{
			URL:            longURL,
			expectedRegion: "",
			expectedErr:    errors.New(fmt.Sprintf("Invalid ecr url %v", longURL)),
		},
		{
			URL:            notEcrURL,
			expectedRegion: "",
			expectedErr:    errors.New(fmt.Sprintf("Invalid ecr url %v", notEcrURL)),
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
