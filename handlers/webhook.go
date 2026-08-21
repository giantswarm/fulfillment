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

	target := "/?" + url.Values{"token": {token}}.Encode()

	http.Redirect(w, r, target, http.StatusSeeOther)
}
