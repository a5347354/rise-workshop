package usecase

import (
	"github.com/a5347354/rise-workshop/internal/relay"

	"context"
)

type relayUsecase struct {
}

func NewRelay() relay.Usecase {
	return &relayUsecase{}
}

func (c relayUsecase) ReceiveMessage(ctx context.Context, msg []byte) error {
	return nil
}
