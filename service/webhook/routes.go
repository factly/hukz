package webhook

import (
	"github.com/factly/hukz/model"
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type webhook struct {
	Name     string         `json:"name"`
	URL      string         `json:"url" validate:"required"`
	Enabled  bool           `json:"enabled"`
	EventIDs []uint         `json:"event_ids" validate:"required"`
	Tags     postgres.Jsonb `json:"tags" swaggertype:"primitive,string"`
}

var userContext model.ContextKey = "webhook_user"

// Router webhooks endpoint router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/logs", logs)
	r.Route("/space/{space_id}", func(r chi.Router) {
		r.Get("/", list)
		r.Post("/", create)
		r.Get("/check", check)
	})
	r.Get("/space/{space_id}/webhook/{webhook_id}/logs", webhooklogs)
	r.Route("/{webhook_id}", func(r chi.Router) {
		r.Get("/", details)
		r.Put("/", update)
		r.Delete("/", delete)
	})

	return r
}
