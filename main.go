package main

import (
	"github.com/nynrathod/automator-api/bootstrap"
	"log"
)

func main() {

	app := bootstrap.NewApplication()
	log.Fatal(app.Listen(":80"))

}
