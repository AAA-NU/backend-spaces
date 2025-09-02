package models

import "github.com/lib/pq"

type CoworkingModel struct {
	ID         string         `json:"id"`
	BookedTime pq.StringArray `gorm:"type:text[]"`
}

type RoomModel struct {
	ID       string `json:"id"`
	IsBooked bool   `json:"is_booked"`
	BookedBy string `json:"booked_by"`
}
