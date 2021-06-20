package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Base contains common columns for all tables.
type Base struct {
	ID        string `gorm:"type:string;primary_key;" json:"id"`
	CreatedAt uint   `json:"created_at"`
	UpdatedAt uint   `json:"updated_at"`
	DeletedAt *uint  `json:"deleted_at"`
}

func (base *Base) BeforeCreate(*gorm.DB) (err error) {
	base.ID = uuid.NewV4().String()
	return nil
}
