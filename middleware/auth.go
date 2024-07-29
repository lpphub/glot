package middleware

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/lpphub/golib/render"
	"github.com/lpphub/golib/zlog"
	"glot/component/errcode"
	"glot/component/utils"
	"strings"
)

const (
	ContextUser       = "LOGIN_USER"
	ContextUserTenant = "LOGIN_USER_TENANT"
)

func CheckAuthLogin(ctx *gin.Context) {
	// Authorization: Bearer token
	authInfo := ctx.GetHeader("Authorization")
	if authInfo == "" {
		render.JsonWithError(ctx, errcode.ErrNotLogin)
		return
	}
	parts := strings.Fields(authInfo)
	if len(parts) < 2 {
		render.JsonWithError(ctx, errcode.ErrNotLogin)
		return
	}

	token := parts[1]
	uk, err := utils.ParseToken(token, utils.JwtSecret)
	if err != nil {
		zlog.Error(ctx, err.Error())
		render.JsonWithError(ctx, errcode.ErrInvalidToken)
		return
	}

	var user map[string]int64
	_ = jsoniter.UnmarshalFromString(uk, &user)
	ctx.Set(ContextUser, user["uid"])
	ctx.Set(ContextUserTenant, user["tenantId"])
	ctx.Next()
}

func GetLoginUid(ctx *gin.Context) int64 {
	return ctx.GetInt64(ContextUser)
}

func GetLoginTenantId(ctx *gin.Context) int64 {
	return ctx.GetInt64(ContextUserTenant)
}
