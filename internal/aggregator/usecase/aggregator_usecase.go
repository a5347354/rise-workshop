package usecase

import (
	"github.com/a5347354/rise-workshop/internal"
	"github.com/a5347354/rise-workshop/internal/aggregator"
	"github.com/a5347354/rise-workshop/internal/client"
	"github.com/a5347354/rise-workshop/internal/client/usecase"
	"github.com/a5347354/rise-workshop/internal/event"

	"context"
	"strings"
	"sync"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type aggregatorUsecase struct {
	eStore      event.Store
	url         []string
	lc          fx.Lifecycle
	limitClient chan client.Usecase
}

func NewAggregator(lc fx.Lifecycle, eStore event.Store) aggregator.Usecase {
	url := strings.Split(viper.GetString("relays.url"), ",")
	return &aggregatorUsecase{eStore, url, lc, make(chan client.Usecase, len(url))}
}

func (u aggregatorUsecase) Collect(ctx context.Context) {
	for _, url := range u.url {
		c := usecase.NewClient(u.lc, u.eStore)
		u.limitClient <- c
		go func(limitClient chan client.Usecase) {
			err := c.Collect(ctx, url)
			if err != nil {
				<-limitClient
			}
		}(u.limitClient)
	}
}

func (u aggregatorUsecase) StartCollect() error {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	u.lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				defer wg.Done()
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				for {
					u.Collect(ctx)
				}

			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			wg.Wait()
			u.Close()
			return nil
		},
	})
	return nil
}

func (u aggregatorUsecase) Close() {
	close(u.limitClient)
}

func (u aggregatorUsecase) ListEventByKeyword(ctx context.Context, keyword string) ([]internal.Event, error) {
	return u.eStore.SearchByContent(ctx, keyword)
}
