package middleware

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS(originStr string) gin.HandlerFunc {
	cfg := cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
			"Origin",
			"Accept",
			"X-Requested-With",
			"Access-Control-Request-Method",
			"Access-Control-Request-Headers",
		},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	trimmed := strings.TrimSpace(originStr)
	if trimmed == "*" || trimmed == "" {
		// Allow any origin dynamically so browsers don't block credentials or custom origins
		cfg.AllowOriginFunc = func(origin string) bool {
			return true
		}
	} else {
		parts := strings.Split(trimmed, ",")
		var allowed []string
		for _, p := range parts {
			if s := strings.TrimSpace(p); s != "" {
				allowed = append(allowed, s)
			}
		}
		if len(allowed) == 0 {
			cfg.AllowOriginFunc = func(origin string) bool {
				return true
			}
		} else {
			cfg.AllowOrigins = allowed
		}
	}

	return cors.New(cfg)
}

