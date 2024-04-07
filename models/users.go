package models

import "time"

type Users struct {
	ID        uint   `gorm:"autoIncrement"`
	Name      string `gorm:"not null"`
	CreatedAt time.Time
}
