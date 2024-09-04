package handlers

import (
	"net/http"
	"testing"
)

func TestSuccess(t *testing.T) {
	tests := []HTTPTest{
		{
			Test:         "Test that getting the success page returns an OK status and the email message",
			Handler:      Success,
			Request:      MustRequest(http.NewRequest("GET", "/success", nil)),
			ExpectedCode: http.StatusOK,
			ExpectedMessages: []string{
				`Please wait, you will receive an email from Giant Swarm within 48 hours from hello@giantswarm.io.`,
			},
		},
	}

	RunTests(tests, t)
}
