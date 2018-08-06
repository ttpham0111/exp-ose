package activities

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ttpham0111/exp-ose/exp/database"
	"github.com/ttpham0111/exp-ose/exp/services"
	"github.com/ttpham0111/exp-ose/exp/util"
)

type Service struct {
	Yelp       services.YelpService
	Eventbrite services.EventbriteService
}

func (s *Service) find(c *gin.Context) {
	var query services.ServiceQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		util.HandleBindError(c, err)
		return
	}

	activitiesChan := make(chan database.Activity, 10)
	errorsChan := make(chan error)
	doneChan := make(chan interface{})

	go s.getActivitiesFromYelp(query, activitiesChan, doneChan, errorsChan)
	go s.getActivitiesFromEventbrite(query, activitiesChan, doneChan, errorsChan)

	var err error
	activities := make([]database.Activity, 0)
	done := 0
	numServices := 2

	for done < numServices && err == nil {
		select {
		case activity := <-activitiesChan:
			activities = append(activities, activity)
		case <-doneChan:
			done++
		case err = <-errorsChan:
		}
	}

	if err != nil {
		if e, ok := err.(services.ClientError); ok {
			c.Header("Content-Type", "application/json")
			c.String(e.StatusCode, e.Error())
			return
		}
		panic(err)
	}

	c.JSON(http.StatusOK, activities)
}

func (s *Service) getActivitiesFromYelp(
	query services.ServiceQuery,
	activities chan<- database.Activity,
	done chan<- interface{},
	errors chan<- error,
) {
	businesses, err := s.Yelp.FindBusinesses(query)
	if err != nil {
		errors <- err
	}

	for _, business := range businesses {
		activities <- database.Activity{
			Name:     business.Name,
			ImageURL: business.ImageURL,
			Source:   database.Yelp,
			SourceMetadata: database.SourceMetadata{
				"url":          business.Url,
				"location":     business.Location,
				"rating":       business.Rating,
				"review_count": business.ReviewCount,
			},
		}
	}

	done <- nil
}

func (s *Service) getActivitiesFromEventbrite(
	query services.ServiceQuery,
	activities chan<- database.Activity,
	done chan<- interface{},
	errors chan<- error,
) {
	events, err := s.Eventbrite.FindEvents(query)
	if err != nil {
		errors <- err
	}

	for _, event := range events {
		activities <- database.Activity{
			Name:     event.Name.Text,
			ImageURL: event.Logo.Url,
			Source:   database.Eventbrite,
			SourceMetadata: database.SourceMetadata{
				"url":          event.Url,
				"description":  event.Description.Text,
				"starts_at":    event.StartsAt.Utc,
				"ends_at":      event.EndsAt.Utc,
				"online_event": event.OnlineEvent,
			},
		}
	}

	done <- nil
}
