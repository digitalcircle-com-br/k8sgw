package cors

import (
	"log"
	"net/http"
	"strings"
)

var corsmode = "SAME"

func Cors(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if next == nil {
			log.Fatal("Empty handler at %s", r.URL.Path)
			return
		}
		if corsmode == "SAME" {
			orig := r.Header.Get("Origin")
			if orig == "" {
				orig = strings.Split(r.Host, ":")[0]
			}
			// base.Debug("Using Origin (SAME): %s", orig)
			(w).Header().Set("Access-Control-Allow-Origin", orig)
		} else {
			// base.Debug("Using Origin (FIXED): %s", corsmode)
			(w).Header().Set("Access-Control-Allow-Origin", corsmode)
		}

		(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		(w).Header().Set("Access-Control-Allow-Headers", "Last-Modified, Expires, Accept, Cache-Control, Content-Type, Content-Language,Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Pragma")
		(w).Header().Set("Access-Control-Allow-Credentials", "true")
		if (*r).Method == http.MethodOptions {
			//base.Debug("Sending pre-flight cors: %+v", w.Header())
			return
		}

		next.ServeHTTP(w, r)
	})
}
