package handlers

import (
	"fmt"
	"net/http"
	"net/url"
)

func Webhook(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Could not parse form: %s", err), http.StatusBadRequest)
		return
	}
	token := r.FormValue("x-amzn-marketplace-token")

	if token == "" {
		http.Error(w, "Missing or invalid token", http.StatusBadRequest)
		return
	}

	// Build the redirect target out of a fixed relative path plus a properly
	// encoded query, rather than interpolating the token into a string. The
	// caller-supplied token can then only ever be a query parameter value, and
	// it survives the round trip intact: AWS marketplace tokens are base64, and
	// an unencoded "+" would otherwise be read back as a space by the / handler.
	redirect := url.URL{
		Path:     "/",
		RawQuery: url.Values{"token": []string{token}}.Encode(),
	}

	http.Redirect(w, r, redirect.String(), http.StatusSeeOther)
}
