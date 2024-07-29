package service

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"glot/component/errcode"
	"glot/component/utils"
	repo "glot/repository"
	"glot/service/consts"
	"glot/service/domain"
)

func Login(ctx *gin.Context, username, password string) (string, error) {
	var user repo.User
	repo.GetDB(ctx).Where("username=? and password=? and status=1", username, password).Take(&user)
	if user.ID == 0 {
		return "", errcode.ErrUserNotFound
	}

	var tenant repo.Tenant
	repo.GetDB(ctx).Where("id=? and status=?", user.TenantId, consts.StatusOn).Take(&tenant)
	if tenant.ID == 0 {
		return "", errcode.ErrUserNotFound
	}

	us, _ := jsoniter.MarshalToString(map[string]int64{
		"uid":      user.ID,
		"tenantId": user.TenantId,
	})
	return utils.GenerateToken(us, utils.JwtSecret)
}

func GetLoginUser(ctx *gin.Context, uid int64) (*domain.LoginUser, error) {
	var user repo.User
	repo.GetDBWithTenant(ctx).Where("id=?", uid).Take(&user)
	if user.ID > 0 {
		roles := user.GetRoleCodes(ctx)
		buttons := user.GetButtons(ctx)
		return &domain.LoginUser{
			Uid:      user.ID,
			Username: user.Username,
			Roles:    roles,
			Buttons:  buttons,
		}, nil
	}
	return nil, errcode.ErrUserNotFound
}
