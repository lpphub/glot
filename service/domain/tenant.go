package domain

import repo "glot/repository"

type TenantQuery struct {
	PageQuery
	Code   string `json:"code" form:"code"`
	Name   string `json:"name" form:"name"`
	Status int8   `json:"status" form:"status"`
}

type TenantVO struct {
	repo.Tenant
	Roles []string `json:"roles"`
}
