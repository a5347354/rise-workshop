package internal

import (
	"gorm.io/gorm"
	"time"
)

type Event struct {
	PK        int64          `gorm:"primaryKey;auto_increment;not_null" json:"-"`
	ID        string         `gorm:"id" json:"id"`
	Kind      int            `gorm:"column:kind;not null;default:0" json:"kind"`
	Content   string         `gorm:"column:content;type:varchar(500);not null" json:"content"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName sets the table name for relational databases
func (Event) TableName() string {
	return "event"
}
