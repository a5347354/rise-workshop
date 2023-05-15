package pkg

import (
	"context"
	"fmt"
	"github.com/a5347354/rise-workshop/internal"
	"github.com/nbd-wtf/go-nostr"
	"github.com/spf13/viper"
	"strings"
)

type nostrClient struct {
	relay      *nostr.Relay
	relayURL   string
	privateKey string
	publicKey  string
}

type NostrClient interface {
	Connect(ctx context.Context) error
	Publish(ctx context.Context, e nostr.Event) (nostr.Status, error)
	Disconnect(ctx context.Context) error
	Subscribe(ctx context.Context, filters nostr.Filters) (*nostr.Subscription, error)
}

func NewNostrClient() NostrClient {
	return &nostrClient{
		nil,
		viper.GetString("relay.url"),
		viper.GetString("private.key"),
		viper.GetString("public.key"),
	}
}

func (c *nostrClient) Connect(ctx context.Context) error {
	r, err := nostr.RelayConnect(ctx, c.relayURL+"/ws")
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

func NostrEventToEvent(e nostr.Event) internal.Event {
	return internal.Event{
		ID:        e.ID,
		Kind:      e.Kind,
		Content:   e.Content,
		CreatedAt: e.CreatedAt.Time(),
	}
}
