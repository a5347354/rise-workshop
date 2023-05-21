package client

import (
	"context"
)

type Usecase interface {
	SendMessage(ctx context.Context) error
	Collect(ctx context.Context, url string) error
}
