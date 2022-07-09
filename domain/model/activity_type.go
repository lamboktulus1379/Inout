package model

type ActivityType struct {
	ID   uint8  `gorm:"primaryKey;column:id;type:int;not null" json:"id"`
	Name string `gorm:"column:name;type:varchar(45)"`
}

func (ActivityType) TableName() string {
	return "ActivityType"
}
