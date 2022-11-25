package model

import "gorm.io/gorm"

type User struct {
	ID            uint64 `gorm:"primaryKey" json:"id"`
	Msisdn        string `gorm:"size:15;not null" json:"msisdn"`
	LatestVoucher string `gorm:"size:100;not null" json:"latest_voucher"`
	Name          string `gorm:"size:200;not null" json:"name"`
	IpAddress     string `gorm:"size:15" json:"ip_address"`
	gorm.Model    `json:"-"`
}
