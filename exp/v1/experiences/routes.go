package experiences

import (
	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, experiencesService *Service) {
	router.GET("", experiencesService.find)
	router.GET(":id", experiencesService.findId)
	router.GET(":id/comments", experiencesService.getCommentsById)
	router.GET(":id/activities", experiencesService.getActivitiesById)

	router.POST("", experiencesService.create)
	router.POST(":id/ratings", experiencesService.addRating)
	router.POST(":id/comments", experiencesService.addComment)
	router.POST(":id/activities", experiencesService.addActivity)

	router.PUT(":id", experiencesService.update)

	router.DELETE(":id/ratings/:ratingId", experiencesService.removeRating)
	router.DELETE(":id/comments/:commentId", experiencesService.removeComment)
	router.DELETE(":id/activities/:activityId", experiencesService.removeActivity)
}
