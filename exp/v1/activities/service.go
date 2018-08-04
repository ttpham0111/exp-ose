package activities

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ttpham0111/exp-ose/exp/services"
)

type Service struct {
	Yelp services.YelpService
}

func (s *Service) find(c *gin.Context) {
	var query services.YelpQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	businesses, err := s.Yelp.FindBusinesses(query)
	if err != nil {
		if e, ok := err.(services.ClientError); ok {
			c.Header("Content-Type", "application/json")
			c.String(e.StatusCode, e.Error())
			return
		}
		panic(err)
	}

	c.JSON(http.StatusOK, businesses)
}
