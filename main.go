package main

import (
	"log"

	"skeleton-v2/bootstrap"
)

func main() {
	app := bootstrap.NewApplication()

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
