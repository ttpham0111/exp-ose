package events

import (
	"net/http"

	"github.com/ttpham0111/exp-ose/exp/services"
	"github.com/ttpham0111/exp-ose/exp/util"
)

type service struct {
	yelp *services.YelpService
}

func (s *service) find(w http.ResponseWriter, r *http.Request) {
	// TODO: goroutine for other API's
	if yelpBusinesses, err := s.yelp.FindBusinesses(r.URL.Query()); err != nil {
		if e, ok := err.(services.ClientError); ok {
			util.JsonResponse(w, e.Message, e.StatusCode)
		} else {
			panic(err)
		}
	} else {
		util.JsonResponse(w, yelpBusinesses, http.StatusOK)
	}
}
