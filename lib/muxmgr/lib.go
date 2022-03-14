package muxmgr

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

var theMux *mux.Router

// var ch = make(chan map[string]string, 0)
var ch = make(chan *mux.Router, 0)
var cli = http.Client{}

func Update(s string) {
	log.Printf("Loading config: %s", s)
	var config = make(map[string]map[string]string)
	err := yaml.Unmarshal([]byte(s), &config)
	if err != nil {
		log.Printf("Error updating config: %s", err.Error())
		return
	}
	tmpmux := mux.NewRouter()
	mapRouters := make(map[string]*mux.Router)
	for h, paths := range config {
		hrouter := mux.NewRouter()
		for path, target := range paths {
			func(p, t string) {
				log.Printf("Setting up route: %s %s => %s", h, p, t)
				hrouter.PathPrefix(p).HandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

					r.Host = t
					r.RequestURI = ""
					urlstr := r.URL.String()
					t = t + "/" + strings.Replace(urlstr, p, "", 1)
					t = strings.ReplaceAll(t, "//", "/")
					urlstr = "http://" + t

					nu, err := url.Parse(urlstr)
					if err != nil {
						log.Printf("Error: %s", err.Error())
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return

					}
					log.Printf("Calling %s => [%s] %s for host %s", r.URL.String(), r.Method, urlstr, r.Host)
					r.URL = nu
					res, err := cli.Do(r)
					if err != nil {
						log.Printf("Error: %s", err.Error())
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return

					}

					for k, v := range res.Header {
						for _, vv := range v {
							w.Header().Add(k, vv)
						}
					}
					w.WriteHeader(res.StatusCode)

					defer res.Body.Close()
					io.Copy(w, res.Body)

				}))
			}(path, target)
		}
		mapRouters[h] = hrouter
	}
	tmpmux.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rt, ok := mapRouters[r.Host]
		if !ok {
			rt, ok = mapRouters["*"]
			if !ok {
				log.Printf("NO router found for: %s %s %s %s", r.Proto, r.Host, r.Method, r.URL.String())
				http.NotFound(w, r)
				return
			} else {
				log.Printf("Calling DEFAULT router for host: %s %s %s %s", r.Proto, r.Host, r.Method, r.URL.String())
			}
		} else {
			log.Printf("Calling router for host: %s %s %s %s", r.Proto, r.Host, r.Method, r.URL.String())

		}
		rt.ServeHTTP(w, r)
	})
	theMux = tmpmux
	ch <- theMux
}

func Ch() chan *mux.Router {
	return ch
}
