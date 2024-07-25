package system

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/render"
	"github.com/lpphub/golib/zlog"
	"github.com/spf13/cast"
	"glot/component/errcode"
	repo "glot/repository"
	"glot/service/consts"
	"glot/service/domain"
	syssrv "glot/service/system"
)

func PageListRole(ctx *gin.Context) {
	var req domain.RoleQuery
	if err := ctx.ShouldBind(&req); err != nil {
		zlog.Warn(ctx, err.Error())
		render.JsonWithError(ctx, errcode.ErrParamInvalid)
		return
	}
	if data, err := syssrv.PageListRole(ctx, req); err != nil {
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, data)
	}
}

func ListAllRole(ctx *gin.Context) {
	if data, err := syssrv.ListAllRole(ctx); err != nil {
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, data)
	}
}

func SaveRole(ctx *gin.Context) {
	var req repo.Role
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zlog.Error(ctx, err.Error())
		render.JsonWithError(ctx, errcode.ErrParamInvalid)
		return
	}
	if err := syssrv.SaveRole(ctx, req); err != nil {
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, "ok")
	}
}

func DelRole(ctx *gin.Context) {
	var req struct {
		Ids []int64 `json:"ids" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zlog.Error(ctx, err.Error())
		render.JsonWithError(ctx, errcode.ErrParamInvalid)
		return
	}
	if err := syssrv.DelRole(ctx, req.Ids); err != nil {
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, "ok")
	}
}

func GetRoleMenu(ctx *gin.Context) {
	roleId := ctx.Query("roleId")
	mode := ctx.Query("mode")
	if roleId == "" || mode == "" {
		render.JsonWithError(ctx, errcode.ErrParamInvalid)
		return
	}

	types := mapMode(cast.ToInt(mode))
	if data, err := syssrv.GetRoleMenu(ctx, cast.ToInt64(roleId), types); err != nil {
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, data)
	}
}

func BindRoleMenu(ctx *gin.Context) {
	var req struct {
		RoleId int64   `json:"roleId" binding:"required"`
		Mode   int     `json:"mode" binding:"required"`
		Ids    []int64 `json:"ids"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zlog.Error(ctx, err.Error())
		render.JsonWithError(ctx, errcode.ErrParamInvalid)
		return
	}

	types := mapMode(req.Mode)
	if err := syssrv.BindRoleMenu(ctx, req.RoleId, types, req.Ids); err != nil {
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, "ok")
	}
}

func mapMode(mode int) []int {
	// 1-菜单 2-按钮
	modeMap := map[int][]int{
		1: {consts.MenuDir, consts.MenuOpt},
		2: {consts.MenuButton},
	}
	if rsc, ok := modeMap[mode]; ok {
		return rsc
	}
	return []int{0}
}
