package domain

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Username     string    `gorm:"type:varchar;uniqueIndex;not null"`
	Email        string    `gorm:"type:varchar;uniqueIndex;not null"`
	PasswordHash string    `gorm:"type:varchar;not null"`
	CreatedAt    time.Time `gorm:"not null;default:now()"`
}

func (User) TableName() string {
	return "users"
}
