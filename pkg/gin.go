package pkg

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func init() {
	viper.SetDefault("PORT", 8080)
}

var traceIDLogFormatter = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}
	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	//logrus.Info(param.Request.UserAgent())
	return fmt.Sprintf("[GIN] %s %3d %s | %13v | %15s | %s %-7s %s %#v\n%s",
		// param.TimeStamp.Format("2006/01/02 15:04:05"),
		// trace.SpanContextFromContext(param.Request.Context()).TraceID(),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

func NewRouter(lc fx.Lifecycle) *gin.Engine {
	if viper.GetBool("debug") {
		gin.ForceConsoleColor()
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	engine.Use(
		gin.LoggerWithFormatter(traceIDLogFormatter),
		gin.Recovery(),
	)

	serverAddr := viper.GetInt("PORT")
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", serverAddr),
		Handler: engine,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logrus.WithError(err).Fatalf("listen and serve error")
				}
			}()
			log.Printf("[GIN] server start and listen (%d)", serverAddr)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Print("[GIN] server graceful stop")
			return srv.Shutdown(ctx)
		},
	})
	return engine
}

func AbortWithErrorJSON(c *gin.Context, code int, err error) {
	c.Error(err)
	c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
}
