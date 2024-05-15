package entity

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID             uint           `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Name           string         `json:"name"`
	Age            int            `json:"age"`
	Weight         int            `json:"weight"`
	Height         int            `json:"height"`
	NIK            int            `json:"nik"`
	Birthday       string         `json:"birthday"`
	Gender         string         `json:"gender"`
	Address        string         `json:"address"`
	Phone          string         `json:"phone"`
	Username       string         `json:"username" gorm:"unique"`
	Password       string         `json:"password" gorm:"column:password"`
	Role           int            `json:"role"`
	DormID         uint           `json:"dormID"`
	ProfilePicture *string        `json:"profilePicture"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index,column:deleted_at"`
}

func (User) TableName() string {
	return "users"
}
