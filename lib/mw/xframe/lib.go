package xframe

import "net/http"

var xframe = "*"

func XFrame(next http.Handler) http.Handler {

	if xframe == "" || xframe == "*" {
		xframe = "SAMEORIGIN"
	}
	//base.Log("Setting XFrame: %s", xframe)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		(w).Header().Set("X-Frame-Options:", xframe)

		next.ServeHTTP(w, r)
	})
}
