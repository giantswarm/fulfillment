package handlers

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestWebhook(t *testing.T) {
	tests := []HTTPTest{
		{
			Test:         "Test that posting the webhook without anything else returns an error",
			Handler:      Webhook,
			Request:      MustRequest(http.NewRequest("POST", "/webhook", nil)),
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Test:    "Test that posting the webhook with a token as a form redirects to the root page",
			Handler: Webhook,
			Request: MustRequest(func() (*http.Request, error) {
				form := url.Values{}
				form.Add("x-amzn-marketplace-token", "foo")
				r := MustRequest(http.NewRequest("POST", "/webhook", strings.NewReader(form.Encode())))
				r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return r, nil
			}()),
			ExpectedCode:     http.StatusSeeOther,
			ExpectedLocation: "/?token=foo",
		},
	}

	RunTests(tests, t)
}
