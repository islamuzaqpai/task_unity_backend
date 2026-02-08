package routes

import (
	"enactus/internal/handler"

	"github.com/gin-gonic/gin"
)

func UserRoutes(userH *handler.UserHandler, router *gin.Engine) {
	router.POST("/signUp", userH.Register)
}
