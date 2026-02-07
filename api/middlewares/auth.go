package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/joyboy1210/stolight/models"
)

func CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("sto-Key")
		if key == "" {
			http.Error(w, "Missing API key", http.StatusUnauthorized)
			return
		}
		user, err := models.GetUserByKey(key)
		if err != nil {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		if user.Role != "admin" {
			if strings.HasPrefix(r.URL.Path, "/api/upload/") {

				bucketName := getNameFromURL(r.URL.Path)
				if bucketName != "" {
					if !isBucketAllowed(user.AllowedBuckets, bucketName) {
						http.Error(w, "Forbidden: Access to this bucket is not allowed", http.StatusForbidden)
						return
					}
				}
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func getNameFromURL(path string) string {
	cleanPath := strings.TrimPrefix(path, "/api/upload/")
	return cleanPath
}

func isBucketAllowed(allowed string, target string) bool {
	if allowed == "*" {
		return true
	}
	parts := strings.Split(allowed, ",")
	for _, p := range parts {
		if strings.TrimSpace(p) == target {
			return true
		}
	}
	return false
}
