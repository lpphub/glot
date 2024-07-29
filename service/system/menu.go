package system

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
	"glot/component/errcode"
	repo "glot/repository"
	"glot/service/consts"
	"glot/service/domain"
	"gorm.io/gorm"
)

func PageListMenu(ctx *gin.Context, param domain.PageQuery) (*domain.Pager, error) {
	var (
		total int64
		list  []repo.Menu
	)
	_db := repo.GetDB(ctx).Model(repo.Menu{}).Where("parent_id = 0 and mode in ?", consts.RouteMenu)

	if err := _db.Count(&total).Error; err != nil {
		return nil, err
	}
	if total > 0 {
		_db.Scopes(repo.Paginate(param.Pn, param.Ps)).Find(&list)

		voList := make([]*domain.MenuVO, 0, len(list))
		for _, rsc := range list {
			rscTree := rsc.GetMenuTree(ctx)
			voList = append(voList, convertMenu(ctx, rscTree))
		}
		return domain.WrapPager(total, voList), nil
	}
	return domain.WrapPager(total, domain.EmptyList{}), nil
}

func convertMenu(ctx *gin.Context, tree *repo.MenuTree) *domain.MenuVO {
	vo := &domain.MenuVO{
		ID:        tree.ID,
		ParentID:  tree.ParentID,
		MenuType:  tree.Mode,
		MenuName:  tree.Name,
		RouteName: tree.RouteName,
		RoutePath: tree.RoutePath,
		Component: tree.Component,
		Status:    tree.Status,
		RouteMeta: tree.RouteMeta,
		BaseModel: tree.BaseModel,
	}

	children := tree.Children
	if len(children) == 0 {
		return vo
	}

	buttons := make([]*domain.MenuButtonVO, 0)
	for i, child := range children {
		if child.Mode == consts.MenuButton {
			buttons = append(buttons, &domain.MenuButtonVO{
				ID:    child.ID,
				Code:  child.Code,
				Label: child.Name,
				Desc:  child.Name,
			})
		} else {
			vo.Children = append(vo.Children, convertMenu(ctx, children[i]))
		}
	}
	vo.Buttons = buttons
	return vo
}

func GetMenuTree(ctx *gin.Context) ([]*domain.MenuTreeVO, error) {
	var list []repo.Menu
	repo.GetDB(ctx).Model(repo.Menu{}).Where("mode in ? and status =?",
		consts.RouteMenu, consts.StatusOn).Order("sort, id").Find(&list)

	menuMap := make(map[int64]*domain.MenuTreeVO)
	for _, menu := range list {
		menuMap[menu.ID] = &domain.MenuTreeVO{
			ID:    menu.ID,
			PID:   menu.ParentID,
			Label: menu.Name,
		}
	}
	menuTree := make([]*domain.MenuTreeVO, 0)
	for _, node := range list {
		if node.ParentID == 0 {
			menuTree = append(menuTree, menuMap[node.ID])
		} else {
			parent := menuMap[node.ParentID]
			parent.Children = append(parent.Children, menuMap[node.ID])
		}
	}
	return menuTree, nil
}

func GetMenuButton(ctx *gin.Context) ([]*domain.MenuButtonVO, error) {
	var list []repo.Menu
	repo.GetDB(ctx).Model(repo.Menu{}).Where("mode =? and status =?",
		consts.MenuButton, consts.StatusOn).Find(&list)

	buttons := make([]*domain.MenuButtonVO, 0)
	for _, node := range list {
		buttons = append(buttons, &domain.MenuButtonVO{
			ID:    node.ID,
			Label: node.Code,
			Code:  node.Code,
		})
	}
	return buttons, nil
}

func SaveMenu(ctx *gin.Context, param domain.MenuVO) error {
	// todo 参数校验
	// 开启事务
	return repo.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		var (
			meta       = param.RouteMeta
			metaStr, _ = jsoniter.MarshalToString(meta)

			base = repo.Menu{
				ID:        param.ID,
				ParentID:  param.ParentID,
				Mode:      param.MenuType,
				Name:      param.MenuName,
				RoutePath: param.RoutePath,
				RouteName: param.RouteName,
				Component: param.Component,
				Code:      param.RouteName,
				Status:    param.Status,
				Sort:      param.Order,
				Meta:      metaStr,
			}
		)

		if base.ID > 0 {
			base.FitUpdated(ctx)
			if err := tx.Model(&repo.Menu{ID: base.ID}).Updates(base).Error; err != nil {
				return err
			}
		} else {
			base.FitCreated(ctx)
			if err := tx.Create(&base).Error; err != nil {
				return err
			}
		}

		// 菜单下的按钮
		var btnIds []int64
		tx.Model(&repo.Menu{}).Where("parent_id =? and mode =?", base.ID, consts.MenuButton).Pluck("id", &btnIds)
		if len(btnIds) > 0 {
			tx.Delete(&repo.Menu{}, "id in ?", btnIds)
			tx.Delete(&repo.RoleMenu{}, "menu_id in ?", btnIds)
		}
		if len(param.Buttons) > 0 {
			btnList := make([]*repo.Menu, 0, len(param.Buttons))
			for _, btn := range param.Buttons {
				var count int64
				tx.Model(&repo.Menu{}).Where("code =? and parent_id =? and mode =?", btn.Code,
					base.ID, consts.MenuButton).Count(&count)
				if count == 0 {
					mbtn := &repo.Menu{
						Mode:     consts.MenuButton,
						ParentID: base.ID,
						Name:     btn.Desc,
						Code:     btn.Code,
						Status:   consts.StatusOn,
					}
					mbtn.FitCreated(ctx)
					btnList = append(btnList, mbtn)
				}
			}
			return tx.Create(&btnList).Error
		}
		return nil
	})
}

