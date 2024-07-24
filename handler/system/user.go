package system

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/render"
	"github.com/lpphub/golib/zlog"
	"glot/component/errcode"
	"glot/service/entity"
	syssrv "glot/service/system"
)

func PageListUser(ctx *gin.Context) {
	var req entity.UserQuery
	if err := ctx.ShouldBind(&req); err != nil {
		zlog.Warn(ctx, err.Error())
		render.JsonWithError(ctx, errcode.ErrParamInvalid)
		return
	}
	if data, err := syssrv.PageListUser(ctx, req); err != nil {
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, data)
	}
}

func SaveUser(ctx *gin.Context) {
	var req entity.User
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zlog.Error(ctx, err.Error())
		render.JsonWithError(ctx, errcode.ErrParamInvalid)
		return
	}
	if err := syssrv.SaveUser(ctx, req); err != nil {
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, "ok")
	}
}

func DelUser(ctx *gin.Context) {
	var req struct {
		Ids []int64 `json:"ids" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zlog.Error(ctx, err.Error())
		render.JsonWithError(ctx, errcode.ErrParamInvalid)
		return
	}
	if err := syssrv.DelUser(ctx, req.Ids); err != nil {
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, "ok")
	}
}
