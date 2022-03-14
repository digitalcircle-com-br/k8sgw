package main

import (
	"log"
	"net/http"
	"time"

	"github.com/digitalcircle-com-br/buildinfo"

	"github.com/digitalcircle-com-br/k8sgw/lib/acme"
	"github.com/digitalcircle-com-br/k8sgw/lib/cfg"
	"github.com/digitalcircle-com-br/k8sgw/lib/muxmgr"
	"github.com/digitalcircle-com-br/k8sgw/lib/mw/cors"
	"github.com/digitalcircle-com-br/k8sgw/lib/mw/csp"
	"github.com/digitalcircle-com-br/k8sgw/lib/mw/helmet"
	"github.com/digitalcircle-com-br/k8sgw/lib/mw/sts"
	"github.com/digitalcircle-com-br/k8sgw/lib/mw/xframe"
	"github.com/gorilla/mux"
)

func main() {
	rootRouter := mux.NewRouter()
	router := mux.NewRouter()
	go func() {
		err := cfg.Setup()
		if err != nil {
			log.Printf("Error: %s", err.Error())
		}
	}()

	go func() {
		for {
			tmprouter, ok := <-muxmgr.Ch()
			if ok {
				log.Printf("Updating router")
				router = tmprouter
			}
		}
	}()
	log.Printf("%s\n", buildinfo.String())

	rootRouter.HandleFunc("/__test", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("TEST: %s %s %s %s", r.Proto, r.Host, r.Method, r.URL.String())
	})
	rootRouter.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Got Request: %s %s %s %s", r.Proto, r.Host, r.Method, r.URL.String())
			h.ServeHTTP(w, r)
		})
	})
	rootRouter.Use(helmet.Helmet)
	rootRouter.Use(cors.Cors)
	rootRouter.Use(csp.Csp)
	rootRouter.Use(sts.Sts)
	rootRouter.Use(xframe.XFrame)

	rootRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router.ServeHTTP(w, r)
	})

	acme.Setup(rootRouter)

	for {
		time.Sleep(time.Minute)
	}
}
