package main

import (
	"context"
	consumer "github.com/a5347354/rise-workshop/internal/aggregator_consumer"
	"github.com/a5347354/rise-workshop/internal/aggregator_consumer/usecase"
	"github.com/a5347354/rise-workshop/internal/event/store/postgres"
	"github.com/a5347354/rise-workshop/pkg"
	"github.com/gin-gonic/gin"
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
			pkg.NewRouter,
			pkg.NewPostgresClient,
			pkg.NewSub,
			postgres.NewEventStore,
			usecase.NewConsumer,
		),
		fx.Invoke(
			func(engine *gin.Engine, u consumer.Usecase) error {
				go u.Consume(context.Background())
				return nil
			},
		),
	).Run()

}
