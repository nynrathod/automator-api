package main

import (
	"crypto/tls"
	"github.com/nynrathod/automator-api/bootstrap"
	"github.com/nynrathod/automator-api/config"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net"
)

func main() {
	app := bootstrap.NewApplication()
	env := config.EnvConfigs.ENV

	var ln net.Listener
	var err error

	if env == "PROD" {
		m := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("api.uoozer.com"),
			Cache:      autocert.DirCache("./certs"),
		}

		cfg := &tls.Config{
			GetCertificate: m.GetCertificate,
			NextProtos: []string{
				"http/1.1", "acme-tls/1",
			},
		}

		ln, err = tls.Listen("tcp", ":443", cfg)
		if err != nil {
			log.Fatalf("Failed to create TLS listener: %v", err)
		}
	} else {
		log.Fatal(app.Listen(":3000"))
	}

	log.Fatal(app.Listener(ln))
}
