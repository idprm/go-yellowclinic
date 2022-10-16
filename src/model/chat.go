package model

import "gorm.io/gorm"

type Chat struct {
	ID         uint64 `gorm:"primaryKey"`
	gorm.Model `json:"-"`
}
