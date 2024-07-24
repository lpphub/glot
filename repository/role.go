package repository

type Role struct {
	BaseModel
	ID     int64  `gorm:"id" json:"id"`
	Name   string `gorm:"name" json:"name"`
	Code   string `gorm:"code" json:"code"`
	Desc   string `gorm:"desc" json:"desc"`
	Status int8   `gorm:"status" json:"status"`
}

func (Role) TableName() string {
	return "tb_role"
}
