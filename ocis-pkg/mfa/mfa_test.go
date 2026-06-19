package mfa_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/mfa"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/test-go/testify/require"
)

func exampleUsage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// In a central place of your service enhance request once.
		// Note: This will not overwrite existing context values so it's safe (but unnecessary) to call multiple times.
		if mfa.IsMFAHeaderTrue(r) { // this is just a condition to set the MFA status. There could be other conditions.
			r = r.WithContext(revactx.SetMFA(r.Context()))
		}

		// somewhere in your code extract the context
		ctx := r.Context()

		// now you can check if the user has MFA enabled
		if !revactx.HasMFA(ctx) {
			// use this line to log access denied information
			// mfa package will not log anything by itself
			mfa.SetRequiredStatus(w)
			return
		}
		// user has MFA enabled, you can now proceed with sensitive operation
	}
}

func TestMFALifecycle(t *testing.T) {
	testCases := []struct {
		Alias         string
		HasMFA        bool
		ShouldHaveMFA bool
		ResponseCode  int
	}{
		{
			Alias:        "simple",
			HasMFA:       true,
			ResponseCode: http.StatusOK,
		},
		{
			Alias:        "denied",
			HasMFA:       false,
			ResponseCode: http.StatusForbidden,
		},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://url&method.doesnt.matter", nil)
		mfa.SetHeader(r, tc.HasMFA)

		exampleUsage().ServeHTTP(w, r)
		res := w.Result()

		require.Equal(t, tc.ResponseCode, res.StatusCode, tc.Alias)
		if tc.ResponseCode == http.StatusForbidden {
			require.Equal(t, "true", res.Header.Get(mfa.MFARequiredHeader), tc.Alias)
		} else {
			require.Empty(t, res.Header.Get(mfa.MFARequiredHeader), tc.Alias)
		}

	}

}
