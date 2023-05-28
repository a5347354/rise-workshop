package pubsub

import (
	"github.com/a5347354/rise-workshop/internal"
	"github.com/a5347354/rise-workshop/internal/event"
	"github.com/a5347354/rise-workshop/pkg"

	"context"
	"encoding/json"
)

type eventStore struct {
	pub pkg.PubsubClient
}

func NewEventStore(pub pkg.PubsubClient) event.AsyncStore {
	return &eventStore{pub}
}

func (e *eventStore) Insert(ctx context.Context, event internal.Event) error {
	b, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return e.pub.Publish(b)
}
