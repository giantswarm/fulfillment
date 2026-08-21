package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/giantswarm/fulfillment/aws"
	"github.com/giantswarm/fulfillment/customer"
	"github.com/giantswarm/fulfillment/slack"
)

// recordingAWS captures the token that finally reaches the AWS lookup.
type recordingAWS struct {
	token string
}

func (m *recordingAWS) FetchCustomerData(token string) (customer.Data, error) {
	m.token = token
	return customer.Data{Identifier: "example-customer-identifier"}, nil
}

var hiddenTokenField = regexp.MustCompile(`<input type="hidden" name="token" value="([^"]*)">`)

// TestTokenSurvivesFullPipeline follows a token through every hop a real
// customer makes: the AWS webhook POST, the redirect to /, the form rendered
// there, and the form POST back to /. The token AWS sent us has to be the token
// we hand back to AWS, byte for byte.
func TestTokenSurvivesFullPipeline(t *testing.T) {
	tokens := []string{
		"simpletoken",
		// Real marketplace tokens are base64, so "+", "/" and "=" all occur.
		"abc+def/ghi==",
		"++//==",
	}

	for _, token := range tokens {
		// 1. AWS POSTs the token to /webhook, which redirects to /.
		form := url.Values{}
		form.Add("x-amzn-marketplace-token", token)
		webhookReq := MustRequest(http.NewRequest("POST", "/webhook", strings.NewReader(form.Encode())))
		webhookReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		webhookRec := httptest.NewRecorder()
		Webhook(webhookRec, webhookReq)

		location := webhookRec.Header().Get("Location")
		if webhookRec.Code != http.StatusSeeOther {
			t.Errorf("token %q: webhook returned %d, want %d", token, webhookRec.Code, http.StatusSeeOther)
			continue
		}

		// 2. The browser follows the redirect and gets the form.
		rootRec := httptest.NewRecorder()
		Root(rootRec, MustRequest(http.NewRequest("GET", location, nil)), &recordingAWS{}, &slack.Mock{})

		if rootRec.Code != http.StatusOK {
			t.Errorf("token %q: GET %q returned %d, want %d", token, location, rootRec.Code, http.StatusOK)
			continue
		}

		match := hiddenTokenField.FindStringSubmatch(rootRec.Body.String())
		if match == nil {
			t.Errorf("token %q: no hidden token field rendered at %q", token, location)
			continue
		}

		// 3. The browser POSTs the form back.
		postForm := url.Values{}
		postForm.Add("email", "customer@example.com")
		postForm.Add("token", match[1])
		postReq := MustRequest(http.NewRequest("POST", "/", strings.NewReader(postForm.Encode())))
		postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		awsService := &recordingAWS{}
		postRec := httptest.NewRecorder()
		Root(postRec, postReq, awsService, &slack.Mock{})

		if postRec.Code != http.StatusSeeOther {
			t.Errorf("token %q: form POST returned %d (%s), want %d",
				token, postRec.Code, strings.TrimSpace(postRec.Body.String()), http.StatusSeeOther)
			continue
		}

		if awsService.token != token {
			t.Errorf("token %q: AWS was queried with %q instead", token, awsService.token)
		}
	}
}

// compile-time check that the test double still satisfies the real interface.
var _ aws.Service = (*recordingAWS)(nil)
