package experiences

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"

	"github.com/ttpham0111/exp-ose/exp/database"
	"github.com/ttpham0111/exp-ose/exp/util"
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
	var query database.ExperienceQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		e := util.FirstError(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid value for " + e.Name})
		return
	}

	var modifier database.QueryModifier
	if err := c.ShouldBindQuery(&modifier); err != nil {
		e := util.FirstError(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid value for " + e.Name})
		return
	}

	experiences, err := s.ExperienceCollection.Find(query, modifier)
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
	if _, exists := c.GetQuery("include_collaborators"); exists {
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

	var modifier database.QueryModifier
	if err := c.ShouldBindQuery(&modifier); err != nil {
		e := util.FirstError(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid value for " + e.Name})
	}

	comments, err := s.ExperienceCollection.GetCommentsById(experienceId, modifier)
	if err != nil {
		if e, ok := err.(database.NoResultFound); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
			return
		}
		panic(err)
	}

	c.JSON(http.StatusOK, comments)
}
