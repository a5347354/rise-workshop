package pkg

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func getTopicName() string {
	arr := strings.Split(viper.GetString("service.id"), "_")
	return arr[0]
}

func getSubscriberName() string {
	randStr, _ := GenerateAlphanumericString(6)
	return fmt.Sprintf("%s_subscriber_%s", viper.GetString("service.id"), randStr)
}

func getProjectName() string {
	return viper.GetString("gcp.project.id")
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
	client         *pubsub.Client
	topic          *pubsub.Topic
	subscriptionID string
}

func NewPub(lc fx.Lifecycle) (PubsubClient, error) {
	ctx := context.Background()
	client, err := NewPubSubWithoutLC(ctx)
	if err != nil {
		return &pubsubClient{}, err
	}
	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				return client.close(context.Background())
			},
		},
	)
	if err := client.createTopic(ctx, getTopicName()); err != nil {
		return &pubsubClient{}, err
	}
	return &client, nil
}

func NewSub(lc fx.Lifecycle) (PubsubClient, error) {
	ctx := context.Background()
	client, err := NewPubSubWithoutLC(ctx)
	if err != nil {
		return &pubsubClient{}, err
	}
	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				logrus.Info(1122333)
				return client.close(ctx)
			},
		},
	)
	if err := client.createTopic(ctx, getTopicName()); err != nil {
		logrus.Warn(err)
		return &pubsubClient{}, err
	}

	if err := client.createSubscription(ctx); err != nil {
		return nil, err
	}
	return &client, nil
}

func (c *pubsubClient) createTopic(ctx context.Context, topic string) error {
	t := c.client.Topic(topic)
	ok, err := t.Exists(ctx)
	if err != nil {
		return err
	}
	if ok {
		c.setTopic(t)
		return nil
	}
	t, err = c.client.CreateTopic(ctx, topic)
	if err != nil {
		return err
	}
	c.setTopic(t)
	return nil
}

func (c *pubsubClient) setTopic(t *pubsub.Topic) {
	c.topic = t
}

func (c *pubsubClient) createSubscription(ctx context.Context) error {
	c.setSubscriptionID()
	sub, err := c.client.CreateSubscription(ctx, c.subscriptionID, pubsub.SubscriptionConfig{
		Topic:       c.topic,
		AckDeadline: 20 * time.Second,
	})
	if err != nil {
		logrus.Warn(err)
		return err
	}
	fmt.Printf("Created subscription: %v\n", sub)
	return nil
}

func (c *pubsubClient) close(ctx context.Context) error {
	if c.subscriptionID != "" {
		c.client.DetachSubscription(ctx, c.subscriptionID)
	}
	return c.client.Close()
}

func (c *pubsubClient) Publish(data []byte) error {
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

func (c *pubsubClient) Consume(ctx context.Context, receive func(ctx context.Context, msg *pubsub.Message)) error {
	sub := c.client.Subscription(c.subscriptionID)
	err := sub.Receive(ctx, receive)
	if err != nil {
		return err
	}
	return nil
}

func (c *pubsubClient) setSubscriptionID() {
	c.subscriptionID = getSubscriberName()
}
