package model

import "gorm.io/gorm"

type Order struct {
	ID         uint64  `gorm:"primarykey" json:"id"`
	UserID     uint64  `json:"-"`
	User       User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	DoctorID   uint    `json:"-"`
	Doctor     Doctor  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Voucher    string  `gorm:"size:100;not null" json:"voucher"`
	Number     string  `gorm:"size:100" json:"number"`
	Total      float32 `gorm:"size:7" json:"total"`
	Status     string  `gorm:"size:25" json:"status"`
	gorm.Model `json:"-"`
}
