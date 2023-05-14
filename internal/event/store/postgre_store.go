package store

import (
	"github.com/a5347354/rise-workshop/internal"
	"github.com/a5347354/rise-workshop/internal/event"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"context"
)

type eventStore struct {
	db *gorm.DB
}

func NewEventStore(db *gorm.DB) event.Store {
	if viper.GetBool("db.automigrate") {
		db.AutoMigrate(internal.Event{})
	}
	return &eventStore{db}
}

func (e eventStore) Insert(ctx context.Context, event internal.Event) error {
	if err := e.db.WithContext(ctx).
		Create(&event).Error; err != nil {
		return err
	}
	return nil
}
