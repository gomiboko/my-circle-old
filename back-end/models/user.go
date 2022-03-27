package models

import "time"

type User struct {
	ID           uint
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Circles      []Circle `gorm:"many2many:users_circles;"`
}
