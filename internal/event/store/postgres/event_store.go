package postgres

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

func (e *eventStore) Insert(ctx context.Context, event internal.Event) error {
	if err := e.db.WithContext(ctx).
		Create(&event).Error; err != nil {
		return err
	}
	return nil
}

func (e *eventStore) SearchByContent(ctx context.Context, keyword string) ([]internal.Event, error) {
	var events []internal.Event

	if err := e.db.WithContext(ctx).
		Model(&internal.Event{}).
		Select("id, kind, content").
		Where("content ILIKE ?", "%"+keyword+"%").
		Find(&events).Error; err != nil {
		return []internal.Event{}, err
	}
	return events, nil
}
