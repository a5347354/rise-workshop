package usecase

import (
	"github.com/a5347354/rise-workshop/internal/event"
	"github.com/a5347354/rise-workshop/internal/relay"
	"github.com/a5347354/rise-workshop/pkg"

	"context"
	"encoding/json"
	"fmt"

	"github.com/nbd-wtf/go-nostr"
	"github.com/sirupsen/logrus"
)

type relayUsecase struct {
	eStore event.Store
}

func NewRelay() relay.Usecase {
	return &relayUsecase{}
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
	// TODO: Verify the event ID and signature
	var evt nostr.Event
	if err := json.Unmarshal(eventJSON, &evt); err != nil {
		return pkg.WebSocketMsg{}, nil
	}
	//serialized := evt.Serialize()
	//hash := sha256.Sum256(serialized)
	//evt.ID = hex.EncodeToString(hash[:])
	//if ok, err := evt.CheckSignature(); err != nil {
	//	return fmt.Errorf("error: failed to verify signature")
	//} else if !ok {
	//	logrus.Warn(fmt.Sprintf("eventID:%s", evt.ID))
	//	logrus.Warn(fmt.Sprintf("pubky:%s", evt.PubKey))
	//	logrus.Warn(fmt.Sprintf("sig:%s", evt.Sig))
	//
	//	return fmt.Errorf("invalid: signature is invalid")
	//}

	// TODO: Store the event in the database

	// Log the event
	logrus.Infof("Received event: %s", string(eventJSON))

	return pkg.WebSocketMsg{}, nil
}

func (c relayUsecase) handleRequestMessage(ctx context.Context, requestData []json.RawMessage) (pkg.WebSocketMsg, error) {
	// TODO: Process the request and send appropriate events to the client(s)

	return pkg.WebSocketMsg{}, nil
}

func (c relayUsecase) handleCloseMessage(ctx context.Context, subscriptionID json.RawMessage) (pkg.WebSocketMsg, error) {
	// TODO: Close the subscription with the given ID

	return pkg.WebSocketMsg{}, nil
}
