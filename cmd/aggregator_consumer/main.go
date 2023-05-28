package main

import (
	"context"
	consumer "github.com/a5347354/rise-workshop/internal/aggregator_consumer"
	"github.com/a5347354/rise-workshop/internal/aggregator_consumer/usecase"
	"github.com/a5347354/rise-workshop/internal/event/store/postgres"
	"github.com/a5347354/rise-workshop/pkg"
	"github.com/spf13/viper"

	"go.uber.org/fx"
	//_ "net/http/pprof"
	"strings"
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
			usecase.NewConsumer,
		),
		fx.Invoke(
			func(u consumer.Usecase) error {
				return u.Consume(context.Background())
			},
		),
	).Run()
}
