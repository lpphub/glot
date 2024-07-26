package tenant

import (
	"github.com/gin-gonic/gin"
	"glot/helper"
	repo "glot/repository"
	"glot/service/domain"
)

func PageList(ctx *gin.Context, param domain.TenantQuery) (*domain.Pager, error) {
	var (
		total int64
		list  []repo.Tenant
	)
	_db := helper.DB.WithContext(ctx).Model(repo.Tenant{})

	if param.Name != "" {
		_db.Where("name like ?", "%"+param.Name+"%")
	}
	if param.Code != "" {
		_db.Where("code =?", param.Code)
	}
	if param.Status > 0 {
		_db.Where("status =?", param.Status)
	}

	if err := _db.Count(&total).Error; err != nil {
		return nil, err
	}
	if total > 0 {
		_db.Order("id desc").Scopes(repo.Paginate(param.Pn, param.Ps)).Find(&list)
		return domain.WrapPager(total, list), nil
	}
	return domain.WrapPager(total, domain.EmptyList{}), nil
}
