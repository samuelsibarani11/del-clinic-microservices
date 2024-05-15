package entity

import (
	"gorm.io/gorm"
	"time"
)

type DoctorReport struct {
	ID            uint           `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	NurseReportID uint           `json:"nurse_report_id"`
	StaffDoctorID uint           `json:"staff_doctor_id"`
	Disease       string         `json:"disease"`
	Amount        int            `json:"amount"`
	MedicineID    int            `json:"medicine_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index,column:deleted_at"`
}

type DoctorReportResponse struct {
	ID            uint   `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Disease       string `json:"disease"`
	NurseReportID uint   `json:"nurse_report_id"`
	MedicineID    uint   `json:"medicine_id"`
	StaffDoctorID uint   `json:"staff_doctor_id"`
}

func (d *DoctorReport) TableName() string {
	return "doctor_reports"
}
