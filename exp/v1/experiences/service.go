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

func (s *Service) find(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"from": "experience/find"})
}

func (s *Service) findId(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"from": "experience/findId"})
}

func (s *Service) findIdEvents(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"from": "experience/findIdEvents"})
}
