package delivery

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

type Metrics interface {
	SuccessTotal(lvs ...string)
	FailTotal(lvs ...string)
	ProcessDuration(t time.Time, lvs ...string)
}

type metrics struct {
	successTotal    *prometheus.CounterVec
	failTotal       *prometheus.CounterVec
	processDuration *prometheus.HistogramVec
}

func NewRelayMetrics() Metrics {
	return &metrics{
		successTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "success_total",
				Help: "How many processed success",
			},
			[]string{"type"},
		),
		failTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "fail_total",
				Help: "How many processed fail, partitioned by type and stage.",
			},
			[]string{"type", "stage"},
		),
		processDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "processed_duration_second",
				Help:    "The processed latencies in second",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"type"},
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
