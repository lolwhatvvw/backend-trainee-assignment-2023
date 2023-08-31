package models

import (
	"time"
)

type User struct {
	ID        int64     `gorm:"primary_key"`
	FirstName string    `gorm:"column:firstname" json:"firstname" validate:"required"`
	LastName  string    `gorm:"column:lastname" json:"lastname" validate:"required"`
	Username  string    `gorm:"size:128;uniqueIndex" json:"username" validate:"required"`
	Segments  []Segment `gorm:"many2many:user_segments" json:"segments,omitempty" validate:"required"`
	CreatedAt time.Time `gorm:"default:now()" json:"-"`
}
