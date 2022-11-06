package models

import "time"

type Circle struct {
	ID         uint
	Name       string
	IconUrl    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	RowVersion uint `gorm:"default:1"`
}
