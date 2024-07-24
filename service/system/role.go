package system

import (
	"github.com/gin-gonic/gin"
	"glot/component/errcode"
	"glot/helper"
	repo "glot/repository"
	"glot/service/consts"
	"glot/service/entity"
	"gorm.io/gorm"
)

func PageListRole(ctx *gin.Context, param entity.RoleQuery) (*entity.Pager, error) {
	var (
		total int64
		list  []repo.Role
	)
	_db := helper.DB.WithContext(ctx).Model(repo.Role{})

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
		return entity.WrapPager(total, list), nil
	}
	return entity.WrapPager(total, entity.EmptyList{}), nil
}

func ListAllRole(ctx *gin.Context) ([]repo.Role, error) {
	var list []repo.Role
	helper.DB.WithContext(ctx).Model(repo.Role{}).Where("status =?", consts.StatusOn).Find(&list)
	return list, nil
}

func SaveRole(ctx *gin.Context, role repo.Role) error {
	if role.ID > 0 {
		role.FitUpdated(ctx)
		return helper.DB.WithContext(ctx).Updates(role).Error
	} else {
		role.FitCreated(ctx)
		return helper.DB.WithContext(ctx).Create(&role).Error
	}
}

func DelRole(ctx *gin.Context, ids []int64) error {
	return helper.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&repo.Role{}, "id in ?", ids).Error; err != nil {
			return err
		}
		return tx.Delete(&repo.UserRole{}, "role_id in ?", ids).Error
	})
}

func GetRoleResource(ctx *gin.Context, roleId int64, rscTypes []int) (ids []int64, err error) {
	var rscIds []int64
	helper.DB.WithContext(ctx).Model(&repo.RoleResource{}).Where("role_id =?", roleId).Pluck("resource_id", &rscIds)
	if len(rscIds) > 0 {
		helper.DB.WithContext(ctx).Model(&repo.Resource{}).Where("id in ? and resource_type in ? and status=?",
			rscIds, rscTypes, consts.StatusOn).Pluck("id", &ids)
	}
	return
}

func BindRoleResource(ctx *gin.Context, roleId int64, rscTypes []int, ids []int64) error {
	var role repo.Role
	helper.DB.WithContext(ctx).Where("id =?", roleId).Take(&role)
	if role.ID == 0 {
		return errcode.ErrParamInvalid
	}

	return helper.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//删掉原记录重新保存
		delSql := "DELETE rr FROM tb_role_resource rr JOIN tb_resource r ON rr.resource_id = r.id " +
			"WHERE role_id =? and r.resource_type in ?"
		err := tx.Exec(delSql, roleId, rscTypes).Error
		if err != nil {
			return err
		}

		if len(ids) == 0 {
			return nil
		}
		rrs := make([]repo.RoleResource, 0, len(ids))
		for _, id := range ids {
			rrs = append(rrs, repo.RoleResource{
				RoleID:     roleId,
				ResourceID: id,
			})
		}
		return tx.Create(&rrs).Error
	})
}
