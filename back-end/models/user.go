package models

import "time"

type User struct {
	ID           uint
	Name         string
	Email        string
	PasswordHash string
	IconUrl      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	RowVersion   uint     `gorm:"default:1"`
	Circles      []Circle `gorm:"many2many:users_circles"`
}
