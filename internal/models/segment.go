package models

import "time"

type Segment struct {
	Name      string    `gorm:"primary_key" json:"name"`
	Users     []User    `gorm:"many2many:user_segments" json:"users,omitempty"`
	CreatedAt time.Time `gorm:"default:now()" json:"-"`
}

func (Segment) TableName() string {
	return "segment"
}
