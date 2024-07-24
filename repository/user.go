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

	var rscIds []int64
	helper.DB.WithContext(ctx).Model(&RoleResource{}).Where("role_id in ?", roleIds).Pluck("resource_id", &rscIds)
	if len(rscIds) > 0 {
		helper.DB.WithContext(ctx).Model(&Resource{}).Where("id in ? and resource_type = ? and status=?", rscIds,
			consts.ResourceButton, consts.StatusOn).Pluck("code", &codes)
	}
	return
}
