package tenant

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/render"
	"github.com/lpphub/golib/zlog"
	"glot/component/errcode"
	"glot/service/domain"
	tenantsrv "glot/service/tenant"
)

func PageListTenant(ctx *gin.Context) {
	var req domain.TenantQuery
	if err := ctx.ShouldBind(&req); err != nil {
		zlog.Warn(ctx, err.Error())
		render.JsonWithError(ctx, errcode.ErrParamInvalid)
		return
	}
	if data, err := tenantsrv.PageList(ctx, req); err != nil {
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, data)
	}
}

func SaveTenant(ctx *gin.Context) {
	var req domain.TenantVO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zlog.Error(ctx, err.Error())
		render.JsonWithError(ctx, errcode.ErrParamInvalid)
		return
	}
	if err := tenantsrv.SaveTenant(ctx, req); err != nil {
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, "ok")
	}
}

func ListRoleScope(ctx *gin.Context) {
	if data, err := tenantsrv.ListRoleScope(ctx); err != nil {
		render.JsonWithError(ctx, err)
	} else {
		render.JsonWithSuccess(ctx, data)
	}
}
