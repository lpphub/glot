package system

import (
	"github.com/gin-gonic/gin"
	"glot/component/errcode"
	repo "glot/repository"
	"glot/service/consts"
	"glot/service/domain"
	"gorm.io/gorm"
)

func PageListRole(ctx *gin.Context, param domain.RoleQuery) (*domain.Pager, error) {
	var (
		total int64
		list  []repo.Role
	)
	_db := repo.GetDB(ctx).Model(repo.Role{})

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

func ListAllRole(ctx *gin.Context) ([]repo.Role, error) {
	var list []repo.Role
	repo.GetDB(ctx).Model(repo.Role{}).Where("status =?", consts.StatusOn).Find(&list)
	return list, nil
}

func SaveRole(ctx *gin.Context, role repo.Role) error {
	if role.ID > 0 {
		role.FitUpdated(ctx)
		return repo.GetDB(ctx).Updates(role).Error
	} else {
		role.FitCreated(ctx)
		return repo.GetDB(ctx).Create(&role).Error
	}
}

func DelRole(ctx *gin.Context, ids []int64) error {
	return repo.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&repo.Role{}, "id in ?", ids).Error; err != nil {
			return err
		}
		return tx.Delete(&repo.UserRole{}, "role_id in ?", ids).Error
	})
}

func GetRoleMenu(ctx *gin.Context, roleId int64, rscTypes []int) (ids []int64, err error) {
	var menuIds []int64
	repo.GetDB(ctx).Model(&repo.RoleMenu{}).Where("role_id =?", roleId).Pluck("menu_id", &menuIds)
	if len(menuIds) > 0 {
		repo.GetDB(ctx).Model(&repo.Menu{}).Where("id in ? and mode in ? and status=?",
			menuIds, rscTypes, consts.StatusOn).Pluck("id", &ids)
	}
	return
}

func BindRoleMenu(ctx *gin.Context, roleId int64, modes []int, ids []int64) error {
	var role repo.Role
	repo.GetDB(ctx).Where("id =?", roleId).Take(&role)
	if role.ID == 0 {
		return errcode.ErrParamInvalid
	}

	return repo.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		//删掉原记录重新保存
		delSql := "DELETE rm FROM tb_role_menu rm JOIN tb_menu m ON rm.menu_id = m.id " +
			"WHERE rm.role_id =? and m.mode in ?"
		err := tx.Exec(delSql, roleId, modes).Error
		if err != nil {
			return err
		}

		if len(ids) == 0 {
			return nil
		}
		rrs := make([]repo.RoleMenu, 0, len(ids))
		for _, id := range ids {
			rrs = append(rrs, repo.RoleMenu{
				RoleID: roleId,
				MenuID: id,
			})
		}
		return tx.Create(&rrs).Error
	})
}
