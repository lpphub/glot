package repository

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"glot/helper"
)

type Menu struct {
	BaseModel
	ID        int64  `gorm:"id" json:"id"`
	ParentID  int64  `gorm:"parent_id" json:"parentId"`
	Mode      int8   `gorm:"mode" json:"mode"`                      // 类型 1-目录 2-菜单 3-页面按钮
	Name      string `gorm:"name" json:"name"`                      // 资源名称
	Code      string `gorm:"code" json:"code"`                      // 资源唯一标识
	RouteName string `gorm:"route_name" json:"routeName,omitempty"` // 路由名称
	RoutePath string `gorm:"route_path" json:"routePath,omitempty"` // 路由路径
	Component string `gorm:"component" json:"component,omitempty"`  // 页面组件
	Meta      string `gorm:"meta" json:"-"`                         // 资源元数据
	Status    int8   `gorm:"status" json:"status"`                  // 状态 1-启用 2-禁用
	Sort      int    `gorm:"sort" json:"sort"`                      // 排序
}

type MenuTree struct {
	Menu
	RouteMeta
	Children []*MenuTree `json:"children"`
}

type RouteMeta struct {
	Title           string `json:"title,omitempty"`           //标题
	I18NKey         string `json:"i18nKey,omitempty"`         //国际化key
	IconType        int8   `json:"iconType"`                  //icon类型
	Icon            string `json:"icon"`                      //icon
	HideInMenu      bool   `json:"hideInMenu,omitempty"`      //是否在菜单中隐藏
	ActiveMenu      string `json:"activeMenu,omitempty"`      //激活的菜单键
	KeepAlive       bool   `json:"keepAlive,omitempty"`       //是否缓存
	Constant        bool   `json:"constant,omitempty"`        //常量路由
	Href            string `json:"href,omitempty"`            //外部链接
	MultiTab        bool   `json:"multiTab,omitempty"`        //是否开启多标签页模式
	FixedIndexInTab int    `json:"fixedIndexInTab,omitempty"` //在标签页中固定索引位置
	Order           int    `json:"order"`                     //排序
	Query           []KV   `json:"query,omitempty"`           //路由参数
	PathParam       string `json:"pathParam,omitempty"`       //路径参数
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (Menu) TableName() string {
	return "tb_menu"
}

func (r Menu) getChildren(ctx *gin.Context, resourceTypes ...int) (res []Menu) {
	_tx := helper.DB.WithContext(ctx).Model(&Menu{}).Where("parent_id =?", r.ID)
	if len(resourceTypes) > 0 {
		_tx.Where("mode in ?", resourceTypes)
	}
	_tx.Find(&res)
	return
}

func (r Menu) GetMenuTree(ctx *gin.Context) *MenuTree {
	var meta RouteMeta
	_ = jsoniter.UnmarshalFromString(r.Meta, &meta)
	tree := &MenuTree{
		Menu:      r,
		RouteMeta: meta,
	}
	// 递归子项
	children := r.getChildren(ctx)
	if len(children) > 0 {
		for _, child := range children {
			tree.Children = append(tree.Children, child.GetMenuTree(ctx))
		}
	}
	return tree
}
