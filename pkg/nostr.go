package pkg

import (
	"context"
	"fmt"
	
	"github.com/nbd-wtf/go-nostr"
	"github.com/spf13/viper"
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
}

func NewNostrClient() NostrClient {
	if viper.GetBool("start") {

	}
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

func (c *nostrClient) Publish(ctx context.Context, e nostr.Event) (nostr.Status, error) {
	e.ID = e.GetID()
	e.Sign(c.privateKey)
	_, err := e.CheckSignature()
	if err != nil {
		return nostr.PublishStatusFailed, err
	}
	fmt.Println(e.ID)
	e.PubKey = c.publicKey
	e.CreatedAt = nostr.Now()
	status, err := c.relay.Publish(ctx, e)
	if err != nil {
		return nostr.PublishStatusFailed, err
	}
	if status != nostr.PublishStatusSucceeded {
		return status, fmt.Errorf("relay no response")
	}
	return nostr.PublishStatusFailed, err
}

func (c *nostrClient) Disconnect(ctx context.Context) error {
	if c.relay != nil {
		return c.relay.Close()
	}
	return nil
}
