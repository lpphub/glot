package repository

type RoleResource struct {
	ID         int64 `gorm:"id"`
	RoleID     int64 `gorm:"role_id"`
	ResourceID int64 `gorm:"resource_id"`
}

func (RoleResource) TableName() string {
	return "tb_role_resource"
}
