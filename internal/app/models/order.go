package models

import "time"

type Order struct {
	ID         string
	UserID     int
	StatusID   int
	UploadedAt time.Time
}
