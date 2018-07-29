package experiences

import (
	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, experiencesService *Service) {
	router.GET("", experiencesService.find)
	router.GET(":id", experiencesService.findId)
	router.GET(":id/events", experiencesService.findIdEvents)
}
