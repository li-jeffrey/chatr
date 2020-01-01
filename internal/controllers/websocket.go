package controllers

import (
	"chatr/internal/coordinator"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

const wsPath = "/ws"

var upgrader = websocket.FastHTTPUpgrader{
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true
	},
}

func upgradeConnection(ctx *fasthttp.RequestCtx) {
	err := upgrader.Upgrade(ctx, onUpgrade)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); ok {
			log.Error("Failed to upgrade ws connection: %s", err)
		}
		return
	}
}

func onUpgrade(ws *websocket.Conn) {
	session := coordinator.CreateSession()

	go func(s *coordinator.Session, conn *websocket.Conn) {
	Loop:
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Error("[%s] Failed to read: err=%s", s.ID, err)
				break
			}

			switch mt {
			case websocket.TextMessage:
				log.Debug("[%s] Received: %s", s.ID, message)
				s.Reads <- message
			case websocket.CloseMessage:
				break Loop
			}
		}
		close(s.Reads)
	}(session, ws)

	func(s *coordinator.Session, conn *websocket.Conn) {
		for msg := range s.Writes {
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Error("[%s] Failed to write: err=%s", s.ID, err)
				break
			}
			log.Debug("[%s] Sent: %s", s.ID, msg)
		}
		conn.Close()
	}(session, ws)
}
