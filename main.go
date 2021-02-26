package main

import (
	"log"
	"net/http"

	"github.com/factly/hukz/config"
	"github.com/factly/hukz/model"
	"github.com/factly/hukz/service"
	"github.com/factly/hukz/util"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	go func() {
		promRouter := chi.NewRouter()
		promRouter.Mount("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":8001", promRouter))
	}()

	r := service.RegisterRoutes()
	if err := http.ListenAndServe(":7790", r); err != nil {
		log.Fatal(err)
	}

}
