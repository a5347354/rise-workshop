package usecase

import (
	"context"
	"github.com/a5347354/rise-workshop/internal/aggregator"
	"github.com/a5347354/rise-workshop/internal/client/usecase"
	"github.com/a5347354/rise-workshop/internal/event"

	"context"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type aggregatorUsecase struct {
	eStore event.Store
	url    []string
	lc     fx.Lifecycle
}

func NewAggregator(lc fx.Lifecycle, eStore event.Store) aggregator.Usecase {
	return &aggregatorUsecase{eStore, strings.Split(viper.GetString("relays.url"), ","), lc}
}

func (u aggregatorUsecase) Collect(ctx context.Context) {
	for _, url := range u.url {
		c := usecase.NewClient(u.lc, u.eStore)
		go c.Collect(ctx, url)
	}
}
