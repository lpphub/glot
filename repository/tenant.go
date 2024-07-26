package repository

type Tenant struct {
	BaseModel
	ID       int64  `gorm:"id" json:"id"`
	Code     string `gorm:"code" json:"code"`
	Name     string `gorm:"name" json:"name"`
	Contacts string `gorm:"contacts" json:"contacts"`
	Phone    string `gorm:"phone" json:"phone"`
	Address  string `gorm:"address" json:"address"`
	Size     string `gorm:"size" json:"size"`
	Status   int8   `gorm:"status" json:"status"`
}

func (Tenant) TableName() string {
	return "tb_tenant"
}
