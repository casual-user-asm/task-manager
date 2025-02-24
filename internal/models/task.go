package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title       string `json:"title" gorm:"not null"`
	Description string `'json:"description" gorm:"not null"`
	CreatedBy   uint   `json:"created_by" gorm:"not null;default:0"`
	User        User   `json:"user" gorm:"foreignKey:CreatedBy;references:id"`
	Date        time.Time
	Status      int `json:"status" gorm:"default:0"` // 0 = new, 1 = ongoing, 2 = completed
}
