package service

import (
	"fmt"
	"net/http"

	"github.com/factly/web-hooks-service/config"
	_ "github.com/factly/web-hooks-service/docs"
	"github.com/factly/web-hooks-service/service/event"
	"github.com/factly/web-hooks-service/service/webhook"
	"github.com/factly/x/healthx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/middlewarex"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
)

// RegisterRoutes registers API routes
func RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(loggerx.Init())
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	if viper.IsSet("mode") && viper.GetString("mode") == "development" {
		r.Get("/swagger/*", httpSwagger.WrapHandler)
		fmt.Println("Swagger @ http://localhost:7790/swagger/index.html")
	}

	sqlDB, _ := config.DB.DB()

	healthx.RegisterRoutes(r, healthx.ReadyCheckers{
		"database": sqlDB.Ping,
	})

	r.With(middlewarex.CheckUser).Group(func(r chi.Router) {
		r.Mount("/webhooks", webhook.Router())
		r.Mount("/events", event.Router())
	})

	return r
}
