package repository

import (
	"github.com/gin-gonic/gin"
	"glot/service/consts"
)

type Tenant struct {
	BaseModel
	ID       int64  `gorm:"id" json:"id"`
	Code     string `gorm:"code" json:"code"`
	Name     string `gorm:"name" json:"name"`
	Contacts string `gorm:"contacts" json:"contacts"`
	Phone    string `gorm:"phone" json:"phone"`
	Address  string `gorm:"address" json:"address"`
	Size     string `gorm:"size" json:"size"`
	Status   int8   `gorm:"status" json:"status"`
}

func (*Tenant) TableName() string {
	return "tb_tenant"
}

func (t *Tenant) GetRoleCodes(ctx *gin.Context) (roles []string) {
	var roleIds []int64
	GetDB(ctx).Model(&TenantRole{}).Where("tenant_id=?", t.ID).Pluck("role_id", &roleIds)
	if len(roleIds) > 0 {
		GetDB(ctx).Model(Role{}).Where("id in ? and status=?", roleIds, consts.StatusOn).
			Pluck("code", &roles)
	}
	return
}
