package events

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ttpham0111/exp-ose/exp/database"
	"github.com/ttpham0111/exp-ose/exp/services"
)

type Service struct {
	Collection database.EventCollectionReader
	Yelp       services.YelpService
}

func (s *Service) find(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"from": "events/find"})
}
