package relay

import "context"

type Usecase interface {
	ReceiveMessage(ctx context.Context, msg []byte) error
}
