package relay

import (
	"github.com/a5347354/rise-workshop/pkg"
	"gopkg.in/olahol/melody.v1"

	"context"
)

type Usecase interface {
	ReceiveMessage(ctx context.Context, msg []byte, session *melody.Session) (pkg.WebSocketMsg, error)
	SendMessageToSubscriber(ctx context.Context) error
}
