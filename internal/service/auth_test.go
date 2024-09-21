package service

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExtractToken(t *testing.T) {
	requestWithToken := httptest.NewRequest(http.MethodGet, "/", nil)
	requestWithToken.Header.Add("Authorization", "Bearer token")

	type testCase struct {
		name           string
		request        *http.Request
		expectedString string
		expectedBool   bool
	}

	testCases := []testCase{
		{
			name:           "Empty header",
			request:        httptest.NewRequest(http.MethodGet, "/", nil),
			expectedString: "",
			expectedBool:   false,
		},
		{
			name:           "Valid header",
			request:        requestWithToken,
			expectedString: "token",
			expectedBool:   true,
		},
	}

	for _, tt := range testCases {
		s, b := ExtractToken(tt.request)

		assert.Equal(t, s, tt.expectedString)
		assert.Equal(t, b, tt.expectedBool)
	}
}
