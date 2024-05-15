package response

import (
	"gorm.io/gorm"
	"time"
)

type UserResponse struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name"`
	Age       int            `json:"age"`
	Weight    int            `json:"weight"`
	Height    int            `json:"height"`
	NIK       int            `json:"nik"`
	Birthday  string         `json:"birthday"`
	Gender    string         `json:"gender"`
	Address   string         `json:"address"`
	Phone     string         `json:"phone"`
	Username  string         `json:"username"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
