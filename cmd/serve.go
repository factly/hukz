package cmd

import (
	"log"
	"net/http"

	"github.com/factly/hukz/config"
	"github.com/factly/hukz/service"
	"github.com/factly/hukz/util"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts server for hukz.",
	Run: func(cmd *cobra.Command, args []string) {
		// setup database
		config.SetupDB()

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
		if err := http.ListenAndServe(":8000", r); err != nil {
			log.Fatal(err)
		}
	},
}
