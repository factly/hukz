package main

import (
	"log"
	"net/http"

	"github.com/factly/web-hooks-service/config"
	"github.com/factly/web-hooks-service/model"
	"github.com/factly/web-hooks-service/service"
	"github.com/factly/web-hooks-service/util"
)

// @title Webhooks API
// @version 1.0
// @description Webhooks Service API

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:7790
// @BasePath /
func main() {
	// setup environment vars
	config.SetupVars()

	// setup database
	config.SetupDB()

	// apply database migrations
	model.Migration()

	util.ConnectNats()
	defer util.NC.Close()

	if err := util.SubscribeExistingEvents(); err != nil {
		log.Fatal(err)
	}

	r := service.RegisterRoutes()
	if err := http.ListenAndServe(":7790", r); err != nil {
		log.Fatal(err)
	}

}
