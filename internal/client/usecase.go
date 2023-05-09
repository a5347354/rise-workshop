package client

import "context"

type Usecase interface {
	SendMessage(ctx context.Context) error
}
