package main

import (
	"github.com/a5347354/rise-workshop/internal/relay/delivery"
	"github.com/a5347354/rise-workshop/internal/relay/usecase"
	"github.com/a5347354/rise-workshop/pkg"

	"strings"

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
			pkg.NewWebsocket,
			pkg.NewRouter,

			usecase.NewRelay,
		),
		fx.Invoke(
			delivery.RegistWebsocketHandler,
		),
	).Run()
}
