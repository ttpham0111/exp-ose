package experiences

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"

	"github.com/ttpham0111/exp-ose/exp/database"
	"github.com/ttpham0111/exp-ose/exp/util"
)

type Service struct {
	ExperienceCollection database.ExperienceCollectionReadWriter
}

func (s *Service) find(c *gin.Context) {
	var query database.ExperienceQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		util.HandleBindError(c, err)
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

	c.JSON(http.StatusOK, experience)
}

func (s *Service) getCommentsById(c *gin.Context) {
	experienceId := c.Param("id")

	var modifier database.QueryModifier
	if err := c.ShouldBindQuery(&modifier); err != nil {
		util.HandleBindError(c, err)
		return
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

func (s *Service) create(c *gin.Context) {
	var experience database.Experience
	if err := c.ShouldBindJSON(&experience); err != nil {
		util.HandleBindError(c, err)
		return
	}

	if err := s.ExperienceCollection.Create(&experience); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, &experience)
}

func (s *Service) addComment(c *gin.Context) {
	experienceId := c.Param("id")

	var comment database.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		util.HandleBindError(c, err)
		return
	}

	if err := s.ExperienceCollection.AddComment(experienceId, &comment); err != nil {
		if e, ok := err.(database.NoResultFound); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
			return
		}
		panic(err)
	}

	c.JSON(http.StatusCreated, &comment)
}

func (s *Service) update(c *gin.Context) {
	experienceId := c.Param("id")

	var experience database.Experience
	if err := c.ShouldBindJSON(&experience); err != nil {
		util.HandleBindError(c, err)
		return
	}

	if err := s.ExperienceCollection.Update(experienceId, &experience); err != nil {
		if e, ok := err.(database.NoResultFound); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
			return
		}
		panic(err)
	}

	c.JSON(http.StatusOK, &experience)
}

func (s *Service) removeComment(c *gin.Context) {
	commentId := c.Param("commentId")

	if err := s.ExperienceCollection.RemoveComment(commentId); err != nil {
		if e, ok := err.(database.NoResultFound); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
			return
		}
		panic(err)
	}

	c.JSON(http.StatusNoContent, nil)
}
