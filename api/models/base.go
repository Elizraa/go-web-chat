// models/base.go
package models

import (
	"time"
)

type BaseModel struct {
	ID        uint      `gorm:"-" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt *time.Time `sql:"index"`
}

// Model function for BaseModel
func (b *BaseModel) Model() *BaseModel {
	return b
}
