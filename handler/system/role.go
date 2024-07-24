package system

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/render"
	"github.com/lpphub/golib/zlog"
	"github.com/spf13/cast"
	"glot/component/errcode"
	repo "glot/repository"
	"glot/service/consts"
	"glot/service/entity"
	syssrv "glot/service/system"
)

func PageListRole(ctx *gin.Context) {
	var req entity.RoleQuery
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

func GetRoleResource(ctx *gin.Context) {
	roleId := ctx.Query("roleId")
	rscType := ctx.Query("rscType")
	if roleId == "" || rscType == "" {
		render.JsonWithError(ctx, errcode.ErrParamInvalid)
		return
	}

	types := mapResourceType(cast.ToInt(rscType))
	if data, err := syssrv.GetRoleResource(ctx, cast.ToInt64(roleId), types); err != nil {
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, data)
	}
}

func BindRoleResource(ctx *gin.Context) {
	var req struct {
		RoleId  int64   `json:"roleId" binding:"required"`
		RscType int     `json:"rscType" binding:"required"`
		Ids     []int64 `json:"ids"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zlog.Error(ctx, err.Error())
		render.JsonWithError(ctx, errcode.ErrParamInvalid)
		return
	}

	types := mapResourceType(req.RscType)
	if err := syssrv.BindRoleResource(ctx, req.RoleId, types, req.Ids); err != nil {
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, "ok")
	}
}

func mapResourceType(rscType int) []int {
	// 1-菜单 2-按钮
	rscMap := map[int][]int{
		1: {consts.ResourceDir, consts.ResourceMenu},
		2: {consts.ResourceButton},
	}
	if rsc, ok := rscMap[rscType]; ok {
		return rsc
	}
	return []int{0}
}
