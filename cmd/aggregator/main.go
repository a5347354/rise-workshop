package main

import (
	"github.com/a5347354/rise-workshop/internal/aggregator"
	"github.com/a5347354/rise-workshop/internal/aggregator/usecase"
	"github.com/a5347354/rise-workshop/internal/event/store/postgres"
	"github.com/a5347354/rise-workshop/pkg"

	"context"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func main() {
	pkg.NewTracerProvider()

	fx.New(
		fx.Provide(
			pkg.NewPostgresClient,
			postgres.NewEventStore,
			usecase.NewAggregator,
		),
		fx.Invoke(
			func(usecase aggregator.Usecase) error {
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				usecase.Collect(ctx)
				time.Sleep(time.Second * 360)
				return nil
			},
		),
	).Run()
}
