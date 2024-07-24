package repository

import "glot/component/utils"

type Tenant struct {
	ID        int64            `gorm:"id" json:"id"`
	Name      string           `gorm:"name" json:"name"`
	Contact   string           `gorm:"contact" json:"contact"`
	Phone     string           `gorm:"phone" json:"phone"`
	Address   string           `gorm:"address" json:"address"`
	Status    int8             `gorm:"status" json:"status"`
	CreatedAt *utils.Timestamp `gorm:"created_at" json:"createdAt"`
	UpdatedAt *utils.Timestamp `gorm:"updated_at" json:"updatedAt"`
}

func (Tenant) TableName() string {
	return "tb_tenant"
}
