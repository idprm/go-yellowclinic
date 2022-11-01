package model

import "gorm.io/gorm"

type User struct {
	ID         uint64 `gorm:"primaryKey" json:"id"`
	Msisdn     string `gorm:"size:15;unique;not null" json:"msisdn"`
	Name       string `gorm:"size:200;not null" json:"name"`
	UserIds    string `gorm:"size:25"  json:"user_id"`
	gorm.Model `json:"-"`
}
