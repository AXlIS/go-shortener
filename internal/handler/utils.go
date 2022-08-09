package handler

import (
	"crypto/rand"

	"github.com/gin-gonic/gin"

	"github.com/AXlIS/go-shortener/internal/utils"
)

func GenerateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

func GetUserID(c *gin.Context) string {
	id := c.GetString(IdentityKey)

	return utils.GenerateString(id)
}
