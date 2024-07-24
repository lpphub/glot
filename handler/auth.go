package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/render"
	"github.com/lpphub/golib/zlog"
	"glot/component/errcode"
	"glot/middleware"
	srv "glot/service"
)

func Login(ctx *gin.Context) {
	var login struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&login); err != nil {
		render.JsonWithError(ctx, errcode.ErrParamInvalid)
		return
	}
	if token, err := srv.Login(ctx, login.Username, login.Password); err != nil {
		zlog.Error(ctx, err.Error())
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, gin.H{"token": token})
	}
}

func LoginByToken(ctx *gin.Context) {
	uid := middleware.GetLoginUid(ctx)
	if user, err := srv.GetLoginUser(ctx, uid); err != nil {
		zlog.Error(ctx, err.Error())
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, user)
	}
}
