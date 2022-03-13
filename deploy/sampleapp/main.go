package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		me := os.Getenv("ME")
		if me == "" {
			me = "Could not find ENV ME"
		}
		log.Printf("Got:%s %s %s", r.Host, r.Method, r.URL.String())
		hn, err := os.Hostname()
		if err != nil {
			log.Printf("error obtaining hostname: %s", err.Error())
		}
		w.Write([]byte(fmt.Sprintf("Hello from %s\nHostname: %s", time.Now().String(), hn)))
		r.Body.Close()
	})

	http.ListenAndServe(":80", nil)
}
