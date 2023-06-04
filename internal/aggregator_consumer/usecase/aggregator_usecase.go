package usecase

import (
	"github.com/a5347354/rise-workshop/internal"
	"github.com/a5347354/rise-workshop/internal/aggregator_consumer"
	"github.com/a5347354/rise-workshop/internal/aggregator_consumer/delivery"
	"github.com/a5347354/rise-workshop/internal/event"
	"github.com/a5347354/rise-workshop/pkg"
	"time"

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
	metrics      delivery.Metrics
}

func NewConsumer(eStore event.Store, pubsubClient pkg.PubsubClient, metrics delivery.Metrics) consumer.Usecase {
	return &consumerUsecase{eStore, pubsubClient, metrics}
}

func (u consumerUsecase) Consume(ctx context.Context) error {
	var mu sync.Mutex
	return u.pubsubClient.Consume(ctx, func(ctx context.Context, msg *pubsub.Message) {
		t := time.Now()
		msg.Ack()
		mu.Lock()
		defer mu.Unlock()
		var event internal.Event
		fmt.Printf("Got message: %q\n", string(msg.Data))
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			logrus.Warn(fmt.Errorf("unmarshal fail: %w", err))
			u.metrics.FailTotal("consume")
		}
		if err := u.eStore.Insert(ctx, event); err != nil {
			u.metrics.FailTotal("consume")
			logrus.Warn(fmt.Errorf("insert fail: %w", err))
		}
		u.metrics.SuccessTotal()
		u.metrics.ProcessDuration(t)
	})
}
