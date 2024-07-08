package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/leedrum/simplebank/util"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "Bearer"
	authorizationPayload    = "authorization_payload"
)

func authMiddleware(tokenMaker *util.TokenMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the value of the Authorization header
		authorization := c.GetHeader(authorizationHeaderKey)
		if !strings.HasPrefix(authorization, authorizationTypeBearer) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Extract the token from the Authorization header
		fields := strings.Fields(authorization)
		if len(fields) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		token := fields[1]
		payload, err := tokenMaker.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		c.Set(authorizationPayload, payload)

		c.Next()
	}
}
