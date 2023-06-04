package delivery

import (
	"github.com/a5347354/rise-workshop/internal/relay"
	"github.com/a5347354/rise-workshop/pkg"
	"github.com/tidwall/gjson"

	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/olahol/melody.v1"
)

type relayHandler struct {
	usecase relay.Usecase
	metrics Metrics
}

func RegistWebsocketHandler(engine *gin.Engine, m *melody.Melody, usecase relay.Usecase, metrics Metrics) {
	r := &relayHandler{
		usecase: usecase,
		metrics: metrics,
	}
	engine.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleConnect(r.handleConnect())
	m.HandleDisconnect(r.handleDisconnect())
	m.HandleError(r.handleError())
	m.HandleMessage(r.message())
}

func (r *relayHandler) handleConnect() func(s *melody.Session) {
	return func(s *melody.Session) {
		fmt.Printf("[Melody] %v | %s | Connect %s\n", time.Now().Format("2006/01/02 - 15:04:05"), s.Request.RemoteAddr, s.Request.RequestURI)

	}
}

func (r *relayHandler) handleDisconnect() func(s *melody.Session) {
	return func(s *melody.Session) {
		fmt.Printf("[Melody] %v | %s | Disconnect %s\n", time.Now().Format("2006/01/02 - 15:04:05"), s.Request.RemoteAddr, s.Request.RequestURI)
	}
}

func (r *relayHandler) handleError() func(s *melody.Session, err error) {
	return func(s *melody.Session, err error) {
		fmt.Printf("[Melody] %v | %s | Error %s %s\n", time.Now().Format("2006/01/02 - 15:04:05"), s.Request.RemoteAddr, s.Request.RequestURI, err)
	}
}

func (r *relayHandler) message() func(s *melody.Session, msg []byte) {
	return func(s *melody.Session, msg []byte) {
		eventType := gjson.GetBytes(msg, "0").String()
		t := time.Now()
		resp, err := r.usecase.ReceiveMessage(context.Background(), msg, s)
		if err != nil {
			r.metrics.FailTotal(eventType, "ReceiveMessage")
			logrus.Panic(err)
		}
		switch resp.Action {
		case pkg.WebSocketMsgTypeNormal:
			s.Write(resp.Msg)
		}
		fmt.Printf("[Melody] %v | %13v | %s | Message %s\n",
			t.Format("2006/01/02 - 15:04:05"),
			time.Since(t),
			s.Request.RemoteAddr,
			s.Request.RequestURI,
		)
		r.metrics.SuccessTotal(eventType)
		r.metrics.ProcessDuration(t, eventType)
	}
}
