package entity

import "time"

type Staff struct {
	ID             uint      `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Name           string    `json:"name"`
	Age            int       `json:"age"`
	Weight         int       `json:"weight"`
	Height         int       `json:"height"`
	NIP            int       `json:"NIP"`
	Birthday       string    `json:"birthday"`
	Gender         string    `json:"gender"`
	Address        string    `json:"address"`
	Phone          string    `json:"phone"`
	Username       string    `json:"username" gorm:"unique"`
	Password       string    `json:"password" gorm:"column:password"`
	Role           int       `json:"role"`
	ProfilePicture *string   `json:"profilePicture"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime" db:"updated_at"`
}

func (Staff) TableName() string {
	return "staffs"
}
