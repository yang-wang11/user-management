package routers

import (
	"api-gateway/pkg/handler"
	"api-gateway/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(server ...interface{}) *gin.Engine {
	ginRouter := gin.Default()
	ginRouter.Use(middleware.Cors(), middleware.InitMiddleware(server))
	v1 := ginRouter.Group("/api/v1")
	v1.GET("ping", func(context *gin.Context) {
		context.JSON(200, "success")
	})
	v1.POST("/user/register", handler.UserRegister)
	v1.POST("/user/login", handler.UserLogin)
	return ginRouter
}
