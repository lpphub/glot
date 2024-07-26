package domain

type TenantQuery struct {
	PageQuery
	Code   string `json:"code" form:"code"`
	Name   string `json:"name" form:"name"`
	Status int8   `json:"status" form:"status"`
}
