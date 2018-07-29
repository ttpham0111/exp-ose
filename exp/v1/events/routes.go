package events

import (
	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, eventsService *Service) {
	router.GET("", eventsService.find)
}
