package pkg

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func getTopicName() string {
	return viper.GetString("service.id")
}

func getSubscriberName() string {
	randStr, _ := GenerateAlphanumericString(6)
	return fmt.Sprintf("%s_subscriber_%s", viper.GetString("service.id"), randStr)
}

func getProjectName() string {
	return viper.GetString("service.id")
}

func NewPubSubWithoutLC(ctx context.Context) (pubsubClient, error) {
	proj := getProjectName()
	client, err := pubsub.NewClient(ctx, proj)
	if err != nil {
		return pubsubClient{}, err
	}
	return pubsubClient{
		client: client,
	}, nil
}

type PubsubClient interface {
	Publish(data []byte) error
	Consume(ctx context.Context, receive func(ctx context.Context, msg *pubsub.Message)) error
}

type pubsubClient struct {
	client *pubsub.Client
	topic  *pubsub.Topic
}

func NewPub(lc fx.Lifecycle) (pubsubClient, error) {
	ctx := context.Background()
	client, err := NewPubSubWithoutLC(ctx)
	if err != nil {
		return pubsubClient{}, err
	}
	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				return client.close()
			},
		},
	)
	if err := client.createTopic(ctx, getTopicName()); err != nil {
		return pubsubClient{}, err
	}
	return client, nil
}

func NewSub(lc fx.Lifecycle) (pubsubClient, error) {
	ctx := context.Background()
	client, err := NewPubSubWithoutLC(ctx)
	if err != nil {
		return pubsubClient{}, err
	}
	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				return client.close()
			},
		},
	)
	if err := client.createTopic(ctx, getTopicName()); err != nil {
		return pubsubClient{}, err
	}
	return client, client.createSubscription(ctx, getSubscriberName())
}

func (c pubsubClient) createTopic(ctx context.Context, topic string) error {
	t := c.client.Topic(topic)
	ok, err := t.Exists(ctx)
	if err != nil {
		return err
	}
	if ok {
		c.topic = t
		return nil
	}
	t, err = c.client.CreateTopic(ctx, topic)
	if err != nil {
		return err
	}
	c.topic = t
	return nil
}

func (c pubsubClient) createSubscription(ctx context.Context, name string) error {
	sub, err := c.client.CreateSubscription(ctx, name, pubsub.SubscriptionConfig{
		Topic:       c.topic,
		AckDeadline: 20 * time.Second,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Created subscription: %v\n", sub)
	return nil
}

func (c pubsubClient) close() error {
	return c.client.Close()
}

func (c pubsubClient) Publish(data []byte) error {
	ctx := context.Background()
	result := c.topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})
	id, err := result.Get(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Published a message; msg ID: %v\n", id)
	return nil
}

func (c pubsubClient) Consume(ctx context.Context, receive func(ctx context.Context, msg *pubsub.Message)) error {
	sub := c.client.Subscription(getSubscriberName())
	err := sub.Receive(ctx, receive)
	if err != nil {
		return err
	}
	return nil
}
