package repository

type TenantRole struct {
	ID       int64 `gorm:"id"`
	TenantID int64 `gorm:"tenant_id"`
	RoleID   int64 `gorm:"role_id"`
}

func (*TenantRole) TableName() string {
	return "tb_tenant_role"
}
