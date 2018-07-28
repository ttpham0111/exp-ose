package events

import (
	"github.com/go-chi/chi"
)

func NewRouter(yelpApiKey) *chi.Mux {
	router := chi.NewRouter()

	events := service{
		yelpApiKey: yelpApiKey,
	}

	router.Get("/events", events.find)
	router.Get("/events/{eventId}", events.get)

	return router
}
