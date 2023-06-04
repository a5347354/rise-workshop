package delivery

import (
	"github.com/a5347354/rise-workshop/internal/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

type metrics struct {
	successTotal              *prometheus.CounterVec
	failTotal                 *prometheus.CounterVec
	processDuration           *prometheus.HistogramVec
	websocketConnectionNumber *prometheus.CounterVec
}

func NewClientMetrics() client.Metrics {
	return &metrics{
		successTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "success_total",
				Help: "How many processed success",
			},
			[]string{},
		),
		failTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "fail_total",
				Help: "How many processed fail, partitioned by stage.",
			},
			[]string{"stage"},
		),
		processDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "processed_duration_second",
				Help:    "The processed latencies in second",
				Buckets: prometheus.DefBuckets,
			},
			[]string{},
		),
		websocketConnectionNumber: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "websocket_connection",
				Help: "How many websocket is connecting to relay",
			},
			[]string{},
		),
	}
}

func (m *metrics) SuccessTotal(lvs ...string) {
	m.successTotal.WithLabelValues(lvs...).Inc()
}

func (m *metrics) FailTotal(lvs ...string) {
	m.failTotal.WithLabelValues(lvs...).Inc()
}

func (m *metrics) ProcessDuration(t time.Time, lvs ...string) {
	m.processDuration.WithLabelValues(lvs...).Observe(float64(time.Since(t)) / float64(time.Second))
}

func (m *metrics) WebsocketConnectionNumber(lvs ...string) {
	m.websocketConnectionNumber.WithLabelValues(lvs...).Inc()
}
