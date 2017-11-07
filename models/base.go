package models

import (
	"time"
)

// Model 用于gorm的基础model类型
type Model struct {
	ID        int `gorm:"AUTO_INCREMENT;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
