package tenant

import (
	"github.com/gin-gonic/gin"
	repo "glot/repository"
	"glot/service/consts"
	"glot/service/domain"
	"gorm.io/gorm"
)

func PageList(ctx *gin.Context, param domain.TenantQuery) (*domain.Pager, error) {
	var (
		total int64
		list  []repo.Tenant
	)
	_db := repo.GetDB(ctx).Model(repo.Tenant{})

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

		voList := make([]*domain.TenantVO, 0)
		for _, tenant := range list {
			roles := tenant.GetRoleCodes(ctx)
			voList = append(voList, &domain.TenantVO{
				Tenant: tenant,
				Roles:  roles,
			})
		}
		return domain.WrapPager(total, voList), nil
	}
	return domain.WrapPager(total, domain.EmptyList{}), nil
}

func SaveTenant(ctx *gin.Context, param domain.TenantVO) error {
	var (
		tenant    = param.Tenant
		roleCodes = param.Roles
	)
	return repo.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		if tenant.ID > 0 {
			tenant.FitUpdated(ctx)
			if err := tx.Updates(tenant).Error; err != nil {
				return err
			}
		} else {
			tenant.FitCreated(ctx)
			if err := tx.Create(&tenant).Error; err != nil {
				return err
			}
		}

		// 绑定租户可分配角色
		tx.Delete(&repo.TenantRole{}, "tenant_id =?", tenant.ID)
		if len(roleCodes) > 0 {
			var roleIds []int64
			tx.Model(&repo.Role{}).Where("code in ?", roleCodes).Select("id").Find(&roleIds)

			trList := make([]*repo.TenantRole, 0)
			for _, id := range roleIds {
				trList = append(trList, &repo.TenantRole{
					TenantID: tenant.ID,
					RoleID:   id,
				})
			}
			return tx.Create(&trList).Error
		}
		return nil
	})
}

func ListRoleScope(ctx *gin.Context) ([]repo.Role, error) {
	var roleIds []int64
	repo.GetDBWithTenant(ctx).Model(&repo.TenantRole{}).Pluck("role_id", &roleIds)
	if len(roleIds) == 0 {
		return nil, nil
	}
	var list []repo.Role
	repo.GetDB(ctx).Model(repo.Role{}).Where("id in ? and status =?", roleIds, consts.StatusOn).Find(&list)
	return list, nil
}
