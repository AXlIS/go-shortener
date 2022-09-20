package handler

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Error string
}

func errorResponse(c *gin.Context, statusCode int, message string) {
	log.Println(message)
	c.AbortWithStatusJSON(statusCode, Error{Error: message})
}
