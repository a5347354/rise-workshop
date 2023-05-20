package aggregator

import (
	"context"
)

type Usecase interface {
	Collect(ctx context.Context)
}
