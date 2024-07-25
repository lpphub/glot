package domain

import repo "glot/repository"

type MenuVO struct {
	repo.RouteMeta
	repo.BaseModel
	ID        int64           `json:"id"`
	ParentID  int64           `json:"parentId"`
	MenuType  int8            `json:"menuType"`
	MenuName  string          `json:"menuName"`
	RouteName string          `json:"routeName,omitempty"`
	RoutePath string          `json:"routePath,omitempty"`
	Component string          `json:"component,omitempty"`
	Status    int8            `json:"status"`
	Buttons   []*MenuButtonVO `json:"buttons,omitempty"`
	Children  []*MenuVO       `json:"children,omitempty"`
}

type MenuButtonVO struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
	Code  string `json:"code"`
	Desc  string `json:"desc"`
}

type MenuTreeVO struct {
	ID       int64         `json:"id"`
	PID      int64         `json:"pId"`
	Label    string        `json:"label"`
	Children []*MenuTreeVO `json:"children,omitempty"`
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
