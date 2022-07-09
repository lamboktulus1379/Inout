package model

import "time"

type UserActivity struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement:true;column:id;not null" json:"id"`
	UserID       uint64     `gorm:"column:user_id;not null" json:"user_id"`
	ActivityType uint8      `gorm:"column:activity_type;type:int;not null" json:"activity_type"`
	Latitude     float64    `gorm:"column:latitude; not null" json:"latitude"`
	Longitude    float64    `gorm:"column:longitude;not null" json:"longitude"`
	CreatedAt    *time.Time `gorm:"autoCreateTime;index:idx_created_at;column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP"`
	CreatedBy    uint64     `gorm:"column:created_by;type:varchar(225);not null"`
	UpdatedAt    *time.Time `gorm:"autoUpdateTime;column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedBy    uint64     `gorm:"column:updated_by;type:varchar(225);not null"`
}

func (UserActivity) TableName() string {
	return "UserActivities"
}

type ReqUserActivity struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	UserID    uint64  `json:"user_id"`
}

type ReqCreateUserActivity struct {
	Latitude     float64 `json:"latitude" binding:"required"`
	Longitude    float64 `json:"longitude" binding:"required"`
	UserID       uint64  `json:"user_id"`
	ActivityType uint8   `json:"activity_type"`
}

type ReqUpdateUserActivity struct {
	Latitude     float64 `json:"latitude" binding:"required"`
	Longitude    float64 `json:"longitude" binding:"required"`
	UserID       uint64  `json:"user_id"`
	ActivityType uint8   `json:"activity_type"`
}
type ReqURIParamUserActivity struct {
	ID     uint64 `json:"id" uri:"id"`
	UserID uint64 `json:"user_id"`
}

type ReqQueryParamUserActivity struct {
	StartDate    string `json:"start_date" form:"start_date"`
	EndDate      string `json:"end_date" form:"end_date"`
	PerPage      int    `json:"per_page" form:"per_page" default:"10"`
	PageNumber   int    `json:"page_number" form:"page_number" default:"1"`
	ActivityType string `json:"activity_type" form:"activity_type"`
	UserID       uint64 `json:"user_id"`
}
