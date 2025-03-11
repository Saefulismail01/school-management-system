package models

import (
	"time"
)

// Pengajar model
type Pengajar struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Nama      string    `gorm:"not null" json:"nama"`
	Email     string    `gorm:"unique" json:"email"`
	NoTelepon string    `json:"no_telepon"`
	Alamat    string    `json:"alamat"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	User User `gorm:"foreignKey:UserID" json:"-"`
}