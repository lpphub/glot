package system

import (
	"github.com/gin-gonic/gin"
	"glot/middleware"
	repo "glot/repository"
	"glot/service/domain"
	"gorm.io/gorm"
)

func PageListUser(ctx *gin.Context, param domain.UserQuery) (*domain.Pager, error) {
	var (
		total int64
		list  []repo.User
	)
	_db := repo.GetDB(ctx).Model(repo.User{})

	if param.Uid > 0 {
		_db.Where("id = ?", param.Uid)
	}
	if param.Username != "" {
		_db.Where("username like ?", "%"+param.Username+"%")
	}
	if param.Nickname != "" {
		_db.Where("nickname like ?", "%"+param.Nickname+"%")
	}
	if param.Phone != "" {
		_db.Where("phone =?", param.Phone)
	}
	if param.Email != "" {
		_db.Where("email =?", param.Email)
	}
	if param.Status > 0 {
		_db.Where("status =?", param.Status)
	}

	if err := _db.Count(&total).Error; err != nil {
		return nil, err
	}
	if total > 0 {
		_db.Order("id desc").Scopes(repo.Paginate(param.Pn, param.Ps)).Find(&list)

		userList := make([]*domain.User, 0)
		for _, user := range list {
			roles := user.GetRoleCodes(ctx)
			userList = append(userList, &domain.User{
				User:  user,
				Roles: roles,
			})
		}
		return domain.WrapPager(total, userList), nil
	}
	return domain.WrapPager(total, domain.EmptyList{}), nil
}

func SaveUser(ctx *gin.Context, param domain.User) error {
	// 开启事务
	return repo.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		user := param.User
		if user.TenantId == 0 {
			user.TenantId = middleware.GetLoginTenantId(ctx)
		}
		if user.ID > 0 {
			user.FitUpdated(ctx)
			if err := tx.Omit("tenant_id").Updates(user).Error; err != nil {
				return err
			}
		} else {
			user.Password = "123456"
			user.FitCreated(ctx)
			if err := tx.Create(&user).Error; err != nil {
				return err
			}
		}

		// 绑定用户角色
		tx.Delete(&repo.UserRole{}, "user_id =?", user.ID)
		if len(param.Roles) > 0 {
			var roleIds []int64
			tx.Model(&repo.Role{}).Where("code in ?", param.Roles).Select("id").Find(&roleIds)

			urList := make([]*repo.UserRole, 0)
			for _, id := range roleIds {
				urList = append(urList, &repo.UserRole{
					UserID: user.ID,
					RoleID: id,
				})
			}
			return tx.Create(&urList).Error
		}
		return nil
	})
}

func DelUser(ctx *gin.Context, ids []int64) error {
	return repo.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&repo.User{}, "id in ?", ids).Error; err != nil {
			return err
		}
		return tx.Delete(&repo.UserRole{}, "user_id in ?", ids).Error
	})
}
