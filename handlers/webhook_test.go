package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func webhookRequest(token string) *http.Request {
	form := url.Values{}
	form.Add("x-amzn-marketplace-token", token)
	r := MustRequest(http.NewRequest("POST", "/webhook", strings.NewReader(form.Encode())))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return r
}

func TestWebhook(t *testing.T) {
	tests := []HTTPTest{
		{
			Test:         "Test that posting the webhook without anything else returns an error",
			Handler:      Webhook,
			Request:      MustRequest(http.NewRequest("POST", "/webhook", nil)),
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Test:             "Test that posting the webhook with a token as a form redirects to the root page",
			Handler:          Webhook,
			Request:          webhookRequest("foo"),
			ExpectedCode:     http.StatusSeeOther,
			ExpectedLocation: "/?token=foo",
		},
		{
			Test:             "Test that a base64 token keeps its + and = when redirected",
			Handler:          Webhook,
			Request:          webhookRequest("a+b/c=="),
			ExpectedCode:     http.StatusSeeOther,
			ExpectedLocation: "/?token=a%2Bb%2Fc%3D%3D",
		},
		{
			Test:             "Test that a token cannot inject extra query parameters or a fragment",
			Handler:          Webhook,
			Request:          webhookRequest("foo&admin=1#bar"),
			ExpectedCode:     http.StatusSeeOther,
			ExpectedLocation: "/?token=foo%26admin%3D1%23bar",
		},
		{
			Test:             "Test that a token cannot redirect off-site",
			Handler:          Webhook,
			Request:          webhookRequest("//evil.example.com"),
			ExpectedCode:     http.StatusSeeOther,
			ExpectedLocation: "/?token=%2F%2Fevil.example.com",
		},
	}

	RunTests(tests, t)
}

// TestWebhookRedirectStaysRelative asserts the redirect target can never leave
// the current origin, whatever the caller puts in the token.
func TestWebhookRedirectStaysRelative(t *testing.T) {
	tokens := []string{
		"https://evil.example.com",
		"//evil.example.com",
		"/\\evil.example.com",
		"\r\nLocation: https://evil.example.com",
		"?x=1",
		"#frag",
	}

	for _, token := range tokens {
		r := webhookRequest(token)
		rr := httptest.NewRecorder()
		Webhook(rr, r)

		location := rr.Header().Get("Location")
		u, err := url.Parse(location)
		if err != nil {
			t.Errorf("token %q: Location %q does not parse: %s", token, location, err)
			continue
		}
		if u.Scheme != "" || u.Host != "" {
			t.Errorf("token %q: redirect left the origin: %q", token, location)
		}
		if u.Path != "/" {
			t.Errorf("token %q: redirect changed the path: %q", token, u.Path)
		}
		if got := u.Query().Get("token"); got != token {
			t.Errorf("token %q: round-tripped as %q", token, got)
		}
	}
}
