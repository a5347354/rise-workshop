package consumer

import (
	"context"
)

type Usecase interface {
	Consume(ctx context.Context) error
}
