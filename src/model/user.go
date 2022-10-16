package model

type User struct {
	ID     uint64 `gorm:"primaryKey" json:"id"`
	Msisdn string
}
