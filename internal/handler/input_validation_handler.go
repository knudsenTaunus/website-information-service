package handler

import (
	"net/http"
	"regexp"
)

// InputValidationHandler HTTP middleware to validate the provided URL
func InputValidationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		website := r.URL.Query().Get("website")
		validationRegex := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b`)

		if !validationRegex.MatchString(website) {
			http.Error(w, "failed to validate URL", http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})

}
