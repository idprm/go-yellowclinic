package model

type Clinic struct {
	ID       uint64 `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"size:150;not null" json:"name"`
	Address  string `gorm:"type:text" json:"address"`
	Phone    string `gorm:"size:15" json:"phone"`
	IsActive bool   `gorm:"type:boolean" json:"is_acitve"`
}
