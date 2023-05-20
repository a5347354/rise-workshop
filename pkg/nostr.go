package pkg

import (
	"github.com/a5347354/rise-workshop/internal"

	"context"
	"fmt"
	"strings"

	"github.com/nbd-wtf/go-nostr"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type nostrClient struct {
	relay      *nostr.Relay
	relayURL   string
	privateKey string
	publicKey  string
}

type NostrClient interface {
	Connect(ctx context.Context) error
	ConnectURL(ctx context.Context, url string) error
	Publish(ctx context.Context, e nostr.Event) (nostr.Status, error)
	Subscribe(ctx context.Context, filters nostr.Filters) (*nostr.Subscription, error)
	Disconnect(ctx context.Context) error
	GetClient() *nostr.Relay
}

func NewNostrClientWithoutLC() NostrClient {
	return &nostrClient{
		nil,
		viper.GetString("relay.url"),
		viper.GetString("private.key"),
		viper.GetString("public.key"),
	}
}

func NewNostrClient(lc fx.Lifecycle) NostrClient {
	c := NewNostrClientWithoutLC()
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			if c.GetClient() != nil {
				return c.GetClient().Close()
			}
			return nil
		},
	})
	return &nostrClient{
		nil,
		viper.GetString("relay.url"),
		viper.GetString("private.key"),
		viper.GetString("public.key"),
	}
}

func (c *nostrClient) Connect(ctx context.Context) error {
	r, err := nostr.RelayConnect(ctx, c.relayURL)
	if err != nil {
		return err
	}
	c.relay = r
	return nil
}

func (c *nostrClient) ConnectURL(ctx context.Context, url string) error {
	r, err := nostr.RelayConnect(ctx, url)
	if err != nil {
		return err
	}
	c.relay = r
	return nil
}

func (c *nostrClient) Publish(ctx context.Context, e nostr.Event) (nostr.Status, error) {
	e.ID = e.GetID()
	fmt.Println(e.ID)
	e.PubKey = c.publicKey
	e.Kind = 1
	e.CreatedAt = nostr.Now()
	e.Content = strings.TrimSpace(e.Content)
	e.Sign(c.privateKey)
	status, err := c.relay.Publish(ctx, e)
	if status != nostr.PublishStatusSucceeded {
		return status, fmt.Errorf("relay no response")
	}
	return nostr.PublishStatusFailed, err
}

func (c *nostrClient) Subscribe(ctx context.Context, filters nostr.Filters) (*nostr.Subscription, error) {
	return c.relay.Subscribe(ctx, filters)
}

func (c *nostrClient) Disconnect(ctx context.Context) error {
	if c.relay != nil {
		return c.relay.Close()
	}
	return nil
}

func (c *nostrClient) GetClient() *nostr.Relay {
	return c.relay
}

func NostrEventToEvent(e nostr.Event) internal.Event {
	return internal.Event{
		ID:        e.ID,
		Kind:      e.Kind,
		Content:   e.Content,
		CreatedAt: e.CreatedAt.Time(),
	}
}
