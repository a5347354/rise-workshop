package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/nbd-wtf/go-nostr"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func main() {
	relayURL := viper.GetString("relay.url")
	pb := viper.GetString("public.key")
	pub := viper.GetString("private.key")
	relay, err := nostr.RelayConnect(context.Background(), relayURL)
	defer relay.Close()
	if err != nil {
		panic(err)
	}

	e := nostr.Event{
		PubKey:    pub,
		CreatedAt: nostr.Now(),
		Kind:      31337,
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
	e.ID = e.GetID()
	e.Sign(pb)
	_, err = e.CheckSignature()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(e.ID)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	status, err := relay.Publish(ctx, e)
	if err != nil {
		logrus.Panic(err)
	}
	if status != nostr.PublishStatusSucceeded {
		logrus.Panic(status)
	}
}
