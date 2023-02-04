package middleware

import (
	"github.com/gin-gonic/gin"
)

// InitMiddleware register all the service into the 'Context.Keys'
func InitMiddleware(service []interface{}) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Keys = make(map[string]interface{})
		context.Keys["user"] = service[0]
		context.Next()
	}
}
