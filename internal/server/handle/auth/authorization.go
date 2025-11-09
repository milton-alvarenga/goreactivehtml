package auth

import (
	"net/http"
	"strings"
)

func Check(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return false
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	return token != ""
}
