package handler

import (
	"api-gateway/pkg/common"
	"api-gateway/pkg/pb"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserRegister(ctx *gin.Context) {
	var userReq pb.UserRequest
	common.PanicIfErrorIsNotNil(ctx.Bind(&userReq))
	// 从context.Keys 里面获取之前注册好的services ()
	userService := ctx.Keys["user"].(pb.UserServiceClient)
	userResp, err := userService.UserRegister(ctx, &userReq)
	common.PanicIfErrorIsNotNil(err)
	r := common.Response{
		Data:   userResp,
		Status: uint(userResp.ReturnCode),
		Error:  "",
	}
	ctx.JSON(http.StatusOK, r)
}

func UserLogin(ctx *gin.Context) {
	var userReq pb.UserRequest
	common.PanicIfErrorIsNotNil(ctx.Bind(&userReq))
	userService := ctx.Keys["user"].(pb.UserServiceClient)
	userResp, err := userService.UserLogin(ctx, &userReq)
	common.PanicIfErrorIsNotNil(err)
	token, err := common.GenerateToken(uint(userResp.UserDetail.UserId))
	r := common.Response{
		Data: common.TokenData{
			User:  userResp.UserDetail,
			Token: token,
		},
		Status: uint(userResp.ReturnCode),
		Error:  "",
	}
	ctx.JSON(http.StatusOK, r)
}
