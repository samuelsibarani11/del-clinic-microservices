package entity

import (
	"gorm.io/gorm"
	"time"
)

type Dorm struct {
	ID        uint           `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Name      string         `json:"name"`
	Status    string         `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index,column:deleted_at"`
}

func (m *Dorm) TableName() string {
	return "dorms"
}
