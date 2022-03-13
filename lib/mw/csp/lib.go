package csp

import "net/http"

var csp = "*"
func Csp(next http.Handler) http.Handler {
	if csp == "" || csp == "*" {
		csp = "default-src: 'self';"
	}
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Content-Security-Policy:", csp)
		next.ServeHTTP(w, r)
	})
}