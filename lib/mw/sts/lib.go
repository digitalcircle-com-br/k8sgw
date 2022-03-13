package sts

import "net/http"

var sts = "*"

func Sts(next http.Handler) http.Handler {

	if sts == "" || sts == "*" {
		sts = "max-age=31536000; includeSubDomains"
	}
	//base.Log("Setting STS: %s", sts)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		(w).Header().Set("Strict-Transport-Security", sts)

		next.ServeHTTP(w, r)
	})
}
