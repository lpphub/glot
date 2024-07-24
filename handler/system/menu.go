package system

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/render"
	"github.com/lpphub/golib/zlog"
	"glot/component/errcode"
	"glot/middleware"
	"glot/service/entity"
	syssrv "glot/service/system"
)

func SaveMenu(ctx *gin.Context) {
	var req entity.Menu
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zlog.Error(ctx, err.Error())
		render.JsonWithError(ctx, errcode.ErrParamInvalid)
		return
	}
	if err := syssrv.SaveMenu(ctx, req); err != nil {
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, "ok")
	}
}

func PageListMenu(ctx *gin.Context) {
	var req entity.PageQuery
	if err := ctx.ShouldBind(&req); err != nil {
		zlog.Warn(ctx, err.Error())
		render.JsonWithError(ctx, errcode.ErrParamInvalid)
		return
	}
	if data, err := syssrv.PageListMenu(ctx, req); err != nil {
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, data)
	}
}

func GetMenuTree(ctx *gin.Context) {
	if data, err := syssrv.GetMenuTree(ctx); err != nil {
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, data)
	}
}

func GetMenuButton(ctx *gin.Context) {
	if data, err := syssrv.GetMenuButton(ctx); err != nil {
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, data)
	}
}

func DelMenu(ctx *gin.Context) {
	var req struct {
		Ids []int64 `json:"ids" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zlog.Error(ctx, err.Error())
		render.JsonWithError(ctx, errcode.ErrParamInvalid)
		return
	}
	if err := syssrv.DelMenu(ctx, req.Ids); err != nil {
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, "ok")
	}
}

func GetUserRoutes(ctx *gin.Context) {
	uid := middleware.GetLoginUid(ctx)
	if user, err := syssrv.GetUserRouteMenus(ctx, uid); err != nil {
		zlog.Error(ctx, err.Error())
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, user)
	}
}

func IsExistRoute(ctx *gin.Context) {
	routeName := ctx.Query("routeName")
	if routeName == "" {
		render.JsonWithError(ctx, errcode.ErrParamInvalid)
		return
	}
	if user, err := syssrv.IsExistRoute(ctx, routeName); err != nil {
		zlog.Error(ctx, err.Error())
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, user)
	}
}
