package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"glot/component/utils"
	"glot/helper"
	"glot/middleware"
	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedBy string           `gorm:"created_by" json:"createdBy"`
	CreatedAt *utils.Timestamp `gorm:"created_at" json:"createdAt"`
	UpdatedBy string           `gorm:"updated_by" json:"updatedBy"`
	UpdatedAt *utils.Timestamp `gorm:"updated_at" json:"updatedAt"`
}

func (model *BaseModel) FitCreated(ctx *gin.Context) {
	model.CreatedBy = cast.ToString(middleware.GetLoginUid(ctx))
	model.CreatedAt = utils.NowTimestamp()
}
func (model *BaseModel) FitUpdated(ctx *gin.Context) {
	model.UpdatedBy = cast.ToString(middleware.GetLoginUid(ctx))
	model.UpdatedAt = utils.NowTimestamp()
}

func Paginate(pn, ps int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pn <= 0 {
			pn = 1
		}
		if ps <= 0 || ps > 200 {
			ps = 20
		}
		offset := (pn - 1) * ps
		return db.Offset(offset).Limit(ps)
	}
}

func GetDB(ctx *gin.Context) *gorm.DB {
	return helper.DB.WithContext(ctx)
}

func GetDBWithTenant(ctx *gin.Context) *gorm.DB {
	return GetDB(ctx).Scopes(ScopeTenant(ctx))
}

func ScopeTenant(ctx *gin.Context) func(db *gorm.DB) *gorm.DB {
	tId := middleware.GetLoginTenantId(ctx)
	return ScopeTenantId(tId)
}

func ScopeTenantId(tenantId int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("tenant_id=?", tenantId)
	}
}
