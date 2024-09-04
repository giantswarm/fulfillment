package handlers

import (
	"net/http"
)

func Success(w http.ResponseWriter, r *http.Request) {
	err := Template.Execute(w, map[string]string{
		"Token":   "",
		"Success": "true",
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}
