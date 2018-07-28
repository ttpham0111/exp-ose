package events

import (
	"github.com/go-chi/chi"

	"github.com/ttpham0111/exp-ose/exp/services"
)

func NewRouter(yelpApiKey string) *chi.Mux {
	router := chi.NewRouter()

	events := service{
		yelp: services.NewYelpService(yelpApiKey),
	}

	router.Get("/", events.find)

	return router
}
