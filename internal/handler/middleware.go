package handler

import (
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	secretKey = []byte("super-secret")
)

const (
	IDLength    = 6
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
	h.Write(data[:IDLength])
	sign := h.Sum(nil)

	if hmac.Equal(data[IDLength:], sign) {
		c.Set(IdentityKey, string(data[:IDLength]))
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
			id, err := GenerateRandom(IDLength)
			if err != nil {
				errorResponse(c, http.StatusInternalServerError, err.Error())
			}

			h := hmac.New(sha256.New, secretKey)
			h.Write(id)
			sign := h.Sum(nil)

			c.SetCookie(IdentityKey, fmt.Sprintf("%x%x", id, sign), 3600, "/", "", false, false)
			c.Set(IdentityKey, string(id))
		}

		c.Next()
	}
}
