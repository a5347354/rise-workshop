package pkg

import "gopkg.in/olahol/melody.v1"

func NewWebsocket() *melody.Melody {
	m := melody.New()
	m.Config.MaxMessageSize = 2097152
	return m
}
