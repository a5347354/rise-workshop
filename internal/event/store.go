package event

import "context"

type Store interface {
	List(ctx context.Context) []string
	Upsert(ctx context.Context, id string) error
	Remove(ctx context.Context, id string) error
}
