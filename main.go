package main

import (
	"log"
	"net/http"

	"github.com/factly/web-hooks-service/config"
	"github.com/factly/web-hooks-service/model"
	"github.com/factly/web-hooks-service/service"
)

func main() {
	// setup environment vars
	config.SetupVars()

	// setup database
	config.SetupDB()

	// apply database migrations
	model.Migration()

	r := service.RegisterRoutes()
	if err := http.ListenAndServe(":7790", r); err != nil {
		log.Fatal(err)
	}

}
