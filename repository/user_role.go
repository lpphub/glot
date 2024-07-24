package repository

type UserRole struct {
	ID     int64 `gorm:"id"`
	UserID int64 `gorm:"user_id"`
	RoleID int64 `gorm:"role_id"`
}

func (UserRole) TableName() string {
	return "tb_user_role"
}
