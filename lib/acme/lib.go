package acme

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
)

func Setup(rr *mux.Router) {
	//go acmecache.Setup()

	certManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache("/caroot"), //&acmecache.K8SAcmeCache{}, //
	}

	server := &http.Server{
		Addr: ":443",

		TLSConfig: &tls.Config{
			PreferServerCipherSuites: true,
			// Only use curves which have assembly implementations
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519, // Go 1.8 only
			},
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

				// Best disabled, as they don't provide Forward Secrecy,
				// but might be necessary for some clients
				// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			},
			GetCertificate:     certManager.GetCertificate,
			InsecureSkipVerify: true,
		},
	}

	//Groutine for HTTP Server
	go func() {

		err := http.ListenAndServe(":80", certManager.HTTPHandler(nil))
		if err != nil {
			log.Printf(err.Error())
		}

	}()

	//Goroutine for HTTPS server
	go func() {
		server.Handler = rr
		err := server.ListenAndServeTLS("", "")
		if err != nil {
			log.Printf(err.Error())
		}
	}()
}
