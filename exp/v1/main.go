package v1

import (
	"github.com/go-chi/chi"

	"github.com/ttpham0111/exp-ose/v1/events"
	"github.com/ttpham0111/exp-ose/v1/experiences"
)

type Blueprint struct {
	yelpApiKey string
}

func (bp Blueprint) Register(router chi.Router) {
	router.Mount("/events", events.NewRouter(bp.yelpApiKey))
	router.Mount("/experiences", experiences.NewRouter())
}
