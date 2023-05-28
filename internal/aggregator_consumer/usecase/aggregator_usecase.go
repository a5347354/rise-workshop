package usecase

import (
	"github.com/a5347354/rise-workshop/internal"
	"github.com/a5347354/rise-workshop/internal/aggregator_consumer"
	"github.com/a5347354/rise-workshop/internal/event"
	"github.com/a5347354/rise-workshop/pkg"

	"context"
	"encoding/json"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/sirupsen/logrus"
)

type consumerUsecase struct {
	eStore       event.Store
	pubsubClient pkg.PubsubClient
}

func NewConsumer(eStore event.Store, pubsubClient pkg.PubsubClient) consumer.Usecase {
	return &consumerUsecase{eStore, pubsubClient}
}

func (u consumerUsecase) Consume(ctx context.Context) error {
	var mu sync.Mutex
	return u.pubsubClient.Consume(ctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		mu.Lock()
		defer mu.Unlock()
		var event internal.Event
		fmt.Printf("Got message: %q\n", string(msg.Data))
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			logrus.Warn(fmt.Errorf("unmarshal fail: %w", err))
		}
		if err := u.eStore.Insert(ctx, event); err != nil {
			logrus.Warn(fmt.Errorf("insert fail: %w", err))
		}
	})
}
