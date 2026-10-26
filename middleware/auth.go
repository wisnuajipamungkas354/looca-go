package middleware

import (
	"net/http"
	"strings"

	"looca-go/utils"

	"github.com/gin-gonic/gin"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Error(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			utils.Error(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		tokenStr := parts[1]
		claims, err := utils.ParseToken(secret, tokenStr)
		if err != nil {
			utils.Error(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		c.Set("admin_id", claims.AdminID)
		c.Set("admin_email", claims.Email)
		c.Next()
	}
}
