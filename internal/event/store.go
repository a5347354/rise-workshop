package event

import (
	"github.com/a5347354/rise-workshop/internal"

	"context"
)

type Store interface {
	Insert(ctx context.Context, event internal.Event) (err error)
	SearchByContent(ctx context.Context, keyword string) ([]internal.Event, error)
}
