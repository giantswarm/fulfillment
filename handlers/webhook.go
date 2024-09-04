package handlers

import (
	"fmt"
	"net/http"
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

	url := fmt.Sprintf("/?token=%s", token)

	http.Redirect(w, r, url, http.StatusSeeOther)
}
