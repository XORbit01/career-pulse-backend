package middleware

import (
	"github.com/gin-gonic/gin"
)

// ConfigMiddleware sets configuration values in the context
func ConfigMiddleware(staticURL, uploadsPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("staticURL", staticURL)
		c.Set("uploadsPath", uploadsPath)
		c.Next()
	}
}
