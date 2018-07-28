package events

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/ttpham0111/exp-ose/exp"
	"github.com/ttpham0111/exp-ose/exp/services/yelp"
)

type service struct {
	yelpApiKey string
}

type businessQuery struct {
	yelp.BusinessQuery
}

func (s *service) find(w http.ResponseWriter, r *http.Request) {
	var q businessQuery
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		panic(err)
	}

	// TODO: goroutine for other API's
	yelpBusinesses, err := yelp.find(q)
	if err != nil {
		panic(err)
	}

	exp.JsonResponse(w, yelpBusinesses, http.StatusOK)
}

func (s *service) get(w http.ResponseWriter, r *http.Request) {
	eventId := chi.URLParam(r, "eventId")
}
