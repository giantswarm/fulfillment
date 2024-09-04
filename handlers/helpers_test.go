package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/giantswarm/fulfillment/aws"
	"github.com/giantswarm/fulfillment/slack"
)

type HTTPTest struct {
	Test string

	Handler       http.HandlerFunc
	MockedHandler func(*aws.Mock, *slack.Mock) func(http.ResponseWriter, *http.Request)
	Request       *http.Request

	ExpectedCode            int
	ExpectedLocation        string
	ExpectedMessages        []string
	ExpectedAWSMockCalled   bool
	ExpectedSlackMockCalled bool
}

func MustRequest(r *http.Request, err error) *http.Request {
	if err != nil {
		panic(err)
	}

	return r
}

func RunTests(tests []HTTPTest, t *testing.T) {
	for _, testCase := range tests {
		rr := httptest.NewRecorder()

		awsMock := &aws.Mock{}
		slackMock := &slack.Mock{}

		var handler http.HandlerFunc
		if testCase.Handler != nil {
			handler = testCase.Handler
		}
		if testCase.MockedHandler != nil {
			handler = testCase.MockedHandler(awsMock, slackMock)
		}

		handler.ServeHTTP(rr, testCase.Request)

		if testCase.ExpectedCode != 0 && testCase.ExpectedCode != rr.Code {
			t.Errorf("%v: Handler returned wrong status code: got %v want %v", testCase.Test, rr.Code, testCase.ExpectedCode)
		}

		if testCase.ExpectedLocation != "" && testCase.ExpectedLocation != rr.Header().Get("Location") {
			t.Errorf("%v: Handler returned wrong redirect location: got %v want %v", testCase.Test, rr.Header().Get("Location"), testCase.ExpectedLocation)
		}

		if len(testCase.ExpectedMessages) != 0 {
			for _, expectedMessage := range testCase.ExpectedMessages {
				if expectedMessage != "" && !strings.Contains(rr.Body.String(), expectedMessage) {
					t.Errorf("%v: Handler returned wrong message: got %v want %v", testCase.Test, rr.Body.String(), expectedMessage)
				}
			}
		}

		if testCase.MockedHandler != nil {
			if awsMock.Called && !testCase.ExpectedAWSMockCalled {
				t.Errorf("%v: Expected the mock AWS method: got %v want %v", testCase.Test, testCase.ExpectedAWSMockCalled, awsMock.Called)
			}
			if slackMock.Called && !testCase.ExpectedSlackMockCalled {
				t.Errorf("%v: Expected the mock Slack method: got %v want %v", testCase.Test, testCase.ExpectedSlackMockCalled, slackMock.Called)
			}
		}
	}
}
