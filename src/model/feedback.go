package model

import "gorm.io/gorm"

type Feedback struct {
	ID         uint   `gorm:"primaryKey"`
	DoctorID   uint   `json:"-"`
	Doctor     Doctor `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID     uint64 `json:"user_id"`
	User       User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Message    string `gorm:"type:text" json:"message"`
	gorm.Model `json:"-"`
}
