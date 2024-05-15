package entity

import (
	"gorm.io/gorm"
	"time"
)

type Medicine struct {
	ID        uint           `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Name      string         `json:"name" validate:"required"`
	Amount    int            `json:"amount" validate:"required"`
	Expired   string         `json:"expired" validate:"required"`
	Image     string         `json:"image"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index,column:deleted_at"`
}

func (m *Medicine) TableName() string {
	return "medicines"
}
