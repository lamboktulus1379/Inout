package dto

import "time"

type ResUserActivity struct {
	Res
	*Pagination
	Data []UserActivity `json:"data"`
}

type UserActivity struct {
	ID           uint64     `json:"id"`
	UserID       uint64     `json:"user_id"`
	ActivityType string     `json:"activity_type"`
	Latitude     float64    `json:"latitude"`
	Longitude    float64    `json:"longitude"`
	CreatedAt    *time.Time `json:"created_at"`
	CreatedBy    uint64     `json:"created_by"`
	UpdatedAt    *time.Time `json:"updated_at"`
	UpdatedBy    uint64     `json:"updated_by"`
}
