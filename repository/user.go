package repository

import (
	"github.com/gin-gonic/gin"
	"glot/helper"
	"glot/service/consts"
)

type User struct {
	BaseModel
	ID       int64  `gorm:"id" json:"id"`
	Username string `gorm:"username" json:"username"`
	Password string `gorm:"password" json:"-"`
	Phone    string `gorm:"phone" json:"phone"`
	Email    string `gorm:"email" json:"email"`
	Nickname string `gorm:"nickname" json:"nickname"`
	Status   int8   `gorm:"status" json:"status"`
	TenantId int64  `gorm:"tenant_id" json:"-"`
}

func (User) TableName() string {
	return "tb_user"
}

func (u User) GetRoleCodes(ctx *gin.Context) (roles []string) {
	var roleIds []int64
	helper.DB.WithContext(ctx).Model(&UserRole{}).Where("user_id=?", u.ID).Pluck("role_id", &roleIds)
	if len(roleIds) > 0 {
		helper.DB.WithContext(ctx).Model(Role{}).Where("id in ? and status=?", roleIds, consts.StatusOn).
			Pluck("code", &roles)
	}
	return
}

func (u User) GetButtons(ctx *gin.Context) (codes []string) {
	var roleIds []int64
	helper.DB.WithContext(ctx).Model(&UserRole{}).Where("user_id=?", u.ID).Pluck("role_id", &roleIds)
	if len(roleIds) == 0 {
		return
	}

	var menuIds []int64
	helper.DB.WithContext(ctx).Model(&RoleMenu{}).Where("role_id in ?", roleIds).Pluck("menu_id", &menuIds)
	if len(menuIds) > 0 {
		helper.DB.WithContext(ctx).Model(&Menu{}).Where("id in ? and mode = ? and status=?", menuIds,
			consts.RouteMenu, consts.StatusOn).Pluck("code", &codes)
	}
	return
}
