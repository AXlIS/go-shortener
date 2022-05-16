package handler

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func DecompressBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				errorResponse(c, http.StatusInternalServerError, err.Error())
				return
			}

			c.Request.Body = ioutil.NopCloser(gz)
		}

		c.Next()
	}
}
