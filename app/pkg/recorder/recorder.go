package recorder

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/MinhNHHH/get-job/pkg/message"

	cfg "github.com/MinhNHHH/get-job/pkg/cfgs"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type Client struct {
	RecordChannel *webrtc.DataChannel
	conn          *webrtc.PeerConnection
	authenticated bool
}

type Recorder struct {
	wsConn *WebSocket
	lock   *sync.Mutex

	clients map[string]*Client
}

func New() *Recorder {
	return &Recorder{
		clients: make(map[string]*Client),
	}
}

func (rc *Recorder) Start(server string) (string, error) {
	sessionId := uuid.NewString()
	log.Printf("New session: %s", sessionId)

	wsURL := GetWSURL(server, sessionId)
	fmt.Println(wsURL)
	wsConn, err := NewWebSocketConnection(wsURL)
	if err != nil {
		log.Printf("Failed to connect to signaling server: %s", err)
		return "", err
	}
	rc.wsConn = wsConn
	go rc.wsConn.Start()

	// send a ping message to keep websocket alive, doesn't expect to receive anything
	// This messages is expected to be broadcast to all client's connections so it keeps them alive too
	go func() {
		for range time.Tick(5 * time.Second) {
			payload := message.Wrapper{
				Type: "Ping",
				Data: []byte{},
			}
			rc.writeWebsocket(payload)
		}
	}()

	rc.wsConn.SetPingHandler(func(appData string) error {
		return rc.wsConn.WriteControl(websocket.PongMessage, []byte{}, time.Time{})
	})

	rc.wsConn.SetCloseHandler(func(code int, text string) error {
		log.Printf("WebSocket connection closed with code %d :%s", code, text)
		rc.Stop("WebSocket connection to server is closed")
		return nil
	})

	return wsURL, nil
}

func (rc *Recorder) Stop(msg string) {
	if rc.wsConn != nil {
		rc.wsConn.WriteControl(websocket.CloseMessage, []byte{}, time.Time{})
	}
}

func (rc *Recorder) writeWebsocket(msg message.Wrapper) error {
	msg.From = cfg.RECORDER_WEBSOCKET_HOST_ID
	if rc.wsConn == nil {
		return fmt.Errorf("Websocket not connected")
	}
	rc.wsConn.Out <- msg
	return nil
}

func (rc *Recorder) startHandleWsMessage() error {
	if rc.wsConn == nil {
		log.Printf("Websocket connection is not connected")
		return fmt.Errorf("Websocket connection is not connected")
	}

	for {
		msg, ok := <-rc.wsConn.In
		if !ok {
			log.Printf("Failed to read websocket message")
			return fmt.Errorf("Failed to read websocket message")
		}

		// skip message that are not send to the host
		if msg.To != cfg.RECORDER_WEBSOCKET_HOST_ID {
			log.Printf("Skip message :%s", msg)
			continue
		}
		err := rc.HandleWebSocketMessage(msg)
		if err != nil {
			log.Printf("Failed to handle message: %v, with error: %s", msg, err)
			continue
		}
	}
}

