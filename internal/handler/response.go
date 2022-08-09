package handler

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Message string
}

func errorResponse(c *gin.Context, statusCode int, message string) {
	log.Println(message)
	c.AbortWithStatusJSON(statusCode, Error{Message: message})
}
