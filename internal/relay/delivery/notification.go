package delivery

import (
	"github.com/a5347354/rise-workshop/internal/relay"
	"github.com/sirupsen/logrus"

	"context"
	"sync"

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

func (n *Notification) Broadcast(ctx context.Context, msg []byte) {
	n.RLock()
	defer n.RUnlock()
	for _, session := range n.sessions {
		err := n.m.BroadcastOthers(msg, session)
		logrus.Warn(err)
	}
}
