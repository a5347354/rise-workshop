package main

import (
	"github.com/a5347354/rise-workshop/internal/aggregator"
	"github.com/a5347354/rise-workshop/internal/aggregator/delivery"
	"github.com/a5347354/rise-workshop/internal/aggregator/usecase"
	"github.com/a5347354/rise-workshop/internal/event/store/postgres"
	"github.com/a5347354/rise-workshop/internal/event/store/pubsub"
	"github.com/a5347354/rise-workshop/pkg"

	"strings"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	//_ "net/http/pprof"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func main() {
	pkg.NewTracerProvider()

	fx.New(
		fx.Provide(
			pkg.NewRouter,
			pkg.NewPostgresClient,
			pkg.NewPub,
			postgres.NewEventStore,
			pubsub.NewEventStore,
			usecase.NewAggregator,
		),
		fx.Invoke(
			delivery.RegisterAggregatorHandler,
			func(u aggregator.Usecase) error {
				// debug goruntine
				//go func() {
				//	log.Println(http.ListenAndServe("localhost:6060", nil))
				//}()
				u.StartCollect()
				return nil
			},
		),
	).Run()
}
