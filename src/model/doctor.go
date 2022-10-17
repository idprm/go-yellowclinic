package model

import (
	"time"

	"gorm.io/gorm"
)

type Doctor struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	Name                 string    `gorm:"size:100;not null" json:"name"`
	Photo                string    `gorm:"size:150;not null" json:"photo"`
	Number               string    `gorm:"size:100" json:"number"`
	Experience           int       `gorm:"size:2" json:"experience"`
	GraduatedFrom        string    `gorm:"size:150" json:"graduated_from"`
	ConsultationSchedule string    `gorm:"size:250" json:"consultation_schedule"`
	PlacePractice        string    `gorm:"size:250" json:"place_practice"`
	Phone                string    `gorm:"size:15" json:"phone"`
	Username             string    `gorm:"size:25" json:"username"`
	Start                time.Time `gorm:"type:time" json:"start"`
	End                  time.Time `gorm:"type:time" json:"end"`
	gorm.Model           `json:"-"`
}
