package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/giantswarm/fulfillment/aws"
	"github.com/giantswarm/fulfillment/slack"
)

func Root(w http.ResponseWriter, r *http.Request, c aws.Service, s slack.Service) {
	switch r.Method {
	case http.MethodGet:
		rootGet(w, r)
	case http.MethodPost:
		rootPost(w, r, c, s)
	}
}

func rootGet(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token == "" {
		http.Error(w, "Missing or invalid token", http.StatusBadRequest)
		return
	}

	escapedToken := url.QueryEscape(token)

	err := Template.Execute(w, map[string]string{"Token": escapedToken})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func rootPost(w http.ResponseWriter, r *http.Request, c aws.Service, s slack.Service) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %s", err), http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	token := r.FormValue("token")

	unescapedToken, err := url.QueryUnescape(token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error unescaping token: %s", err), http.StatusBadRequest)
		return
	}

	// Note: AWS sends us an opaque token, but it's essentially a base64 string,
	// and + is a valid base64 character (see https://datatracker.ietf.org/doc/html/rfc4648#section-4),
	// but browsers will convert it to a space when POSTed as part of a form. So, we hackily convert it back here.
	unescapedToken = strings.ReplaceAll(unescapedToken, " ", "+")

	customerData, err := c.FetchCustomerData(unescapedToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error resolving customer: %s", err), http.StatusInternalServerError)
		return
	}

	customerData.Email = email

	if err := s.PostCustomerData(customerData); err != nil {
		http.Error(w, fmt.Sprintf("Error posting to Slack: %s", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/success", http.StatusSeeOther)
}