func DelMenu(ctx *gin.Context, ids []int64) error {
	var count int64
	repo.GetDB(ctx).Model(&repo.Menu{}).Where("parent_id in ? and mode in ?", ids, consts.RouteMenu).Count(&count)
	if count > 0 { // 避免递归删除
		return errcode.ErrToast.Sprintf("请先删除选中项下的子菜单")
	}

	return repo.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&repo.Menu{}, "id in ?", ids).Error; err != nil {
			return err
		}
		if err := tx.Delete(&repo.RoleMenu{}, "menu_id in ?", ids).Error; err != nil {
			return err
		}

		// 删除菜单下的按钮
		var btnIds []int64
		tx.Model(&repo.Menu{}).Where("mode =? and parent_id in ?", consts.MenuButton, ids).Pluck("id", &btnIds)
		if len(btnIds) > 0 {

		}
		if err := tx.Delete(&repo.Menu{}, "id in ?", btnIds).Error; err != nil {
			return err
		}
		return tx.Delete(&repo.RoleMenu{}, "menu_id in ?", btnIds).Error
	})
}

func IsExistRoute(ctx *gin.Context, routeName string) (bool, error) {
	var count int64
	repo.GetDB(ctx).Model(&repo.Menu{}).Where("route_name = ?", routeName).Count(&count)
	return count > 0, nil
}

// GetUserRouteMenus 查询用户授权的菜单路由
func GetUserRouteMenus(ctx *gin.Context, uid int64) (*domain.UserRoute, error) {
	var user repo.User
	repo.GetDBWithTenant(ctx).Where("id=?", uid).Take(&user)
	if user.ID == 0 {
		return nil, errcode.ErrUserNotFound
	}

	var routes []repo.Menu
	repo.GetDB(ctx).Model(&repo.Menu{}).Where("mode in ? and status =?", consts.RouteMenu, consts.StatusOn).Find(&routes)

	// 过滤角色限制
	var roleIds []int64
	repo.GetDB(ctx).Model(&repo.UserRole{}).Where("user_id=?", uid).Pluck("role_id", &roleIds)
	filterRoutes := filterRoutesRole(ctx, routes, roleIds)

	route := &domain.UserRoute{
		Home:   "home",
		Routes: processRoutes(filterRoutes),
	}
	return route, nil
}

func filterRoutesRole(ctx *gin.Context, routes []repo.Menu, roleIds []int64) []repo.Menu {
	ids := lo.Map(routes, func(item repo.Menu, index int) int64 {
		return item.ID
	})
	var rrList []repo.RoleMenu
	repo.GetDB(ctx).Model(&repo.RoleMenu{}).Where("menu_id in ?", ids).Find(&rrList)

	if len(rrList) > 0 {
		// 菜单对应的角色限制
		var rrMap = make(map[int64][]int64)
		for _, rr := range rrList {
			if _, ok := rrMap[rr.MenuID]; ok {
				rrMap[rr.MenuID] = append(rrMap[rr.MenuID], rr.RoleID)
			} else {
				rrMap[rr.MenuID] = []int64{rr.RoleID}
			}
		}

		filterRoutes := make([]repo.Menu, 0, len(routes))
		for _, route := range routes {
			if rrIds, ok := rrMap[route.ID]; ok {
				if lo.Some(rrIds, roleIds) {
					filterRoutes = append(filterRoutes, route)
				}
			}
		}
		return filterRoutes
	}
	return routes
}

func processRoutes(routes []repo.Menu) []*domain.RouteMenu {
	routeMap := make(map[int64]*domain.RouteMenu)
	for _, route := range routes {
		var meta repo.RouteMeta
		_ = jsoniter.UnmarshalFromString(route.Meta, &meta)

		routeMap[route.ID] = &domain.RouteMenu{
			ID:        route.ID,
			Name:      route.RouteName,
			Path:      route.RoutePath,
			Component: route.Component,
			Meta:      meta,
		}
	}
	rmList := make([]*domain.RouteMenu, 0)
	for _, r := range routes {
		if r.ParentID == 0 {
			rmList = append(rmList, routeMap[r.ID])
		} else {
			if parent, ok := routeMap[r.ParentID]; ok {
				parent.Children = append(parent.Children, routeMap[r.ID])
			} else {
				rmList = append(rmList, routeMap[r.ID])
			}
		}
	}
	return rmList
}
