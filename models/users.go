package models

import "time"

type Users struct {
	ID        uint   `gorm:"autoIncrement"`
	Username  string `gorm:"not null;unique"`
	Password  string `gorm:"not null" json:"-" `
	CreatedAt time.Time
}
