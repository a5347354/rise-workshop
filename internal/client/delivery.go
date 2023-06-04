package client

import "time"

type Metrics interface {
	SuccessTotal(lvs ...string)
	FailTotal(lvs ...string)
	ProcessDuration(t time.Time, lvs ...string)
	WebsocketConnectionNumber(...string)
}
