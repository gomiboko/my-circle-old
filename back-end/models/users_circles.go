package models

import "time"

type UsersCircles struct {
	UserID    uint
	CircleID  uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
