package delivery

import (
	"context"
	"github.com/a5347354/rise-workshop/internal/relay"
	"github.com/a5347354/rise-workshop/pkg"
	"time"
)

type relayCronjobHandler struct {
	usecase relay.Usecase
	metrics Metrics
}

func RegistCronjobHandler(usecase relay.Usecase, metrics Metrics, cronjob pkg.Cronjob) {
	r := &relayCronjobHandler{
		usecase: usecase,
		metrics: metrics,
	}
	cronjob.AddFunc("* * * * *", r.SendToSubscriber)
}

func (h *relayCronjobHandler) SendToSubscriber() {
	for i := 1; i <= 8; i++ {
		go func(h *relayCronjobHandler) {
			t := time.Now()
			err := h.usecase.SendMessageToSubscriber(context.Background())
			if err != nil {
				h.metrics.FailTotal("broadcast", "send")
				return
			}
			h.metrics.SuccessTotal("broadcast")
			h.metrics.ProcessDuration(t, "broadcast")
		}(h)
		time.Sleep(5 * time.Second)
	}
}
