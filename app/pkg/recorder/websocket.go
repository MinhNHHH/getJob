package recorder

import (
	"log"
	"time"

	"github.com/MinhNHHH/get-job/pkg/message"

	"github.com/gorilla/websocket"
)

type WebSocket struct {
	*websocket.Conn
	In             chan message.Wrapper
	Out            chan message.Wrapper
	lastActiveTime time.Time
	active         bool
}

func NewWebSocketConnection(url string) (*WebSocket, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		return nil, err
	}

	return &WebSocket{
		Conn:   conn,
		In:     make(chan message.Wrapper),
		Out:    make(chan message.Wrapper),
		active: true,
	}, nil
}

func (ws *WebSocket) Start() {
	// Receive message coroutine
	go func() {
		for {
			msg, ok := <-ws.Out
			ws.lastActiveTime = time.Now()
			if ok {
				err := ws.WriteJSON(msg)
				if err != nil {
					log.Printf("Failed to send message: %s", err)
					ws.Stop()
				}
			}
		}
	}()

	// Send message coroutine
	for {
		msg := message.Wrapper{}
		err := ws.ReadJSON(&msg)
		if err == nil {
			ws.In <- msg
		} else {
			log.Printf("Failed to read message. Closing connection: %s", err)
			ws.Stop()
		}
	}
}

func (ws *WebSocket) Stop() {
	if ws.active {
		ws.active = false
		log.Printf("Closing client")
		ws.WriteControl(websocket.CloseMessage, []byte{}, time.Time{})
		close(ws.In)
		close(ws.Out)
		ws.Close()
	}
}
