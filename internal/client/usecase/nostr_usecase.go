package usecase

import (
	"github.com/a5347354/rise-workshop/internal/client"
	"github.com/a5347354/rise-workshop/internal/event"
	"github.com/a5347354/rise-workshop/pkg"

	"context"
	"fmt"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type clientUsecase struct {
	client  pkg.NostrClient
	eStore  event.AsyncStore
	metrics client.Metrics
}

func NewClient(lc fx.Lifecycle, eStore event.AsyncStore, metrics client.Metrics) client.Usecase {
	return &clientUsecase{
		client:  pkg.NewNostrClient(lc),
		eStore:  eStore,
		metrics: metrics,
	}
}

func (c clientUsecase) Collect(ctx context.Context, url string) error {
	err := c.client.ConnectURL(ctx, url)
	if err != nil {
		logrus.WithField("action", "sleep 30 sec").Warn(err)
		c.metrics.FailTotal("connect")
		time.Sleep(30 * time.Second)
		return err
	}
	sub, err := c.client.Subscribe(ctx, nostr.Filters{nostr.Filter{Kinds: []int{nostr.KindTextNote}}})
	if err != nil {
		c.metrics.FailTotal("subscribe")
		logrus.Error(err)
		return err
	}
	logrus.Info(fmt.Sprintf("[START] Collect from: %s", url))
	for ev := range sub.Events {
		t := time.Now()
		logrus.Info(*ev)
		err := c.eStore.Insert(ctx, pkg.NostrEventToEvent(*ev))
		if err != nil {
			c.metrics.FailTotal("insert")
			logrus.Error(err)
		}
		c.metrics.SuccessTotal()
		c.metrics.ProcessDuration(t)
	}
	return nil
}

func (c clientUsecase) SendMessage(ctx context.Context) error {
	err := c.client.Connect(ctx)
	if err != nil {
		logrus.Panic(err)
	}

	e := nostr.Event{
		Kind: 31337,
		Tags: nostr.Tags{
			nostr.Tag{
				"d",
				"b07v7s2ic0haospgmeg73i",
			},
			nostr.Tag{
				"media",
				"https://media.zapstr.live:3118/d91191e30e00444b942c0e82cad470b32af171764c2275bee0bd99377efd4075/naddr1qqtxyvphwcmhxvnfvvcxsct0wdcxwmt9vumnx6gzyrv3ry0rpcqygju59s8g9jk5wzej4ut3wexzyad7uz7ejdm7l4q82qcyqqq856g4xnp7j",
				"http",
			},
			nostr.Tag{
				"p",
				"d91191e30e00444b942c0e82cad470b32af171764c2275bee0bd99377efd4075",
				"Host",
			},
			nostr.Tag{
				"p",
				"fa984bd7dbb282f07e16e7ae87b26a2a7b9b90b7246a44771f0cf5ae58018f52",
				"Guest",
			},
			nostr.Tag{
				"c",
				"Podcast",
			},
			nostr.Tag{
				"price",
				"402",
			},
			nostr.Tag{
				"cover",
				"https://s3-us-west-2.amazonaws.com/anchor-generated-image-bank/production/podcast_uploaded_nologo400/36291377/36291377-1673187804611-64b4f8e9f1687.jpg",
			},
			nostr.Tag{
				"subject",
				"Nostrovia | The Pablo Episode",
			},
		},
		Content: "Nostrovia | The Pablo Episode\n\nhttps://s3-us-west-2.amazona",
	}
	c.client.Subscribe(ctx, []nostr.Filter{{
		Kinds:   []int{nostr.KindTextNote},
		Authors: []string{"sss"},
	}})
	_, err = c.client.Publish(ctx, e)
	return err
}
