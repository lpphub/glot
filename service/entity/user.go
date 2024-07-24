package entity

import repo "glot/repository"

type UserQuery struct {
	PageQuery
	Uid      int64  `json:"uid" form:"uid"`
	Username string `json:"username" form:"username"`
	Nickname string `json:"nickname" form:"nickname"`
	Phone    string `json:"phone" form:"phone"`
	Email    string `json:"email" form:"email"`
	Status   int8   `json:"status" form:"status"`
}

type User struct {
	repo.User
	Roles []string `json:"roles"`
}

type LoginUser struct {
	Uid      int64    `json:"uid"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Buttons  []string `json:"buttons"`
}

type RoleQuery struct {
	PageQuery
	Name   string `json:"name" form:"name"`
	Code   string `json:"code" form:"code"`
	Status int8   `json:"status" form:"status"`
}
