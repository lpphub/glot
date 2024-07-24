package entity

import repo "glot/repository"

type Menu struct {
	repo.RouteMeta
	repo.BaseModel
	ID        int64        `json:"id"`
	ParentID  int64        `json:"parentId"`
	MenuType  int8         `json:"menuType"`
	MenuName  string       `json:"menuName"`
	RouteName string       `json:"routeName,omitempty"`
	RoutePath string       `json:"routePath,omitempty"`
	Component string       `json:"component,omitempty"`
	Status    int8         `json:"status"`
	Buttons   []MenuButton `json:"buttons,omitempty"`
	Children  []Menu       `json:"children,omitempty"`
}

type MenuButton struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
	Code  string `json:"code"`
	Desc  string `json:"desc"`
}

type MenuTree struct {
	ID       int64       `json:"id"`
	PID      int64       `json:"pId"`
	Label    string      `json:"label"`
	Children []*MenuTree `json:"children,omitempty"`
}

type UserRoute struct {
	Routes []*RouteMenu `json:"routes"`
	Home   string       `json:"home"`
}

type RouteMenu struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	Path      string         `json:"path"`
	Component string         `json:"component,omitempty"`
	Children  []*RouteMenu   `json:"children,omitempty"`
	Meta      repo.RouteMeta `json:"meta"`
}
