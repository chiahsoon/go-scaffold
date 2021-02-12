package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Base contains common columns for all tables.
type Base struct {
	ID        string `gorm:"type:string;primary_key;"`
	CreatedAt uint
	UpdatedAt uint
	DeletedAt *uint
}

func (base *Base) BeforeCreate(*gorm.DB) (err error) {
	base.ID = uuid.NewV4().String()
	return nil
}
