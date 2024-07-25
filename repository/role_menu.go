package repository

type RoleMenu struct {
	ID     int64 `gorm:"id"`
	RoleID int64 `gorm:"role_id"`
	MenuID int64 `gorm:"menu_id"`
}

func (RoleMenu) TableName() string {
	return "tb_role_menu"
}
