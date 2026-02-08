package handler

import (
	"enactus/internal/models"
	"enactus/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

type UserHandlerInterface interface {
	Register(ctx *gin.Context)
}

type UserHandler struct {
	UserS *service.UserService
}

func (userH *UserHandler) Register(c *gin.Context) {
	var registerInput models.RegisterInput

	err := c.ShouldBindJSON(&registerInput)
	if err != nil {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	_, err = userH.UserS.Register(c.Request.Context(), registerInput)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		log.Printf("fail: %v", err)
		return
	}

	c.JSON(201, gin.H{"status": "ok"})
}
