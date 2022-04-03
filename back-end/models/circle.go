package models

import "time"

type Circle struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
