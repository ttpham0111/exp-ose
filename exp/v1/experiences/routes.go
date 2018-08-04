package experiences

import (
	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, experiencesService *Service) {
	router.GET("", experiencesService.find)
	router.GET(":id", experiencesService.findId)
	router.GET(":id/comments", experiencesService.getCommentsById)

	router.POST("", experiencesService.create)
	router.POST(":id/comments", experiencesService.addComment)

	router.PUT(":id", experiencesService.update)

	router.DELETE(":id/comments/:commentId", experiencesService.removeComment)
}
