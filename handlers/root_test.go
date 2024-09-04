package handlers

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/giantswarm/fulfillment/aws"
	"github.com/giantswarm/fulfillment/slack"
)

func getMockRootHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Root(w, r, &aws.Mock{}, &slack.Mock{})
	}
}

func TestRoot(t *testing.T) {
	tests := []HTTPTest{
		{
			Test:         "Test that getting the root URL without a token returns an error",
			Handler:      getMockRootHandler(),
			Request:      MustRequest(http.NewRequest("GET", "/", nil)),
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Test:         "Test that getting the root URL with a token presents a form",
			Handler:      getMockRootHandler(),
			Request:      MustRequest(http.NewRequest("GET", "/?token=foo", nil)),
			ExpectedCode: http.StatusOK,
			ExpectedMessages: []string{
				`Please enter your email address below, and Giant Swarm will get in touch within 48 hours to discuss fulfillment of your Giant Swarm AWS SaaS offering.`,
				`<input type="hidden" name="token" value="foo">`,
			},
		},
		{
			Test: "Test that posting the root url with a form calls the AWS ResolveCustomer endpoint, posts an update to Slack, and redirects to the success page",
			MockedHandler: func(awsMock *aws.Mock, slackMock *slack.Mock) func(http.ResponseWriter, *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					Root(w, r, awsMock, slackMock)
				}
			},
			Request: MustRequest(func() (*http.Request, error) {
				form := url.Values{}
				form.Add("email", "user@example.com")
				form.Add("token", "foo")
				r := MustRequest(http.NewRequest("POST", "/", strings.NewReader(form.Encode())))
				r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return r, nil
			}()),
			ExpectedCode:            http.StatusSeeOther,
			ExpectedLocation:        "/success",
			ExpectedAWSMockCalled:   true,
			ExpectedSlackMockCalled: true,
		},
	}

	RunTests(tests, t)
}
