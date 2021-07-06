package event

import (
	"github.com/factly/hukz/model"
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type event struct {
	Name  string         `json:"name"`
	Event string         `json:"event" validate:"required"`
	Tags  postgres.Jsonb `json:"tags" swaggertype:"primitive,string"`
}

var userContext model.ContextKey = "event_user"

// Router events endpoint router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)
	r.Post("/", create)

	r.Route("/{event_id}", func(r chi.Router) {
		r.Get("/", details)
		r.Put("/", update)
		r.Delete("/", delete)
	})

	return r
}
