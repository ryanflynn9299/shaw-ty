package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)
import "golang.org/x/time/rate"

func RateLimiter() gin.HandlerFunc {
	limiter := rate.NewLimiter(3, 10) // accept 3 events per second, with burst load of 10
	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()
		} else {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
		}
	}
}
