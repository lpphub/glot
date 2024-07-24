package system

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
	"glot/component/errcode"
	"glot/helper"
	repo "glot/repository"
	"glot/service/consts"
	"glot/service/entity"
	"gorm.io/gorm"
)

func PageListMenu(ctx *gin.Context, param entity.PageQuery) (*entity.Pager, error) {
	var (
		total int64
		list  []repo.Resource
	)
	_db := helper.DB.WithContext(ctx).Model(repo.Resource{}).Where("parent_id = 0 and resource_type in ?", consts.MenuRoute)

	if err := _db.Count(&total).Error; err != nil {
		return nil, err
	}
	if total > 0 {
		_db.Scopes(repo.Paginate(param.Pn, param.Ps)).Find(&list)

		voList := make([]entity.Menu, 0, len(list))
		for _, rsc := range list {
			rscTree := rsc.GetResourceTree(ctx)
			voList = append(voList, convertMenu(ctx, rscTree))
		}
		return entity.WrapPager(total, voList), nil
	}
	return entity.WrapPager(total, entity.EmptyList{}), nil
}

func convertMenu(ctx *gin.Context, tree *repo.ResourceTree) entity.Menu {
	vo := entity.Menu{
		ID:        tree.ID,
		ParentID:  tree.ParentID,
		MenuType:  tree.ResourceType,
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

	buttons := make([]entity.MenuButton, 0)
	for i, child := range children {
		if child.ResourceType == consts.ResourceButton {
			buttons = append(buttons, entity.MenuButton{
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

func GetMenuTree(ctx *gin.Context) ([]*entity.MenuTree, error) {
	var list []repo.Resource
	helper.DB.WithContext(ctx).Model(repo.Resource{}).Where("resource_type in ? and status =?",
		consts.MenuRoute, consts.StatusOn).Order("sort, id").Find(&list)

	menuMap := make(map[int64]*entity.MenuTree)
	for _, menu := range list {
		menuMap[menu.ID] = &entity.MenuTree{
			ID:    menu.ID,
			PID:   menu.ParentID,
			Label: menu.Name,
		}
	}
	menuTree := make([]*entity.MenuTree, 0)
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

func GetMenuButton(ctx *gin.Context) ([]entity.MenuButton, error) {
	var list []repo.Resource
	helper.DB.WithContext(ctx).Model(repo.Resource{}).Where("resource_type =? and status =?",
		consts.ResourceButton, consts.StatusOn).Find(&list)

	buttons := make([]entity.MenuButton, 0)
	for _, node := range list {
		buttons = append(buttons, entity.MenuButton{
			ID:    node.ID,
			Label: node.Code,
			Code:  node.Code,
		})
	}
	return buttons, nil
}

func SaveMenu(ctx *gin.Context, param entity.Menu) error {
	// todo 参数校验
	// 开启事务
	return helper.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var (
			meta       = param.RouteMeta
			metaStr, _ = jsoniter.MarshalToString(meta)

			base = repo.Resource{
				ID:           param.ID,
				ParentID:     param.ParentID,
				ResourceType: param.MenuType,
				Name:         param.MenuName,
				RoutePath:    param.RoutePath,
				RouteName:    param.RouteName,
				Component:    param.Component,
				Code:         param.RouteName,
				Status:       param.Status,
				Sort:         param.Order,
				Meta:         metaStr,
			}
		)

		if base.ID > 0 {
			base.FitUpdated(ctx)
			if err := tx.Model(&repo.Resource{ID: base.ID}).Updates(base).Error; err != nil {
				return err
			}
		} else {
			base.FitCreated(ctx)
			if err := tx.Create(&base).Error; err != nil {
				return err
			}
		}

		// 菜单下的按钮
		if len(param.Buttons) == 0 {
			var btnIds []int64
			tx.Model(&repo.Resource{}).Where("parent_id =? and resource_type =?", base.ID, consts.ResourceButton).
				Pluck("id", &btnIds)
			if len(btnIds) > 0 {
				tx.Delete(&repo.Resource{}, "id in ?", btnIds)
				tx.Delete(&repo.RoleResource{}, "resource_id in ?", btnIds)
			}
		} else {
			btnList := make([]*repo.Resource, 0, len(param.Buttons))
			for _, btn := range param.Buttons {
				var count int64
				tx.Model(&repo.Resource{}).Where("code =? and parent_id =? and resource_type =?", btn.Code,
					base.ID, consts.ResourceButton).Count(&count)
				if count == 0 {
					mbtn := &repo.Resource{
						ResourceType: consts.ResourceButton,
						ParentID:     base.ID,
						Name:         btn.Desc,
						Code:         btn.Code,
						Status:       consts.StatusOn,
					}
					mbtn.FitCreated(ctx)
					btnList = append(btnList, mbtn)
				}
			}
			return tx.Create(btnList).Error
		}
		return nil
	})
}

func DelMenu(ctx *gin.Context, ids []int64) error {
	var count int64
	helper.DB.WithContext(ctx).Model(&repo.Resource{}).Where("parent_id in ? and resource_type !=?", ids, consts.ResourceButton).Count(&count)
	if count > 0 { // 避免递归删除
		return errcode.ErrToast.Sprintf("请先删除选中项下的子菜单")
	}

	return helper.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&repo.Resource{}, "id in ?", ids).Error; err != nil {
			return err
		}
		if err := tx.Delete(&repo.RoleResource{}, "resource_id in ?", ids).Error; err != nil {
			return err
		}

		// 删除菜单下的按钮
		var btnIds []int64
		tx.Model(&repo.Resource{}).Where("resource_type =? and parent_id in ?", consts.ResourceButton, ids).Pluck("id", &btnIds)
		if len(btnIds) > 0 {

		}
		if err := tx.Delete(&repo.Resource{}, "id in ?", btnIds).Error; err != nil {
			return err
		}
		return tx.Delete(&repo.RoleResource{}, "resource_id in ?", btnIds).Error
	})
}

func IsExistRoute(ctx *gin.Context, routeName string) (bool, error) {
	var count int64
	helper.DB.WithContext(ctx).Model(&repo.Resource{}).Where("route_name = ?", routeName).Count(&count)
	return count > 0, nil
}

// GetUserRouteMenus 查询用户授权的菜单路由
func GetUserRouteMenus(ctx *gin.Context, uid int64) (*entity.UserRoute, error) {
	var user repo.User
	repo.DBWithTenant(ctx).Where("id=?", uid).Take(&user)
	if user.ID == 0 {
		return nil, errcode.ErrUserNotFound
	}

	var routes []repo.Resource
	helper.DB.WithContext(ctx).Model(&repo.Resource{}).Where("resource_type in ? and status =?",
		consts.MenuRoute, consts.StatusOn).Find(&routes)

	// 过滤角色限制
	var roleIds []int64
	helper.DB.WithContext(ctx).Model(&repo.UserRole{}).Where("user_id=?", uid).Pluck("role_id", &roleIds)
	filterRoutes := filterRoutesRole(ctx, routes, roleIds)

	route := &entity.UserRoute{
		Home:   "home",
		Routes: processRoutes(filterRoutes),
	}
	return route, nil
}

func filterRoutesRole(ctx *gin.Context, routes []repo.Resource, roleIds []int64) []repo.Resource {
	ids := lo.Map(routes, func(item repo.Resource, index int) int64 {
		return item.ID
	})
	var rrList []repo.RoleResource
	helper.DB.WithContext(ctx).Model(&repo.RoleResource{}).Where("resource_id in ?", ids).Find(&rrList)

	if len(rrList) > 0 {
		// 菜单对应的角色限制
		var rrMap = make(map[int64][]int64)
		for _, rr := range rrList {
			if _, ok := rrMap[rr.ResourceID]; ok {
				rrMap[rr.ResourceID] = append(rrMap[rr.ResourceID], rr.RoleID)
			} else {
				rrMap[rr.ResourceID] = []int64{rr.RoleID}
			}
		}

		filterRoutes := make([]repo.Resource, 0, len(routes))
		for _, route := range routes {
			if rrIds, ok := rrMap[route.ID]; ok {
				if lo.Some(rrIds, roleIds) {
					filterRoutes = append(filterRoutes, route)
				}
			} else { // 无角色限制
				filterRoutes = append(filterRoutes, route)
			}
		}
		return filterRoutes
	}
	return routes
}

func processRoutes(routes []repo.Resource) []*entity.RouteMenu {
	routeMap := make(map[int64]*entity.RouteMenu)
	for _, route := range routes {
		var meta repo.RouteMeta
		_ = jsoniter.UnmarshalFromString(route.Meta, &meta)

		routeMap[route.ID] = &entity.RouteMenu{
			ID:        route.ID,
			Name:      route.RouteName,
			Path:      route.RoutePath,
			Component: route.Component,
			Meta:      meta,
		}
	}
	rmList := make([]*entity.RouteMenu, 0)
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
