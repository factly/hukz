package event

import (
	"github.com/factly/hukz/model"
	"github.com/go-chi/chi"
)

type event struct {
	Name string `json:"name" validate:"required"`
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
