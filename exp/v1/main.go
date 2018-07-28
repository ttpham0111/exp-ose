package v1

import (
	"github.com/go-chi/chi"

	"github.com/ttpham0111/exp-ose/exp/v1/events"
	"github.com/ttpham0111/exp-ose/exp/v1/experiences"
)

type Blueprint struct {
	YelpApiKey string
}

func (bp Blueprint) Register(router chi.Router) {
	router.Mount("/events", events.NewRouter(bp.YelpApiKey))
	router.Mount("/experiences", experience.NewRouter())
}
