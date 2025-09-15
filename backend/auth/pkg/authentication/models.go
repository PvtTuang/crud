package authentication

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Email        string `gorm:"unique"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
