package activities

import (
	"net/http"
	"time"

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

	var startsAt, endsAt time.Time
	const layout = "2006-01-02T15:04:05Z"

	for _, event := range events {
		startsAt, err = time.Parse(layout, event.StartsAt.Utc)
		if err != nil {
			errors <- err
		}

		endsAt, err = time.Parse(layout, event.EndsAt.Utc)
		if err != nil {
			errors <- err
		}

		activities <- database.Activity{
			Name:     event.Name.Text,
			ImageURL: event.Logo.Url,
			StartsAt: &startsAt,
			EndsAt:   &endsAt,
			Source:   database.Eventbrite,
			SourceMetadata: database.SourceMetadata{
				"url": event.Url,
				// "description":  event.Description.Text,
				"online_event": event.OnlineEvent,
				"venue":        event.Venue,
			},
		}
	}

	done <- nil
}
