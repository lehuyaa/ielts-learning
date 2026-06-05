package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	sharedjwt "ielts-learning/backend/internal/shared/jwt"
	"ielts-learning/backend/internal/shared/response"
)

const UserIDContextKey = "userID"

func Auth(jwtManager sharedjwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
			c.Abort()
			return
		}

		token, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || strings.TrimSpace(token) == "" {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authorization header")
			c.Abort()
			return
		}

		claims, err := jwtManager.Verify(strings.TrimSpace(token))
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set(UserIDContextKey, claims.UserID)
		c.Next()
	}
}

func GetUserID(c *gin.Context) (uint, bool) {
	value, exists := c.Get(UserIDContextKey)
	if !exists {
		return 0, false
	}

	userID, ok := value.(uint)
	return userID, ok
}
