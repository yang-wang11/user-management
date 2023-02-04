package middleware

//
//import (
//	"api-gateway/pkg/util"
//	"github.com/gin-gonic/gin"
//	"time"
//)
//
//// JWT token验证中间件
//func JWT() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		var code int
//		var data interface{}
//		code = 200
//		token := c.GetHeader("Authorization")
//		if token == "" {
//			code = 404
//		} else {
//			claims, err := util.ParseToken(token)
//			if err != nil {
//				code = -1 // ErrorAuthCheckTokenFail
//			} else if time.Now().Unix() > claims.ExpiresAt {
//				code = -2 // e.ErrorAuthCheckTokenTimeout
//			}
//		}
//		if code != 0 { // SUCCESS
//			c.JSON(200, gin.H{
//				"status": code,
//				"data":   data,
//			})
//			c.Abort()
//			return
//		}
//		c.Next()
//	}
//}
