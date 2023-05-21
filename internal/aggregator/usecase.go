package aggregator

import (
	"github.com/a5347354/rise-workshop/internal"

	"context"
)

type Usecase interface {
	Collect(ctx context.Context)
	StartCollect() error
	ListEventByKeyword(ctx context.Context, keyword string) ([]internal.Event, error)
}
