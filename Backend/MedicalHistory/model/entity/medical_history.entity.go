package entity

import (
	"gorm.io/gorm"
	"time"
)

type MedicalHistory struct {
	ID             uint           `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	UserID         uint           `json:"UserID"`
	DoctorReportID uint           `json:"DoctorReportID"`
	DoctorReport   uint           `json:"doctorReport" gorm:"foreignKey:DoctorReportID"`
	User           uint           `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index,column:deleted_at"`
}

type MedicalHistoryResponse struct {
	UserID         uint `json:"UserID"`
	DoctorReportID uint `json:"DoctorReportID"`
}

func (m *MedicalHistory) TableName() string {
	return "medical_histories"
}
