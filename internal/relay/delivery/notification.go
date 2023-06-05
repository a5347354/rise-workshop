package delivery

import (
	"github.com/a5347354/rise-workshop/internal/relay"

	"context"
	"encoding/json"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/olahol/melody.v1"
)

type Notification struct {
	m        *melody.Melody
	sessions map[string]*melody.Session
	*sync.RWMutex
}

func NewNotification(m *melody.Melody) relay.Notification {
	return &Notification{m, make(map[string]*melody.Session), &sync.RWMutex{}}
}

func (n *Notification) Subscribe(ctx context.Context, id string, s *melody.Session) {
	n.Lock()
	defer n.Unlock()
	n.sessions[id] = s
}

func (n *Notification) UnSubscribe(ctx context.Context, id string) {
	n.Lock()
	defer n.Unlock()
	delete(n.sessions, id)
}

func (n *Notification) Broadcast(ctx context.Context, msg []interface{}) {
	n.RLock()
	defer n.RUnlock()
	for id := range n.sessions {
		msg[1] = id
		b, err := json.Marshal(msg)
		if err != nil {
			logrus.Warn(err)
		}
		err = n.m.Broadcast(b)
		if err != nil {
			logrus.Warn(err)
		}
	}
}
