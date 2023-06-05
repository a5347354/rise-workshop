package relay

import (
	"context"

	"gopkg.in/olahol/melody.v1"
)

type Notification interface {
	Subscribe(ctx context.Context, id string, s *melody.Session)
	UnSubscribe(ctx context.Context, id string)
	Broadcast(ctx context.Context, msg []interface{})
}
