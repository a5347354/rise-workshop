package main

import (
	"github.com/a5347354/rise-workshop/internal/client/delivery"
	"github.com/a5347354/rise-workshop/internal/client/usecase"
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
			pkg.NewRouter,

			usecase.NewClient,
		),
		fx.Invoke(
			delivery.RegisterClientHandler,
		),
	).Run()
}
