package entity

import (
	"gorm.io/gorm"
	"time"
)

type Appointment struct {
	ID          uint           `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Date        string         `json:"date" form:"date"`
	Time        string         `json:"time" form:"time"`
	Complaint   string         `json:"complaint" form:"complaint"`
	ApprovedID  *uint          `json:"approved_id" gorm:"default:null"` // Changed to *uint and allow null
	RequestedID uint           `json:"requested_id" form:"requested_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index,column:deleted_at"`
}

type AppointmentByAuth struct {
	ID         uint   `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Date       string `json:"date" form:"date"`
	Time       string `json:"time" form:"time"`
	Complaint  string `json:"complaint" form:"complaint"`
	ApprovedID *uint  `json:"approved_id" gorm:"default:null"` // Changed to *uint and allow null
}

type AppointmentResponse struct {
	ID          uint   `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Date        string `json:"date" form:"date"`
	Time        string `json:"time" form:"time"`
	Complaint   string `json:"complaint" form:"complaint"`
	ApprovedID  *uint  `json:"approved_id" form:"approved_id"`
	RequestedID uint   `json:"requested_id" form:"requested_id"`
}

type AppointmentUpdate struct {
	Date        string `json:"date" form:"date"`
	Time        string `json:"time" form:"time"`
	Complaint   string `json:"complaint" form:"complaint"`
	ApprovedID  *uint  `json:"approved_id" form:"dateapproved_id"` // Changed to *uint
	RequestedID *uint  `json:"requested_id" form:"requested_id"`
}

func (AppointmentResponse) TableName() string {
	return "appointments"
}
