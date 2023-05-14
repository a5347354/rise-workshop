package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/a5347354/rise-workshop/internal/event"
	"github.com/a5347354/rise-workshop/internal/relay"
	"github.com/a5347354/rise-workshop/pkg"

	"context"
	"encoding/json"
	"fmt"

	"github.com/nbd-wtf/go-nostr"
)

type relayUsecase struct {
	eStore event.Store
}

func NewRelay(eStore event.Store) relay.Usecase {
	return &relayUsecase{eStore}
}

func (c relayUsecase) ReceiveMessage(ctx context.Context, msg []byte) (pkg.WebSocketMsg, error) {
	var message []json.RawMessage
	if err := json.Unmarshal(msg, &message); err != nil {
		return pkg.WebSocketMsg{}, nil
	}

	if len(message) < 2 {
		return pkg.WebSocketMsg{}, nil
	}

	var messageType string

	if err := json.Unmarshal(message[0], &messageType); err != nil {
		return pkg.WebSocketMsg{}, nil
	}

	// Handle different message types
	switch messageType {
	case "EVENT":
		return c.handleEventMessage(ctx, message[1])
	case "REQ":
		return c.handleRequestMessage(ctx, message[1:])
	case "CLOSE":
		return c.handleCloseMessage(ctx, message[1])
	default:
		return pkg.WebSocketMsg{}, fmt.Errorf("unknown message type: %s", messageType)
	}
}

func (c relayUsecase) handleEventMessage(ctx context.Context, eventJSON json.RawMessage) (pkg.WebSocketMsg, error) {
	event, msg, err := c.verify(eventJSON)
	if err != nil {
		return msg, err
	}

	c.eStore.Insert(ctx, pkg.NostrEventToEvent(event))
	return pkg.WebSocketMsg{
		Action: pkg.WebSocketMsgTypeBroadcast,
		Msg:    eventJSON,
	}, nil
}

func (c relayUsecase) verify(eventJSON json.RawMessage) (nostr.Event, pkg.WebSocketMsg, error) {
	var evt nostr.Event
	if err := json.Unmarshal(eventJSON, &evt); err != nil {
		return nostr.Event{}, pkg.WebSocketMsg{}, nil
	}
	serialized := evt.Serialize()
	hash := sha256.Sum256(serialized)
	evt.ID = hex.EncodeToString(hash[:])
	if ok, err := evt.CheckSignature(); err != nil {
		return nostr.Event{}, pkg.WebSocketMsg{}, fmt.Errorf("error: failed to verify signature")
	} else if !ok {
		return nostr.Event{}, pkg.WebSocketMsg{}, fmt.Errorf("invalid: signature is invalid")
	}
	return evt, pkg.WebSocketMsg{}, nil
}

func (c relayUsecase) handleRequestMessage(ctx context.Context, requestData []json.RawMessage) (pkg.WebSocketMsg, error) {
	// TODO: Process the request and send appropriate events to the client(s)

	return pkg.WebSocketMsg{}, nil
}

func (c relayUsecase) handleCloseMessage(ctx context.Context, subscriptionID json.RawMessage) (pkg.WebSocketMsg, error) {
	// TODO: Close the subscription with the given ID

	return pkg.WebSocketMsg{}, nil
}
