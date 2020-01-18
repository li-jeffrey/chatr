package controllers

import (
	"chatr/internal/sessions"

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
	sessionID := ctx.QueryArgs().Peek("sessionID")
	err := upgrader.Upgrade(ctx, func(ws *websocket.Conn) { onUpgrade(ws, sessionID) })
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); ok {
			log.Error("Failed to upgrade ws connection: %s", err)
		}
		return
	}
}

func onUpgrade(ws *websocket.Conn, sessionID []byte) {
	session := sessions.GetSession(sessionID)

	go func(s *sessions.Session, conn *websocket.Conn) {
	Loop:
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Error("[%s] Failed to read: err=%s", s.ID, err)
				break Loop
			}

			switch mt {
			case websocket.TextMessage:
				log.Debug("[%s] Received: %s", s.ID, message)
				s.Reads <- message
			}
		}
		close(s.Reads)
	}(session, ws)

	func(s *sessions.Session, conn *websocket.Conn) {
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
