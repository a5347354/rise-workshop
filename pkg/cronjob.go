package pkg

import (
	"context"

	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
)

type cronjob struct {
	cron *cron.Cron
}

type Cronjob interface {
	AddFunc(spec string, cmd func()) (int, error)
}

func NewCronjob(lc fx.Lifecycle) Cronjob {
	c := &cronjob{cron: cron.New()}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			c.cron.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			c.cron.Stop()
			return nil
		},
	})
	return c
}

func (c *cronjob) AddFunc(spec string, cmd func()) (int, error) {
	entryID, err := c.cron.AddFunc(spec, cmd)
	return int(entryID), err
}
