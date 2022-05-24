package handler

import (
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

var (
	secretKey = []byte("super-secret")
)

const (
	IdLength    = 6
	IdentityKey = "userId"
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

func checkCookie(c *gin.Context) (bool, error) {

	cookie, err := c.Cookie(IdentityKey)
	if errors.Is(http.ErrNoCookie, err) {
		return false, nil
	}

	data, err := hex.DecodeString(cookie)
	if err != nil {
		return false, err
	}

	h := hmac.New(sha256.New, secretKey)
	h.Write(data[:IdLength])
	sign := h.Sum(nil)

	if hmac.Equal(data[IdLength:], sign) {
		c.Set(IdentityKey, string(data[:IdLength]))
		return true, nil
	}

	return false, nil
}

func CookieHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieIsCorrect, err := checkCookie(c)
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		if !cookieIsCorrect {
			id, err := GenerateRandom(IdLength)
			if err != nil {
				errorResponse(c, http.StatusInternalServerError, err.Error())
			}

			h := hmac.New(sha256.New, secretKey)
			h.Write(id)
			sign := h.Sum(nil)

			c.SetCookie(IdentityKey, fmt.Sprintf("%x%x", id, sign), 3600, "/", "", false, true)
			c.Set(IdentityKey, string(id))
		}

		c.Next()
	}
}
