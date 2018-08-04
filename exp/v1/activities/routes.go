package activities

import (
	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, activitiesService *Service) {
	router.GET("", activitiesService.find)
}
