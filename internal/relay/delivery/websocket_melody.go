package delivery

import (
	"context"
	"github.com/a5347354/rise-workshop/internal/relay"
	"github.com/sirupsen/logrus"

	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gopkg.in/olahol/melody.v1"
)

type Metrics struct {
	successTotal    *prometheus.CounterVec
	failTotal       *prometheus.CounterVec
	processDuration *prometheus.HistogramVec
}
type relayHandler struct {
	usecase relay.Usecase
	metrics Metrics
}

func RegistWebsocketHandler(engine *gin.Engine, m *melody.Melody, usecase relay.Usecase) {
	r := &relayHandler{
		usecase: usecase,
		metrics: Metrics{
			successTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "success_total",
					Help: "How many processed success, partitioned by database, table and type.",
				},
				[]string{},
			),
			failTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "fail_total",
					Help: "How many processed fail, partitioned by database, table, type and stage.",
				},
				[]string{},
			),
			processDuration: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "processed_duration_second",
					Help:    "The bprocessed latencies in second, partitioned by database, table and type.",
					Buckets: prometheus.DefBuckets,
				},
				[]string{},
			),
		},
	}
	engine.GET("/wss", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleConnect(r.handleConnect())
	m.HandleDisconnect(r.handleDisconnect())
	m.HandleError(r.handleError())
	m.HandleMessage(r.message())
}

func (r relayHandler) handleConnect() func(s *melody.Session) {
	return func(s *melody.Session) {
		fmt.Printf("[Melody] | %s | Connect %s\n", s.Request.RemoteAddr, s.Request.RequestURI)
	}
}

func (r relayHandler) handleDisconnect() func(s *melody.Session) {
	return func(s *melody.Session) {
		fmt.Printf("[Melody] | %s | Disconnect %s\n", s.Request.RemoteAddr, s.Request.RequestURI)
	}
}

func (r relayHandler) handleError() func(s *melody.Session, err error) {
	return func(s *melody.Session, err error) {
		fmt.Printf("[Melody] | %s | Error %s %s\n", s.Request.RemoteAddr, s.Request.RequestURI, err)
	}
}

func (r relayHandler) message() func(s *melody.Session, msg []byte) {
	return func(s *melody.Session, msg []byte) {
		t := time.Now()
		err := r.usecase.ReceiveMessage(context.Background(), msg)
		if err != nil {
			logrus.Panic(err)
		}
		fmt.Printf("[Melody] %v | %13v | %s | Message %s\n",
			t.Format("2006/01/02 - 15:04:05"),
			time.Since(t),
			s.Request.RemoteAddr,
			s.Request.RequestURI,
		)

		r.metrics.successTotal.WithLabelValues().Inc()
		r.metrics.processDuration.WithLabelValues().Observe(float64(time.Since(t)) / float64(time.Second))
	}
}