func (rc *Recorder) HandleWebSocketMessage(msg message.Wrapper) error {
	if msg.Type == message.TCConnect {
		clientVersion := msg.Data.(string)
		if clientVersion != cfg.SUPPORTED_VERSION {
			rc.writeWebsocket(message.Wrapper{Type: message.TCUnsupportedVersion, Data: cfg.SUPPORTED_VERSION, To: msg.From})
			return fmt.Errorf("Client is running unsupported version: %s", clientVersion)
		}

		_, err := rc.newClient(msg.From)
		log.Printf("New client with ID: %s", msg.From)
		if err != nil {
			return fmt.Errorf("Failed to create client: %s", err)
		}
		msg := message.Wrapper{
			To: msg.From,
		}
		rc.writeWebsocket(msg)
		return nil
	}

	client, ok := rc.clients[msg.From]
	if !ok {
		return fmt.Errorf("Client with ID: %s is not found", msg.From)
	}

	switch msgType := msg.Type; msgType {
	case message.TRTCOffer:
		offer := webrtc.SessionDescription{}
		if err := json.Unmarshal([]byte(msg.Data.(string)), &offer); err != nil {
			return err
		}
		log.Printf("Get an offer: %v", (string(msg.Data.(string))))

		if err := client.conn.SetRemoteDescription(offer); err != nil {
			return fmt.Errorf("Failed to set remote description: %s", err)
		}

		// send back SDP answer and set it as local description
		answer, err := client.conn.CreateAnswer(nil)
		if err != nil {
			return fmt.Errorf("Failed to create offfer: %s", err)
		}

		if err := client.conn.SetLocalDescription(answer); err != nil {
			return fmt.Errorf("Failed to set local description: %s", err)
		}
		answerByte, _ := json.Marshal(answer)
		payload := message.Wrapper{
			Type: message.TRTCAnswer,
			Data: string(answerByte),
			To:   msg.From,
		}
		rc.writeWebsocket(payload)
	case message.TRTCCandidate:
		candidate := webrtc.ICECandidateInit{}
		if err := json.Unmarshal([]byte(msg.Data.(string)), &candidate); err != nil {
			return fmt.Errorf("Failed to unmarshall icecandidate: %s", err)
		}

		if err := client.conn.AddICECandidate(candidate); err != nil {
			return fmt.Errorf("Failed to add ice candidate: %s", err)
		}
	default:
		return fmt.Errorf("Not implemented to handle message type: %s", msg.Type)
	}
	return nil
}

func (rc *Recorder) removeClient(id string) {
	if client, ok := rc.clients[id]; ok {
		rc.lock.Lock()
		defer rc.lock.Unlock()

		if client.RecordChannel != nil {
			client.RecordChannel.Close()
			client.RecordChannel = nil
		}

		if client.conn != nil {
			client.conn.Close()
		}

		delete(rc.clients, id)
	}
}

func (rc *Recorder) newClient(id string) (*Client, error) {
	// Initiate peer connection
	ICEServers := cfg.TRANSFER_ICE_SERVER_STUNS

	var cfgs = webrtc.Configuration{
		ICEServers: ICEServers,
	}

	client := &Client{
		authenticated: false,
	}

	rc.lock.Lock()
	rc.clients[id] = client
	rc.lock.Unlock()

	peerConn, err := webrtc.NewPeerConnection(cfgs)
	if err != nil {
		fmt.Printf("Failed to create peer connection: %s", err)
		return nil, err
	}
	client.conn = peerConn

	// Create channel that is blocked until ICE Gathering is complete
	peerConn.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		log.Printf("Peer connection state has changed: %s", s.String())
		switch s {
		// case webrtc.PeerConnectionStateConnected:
		case webrtc.PeerConnectionStateClosed, webrtc.PeerConnectionStateDisconnected:
			log.Printf("Removing client: %s", id)
			rc.removeClient(id)
		}
	})

	peerConn.OnDataChannel(func(dc *webrtc.DataChannel) {
		log.Printf("New DataChannel %s %d\n", dc.Label(), dc.ID())
		dc.OnOpen(func() {
			log.Printf("New DataChannel %s %d\n", dc.Label(), dc.ID())
			switch label := dc.Label(); label {
			case cfg.RECORDER_WEBRTC_DATA_CHANNEL:
				dc.OnMessage(func(msg webrtc.DataChannelMessage) {
					log.Printf("Received message: %s", msg.Data)
					dc.Send([]byte("aasdhalsdhasljdhlas"))
				})
				rc.clients[id].RecordChannel = dc
			default:
				log.Printf("Unhanle data channel with label: %s", dc.Label())

			}
		})
	})

	peerConn.OnICECandidate(func(ice *webrtc.ICECandidate) {
		if ice == nil {
			return
		}
		candidate, err := json.Marshal(ice.ToJSON())
		if err != nil {
			log.Printf("Failed to decode ice candidate: %s", err)
			return
		}

		msg := message.Wrapper{
			Type: message.TRTCCandidate,
			Data: string(candidate),
			To:   id,
		}
		rc.writeWebsocket(msg)
	})
	return client, nil
}
