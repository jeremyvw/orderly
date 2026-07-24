package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const CtxUserID = "userID"

type TokenParser func(token string) (int64, error)

func Auth(parse TokenParser) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			abortUnauthorized(c)
			return
		}

		parts := strings.Fields(header)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			abortUnauthorized(c)
			return
		}

		userID, err := parse(parts[1])
		if err != nil {
			abortUnauthorized(c)
			return
		}

		c.Set(CtxUserID, userID)
		c.Next()
	}
}

func UserIDFrom(c *gin.Context) (int64, bool) {
	value, exists := c.Get(CtxUserID)
	if !exists {
		return 0, false
	}

	userID, ok := value.(int64)
	return userID, ok
}

func abortUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
}
