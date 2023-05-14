package internal

import (
	"gorm.io/gorm"
	"time"
)

type Event struct {
	ID        string         `json:"id"`
	Kind      int            `gorm:"column:kind;not null;default:0`
	Content   string         `gorm:"column:content;type:varchar(500);not null`
	CreatedAt time.Time      `gorm:"column:created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName sets the table name for relational databases
func (Event) TableName() string {
	return "event"
}
