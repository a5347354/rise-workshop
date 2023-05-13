package relay

import (
	"github.com/a5347354/rise-workshop/pkg"

	"context"
)

type Usecase interface {
	ReceiveMessage(ctx context.Context, msg []byte) (pkg.WebSocketMsg, error)
}
