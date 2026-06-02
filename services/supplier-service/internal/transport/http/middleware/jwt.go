package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pepshot/SoftPlace/shared/auth"
)

const (
	userIDKey = "userID"
	roleKey   = "role"
)

type errorResponse struct {
	Error string `json:"error"`
}

func JWTAuth(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			c.JSON(http.StatusUnauthorized, errorResponse{Error: "authorization required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, errorResponse{Error: "invalid authorization header"})
			c.Abort()
			return
		}

		claims, err := jwtManager.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, errorResponse{Error: "invalid token"})
			c.Abort()
			return
		}

		c.Set(userIDKey, claims.Subject)
		c.Set(roleKey, claims.Role)
		c.Next()
	}
}
