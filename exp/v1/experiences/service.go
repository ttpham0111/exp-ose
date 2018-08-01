package experiences

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ttpham0111/exp-ose/exp/database"
)

type Service struct {
	ExperienceCollection database.ExperienceCollectionReader
	EventCollection      database.EventCollectionReader
}

type ExperienceResponse struct {
	*database.Experience
	Events        []*database.ExperienceEvent `json:"events"`
	Collaborators []database.UserId           `json:"collaborators"`
}

func (s *Service) find(c *gin.Context) {
	owner := c.GetString("owner")
	isPublic := c.GetBool("is_public")
	name := c.GetString("name")
	tags := c.GetStringSlice("tags")

	sortBy := c.GetString("sort_by")
	sortAsc := c.GetBool("sort_asc")

	skip := c.GetInt("skip")
	limit := c.GetInt("limit")

	experiences, err := s.ExperienceCollection.Find(
		database.ExperienceQuery{
			Owner:    owner,
			IsPublic: isPublic,
			Name:     name,
			Tags:     tags,
		},
		database.QueryModifier{
			SortBy:  sortBy,
			SortAsc: sortAsc,
			Skip:    skip,
			Limit:   limit,
		},
	)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, experiences)
}

func (s *Service) findId(c *gin.Context) {
	// TODO: goroutine
	experienceId := c.Param("id")

	experience, err := s.ExperienceCollection.FindId(experienceId)
	if err != nil {
		if e, ok := err.(database.NoResultFound); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
			return
		}
		panic(err)
	}

	events, err := s.ExperienceCollection.GetEventsById(experienceId)

	var collaborators []database.UserId
	if c.GetBool("include_collaborators") {
		collaborators, err = s.ExperienceCollection.GetCollaboratorsById(experienceId)
		if err != nil {
			panic(err)
		}
	}

	c.JSON(http.StatusOK, ExperienceResponse{
		Experience:    experience,
		Events:        events,
		Collaborators: collaborators,
	})
}

func (s *Service) getCommentsById(c *gin.Context) {
	experienceId := c.Param("id")
	skip := c.GetInt("skip")
	limit := c.GetInt("limit")

	comments, err := s.ExperienceCollection.GetCommentsById(
		experienceId,
		database.QueryModifier{
			Skip:  skip,
			Limit: limit,
		},
	)
	if err != nil {
		if e, ok := err.(database.NoResultFound); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
			return
		}
		panic(err)
	}

	c.JSON(http.StatusOK, comments)
}
