package main

import (
	"crypto/tls"
	"github.com/nynrathod/automator-api/bootstrap"
	"golang.org/x/crypto/acme/autocert"
	"log"
)

func main() {

	app := bootstrap.NewApplication()
	// Setup automatic certificate manager
	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("api.example.com"), // Replace with your domain
		Cache:      autocert.DirCache("./certs"),              // Directory to cache the certificates
	}

	// Setup TLS configuration
	tlsConfig := &tls.Config{
		GetCertificate: m.GetCertificate,
		NextProtos:     []string{"http/1.1", "acme-tls/1"},
	}

	// Create a listener with TLS
	ln, err := tls.Listen("tcp", ":443", tlsConfig)
	if err != nil {
		log.Fatalf("Error creating TLS listener: %v", err)
	}

	// Start the Fiber app with HTTPS
	log.Fatal(app.Listener(ln))

}
