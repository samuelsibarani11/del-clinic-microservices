package entity

import (
	"gorm.io/gorm"
	"time"
)

type Reminder struct {
	ID         uint           `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Date       string         `json:"date"`
	Time       string         `json:"time"`
	Frequency  string         `json:"frequency"`
	Duration   string         `json:"duration"`
	MedicineID uint           `json:"medicine_id"`
	Medicine   uint           `json:"medicine" gorm:"foreignKey:MedicineID"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index,column:deleted_at"`
}

func (r *Reminder) TableName() string {
	return "reminders"
}
