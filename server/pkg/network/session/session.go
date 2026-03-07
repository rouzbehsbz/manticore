package session

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rouzbehsbz/manticore/server/pkg/network/protocol"
)

const (
	SendBufferSize = 256
)

type ReceivedPacket struct {
	Packet  *protocol.Packet
	Session *Session
}

type Session struct {
	id uint32

	AccountId   uint32
	CharacterId uint32

	conn    *websocket.Conn
	manager *SessionManager

	frameBuf *protocol.Frame
	mu       sync.Mutex

	send chan protocol.Frame

	closeOnce sync.Once
}

func NewSession(id uint32, conn *websocket.Conn, manager *SessionManager) *Session {
	s := &Session{
		id:          id,
		AccountId:   0,
		CharacterId: 0,
		conn:        conn,
		frameBuf:    protocol.NewFrame(),
		mu:          sync.Mutex{},
		send:        make(chan protocol.Frame, SendBufferSize),
		manager:     manager,
		closeOnce:   sync.Once{},
	}

	go s.receiveLoop()
	go s.sendLoop()

	return s
}

func (s *Session) Write(packet *protocol.Packet) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.frameBuf.Append(packet)
}

func (s *Session) IsAuthenticated() bool {
	return s.AccountId != 0
}

func (s *Session) IsCharacterSelected() bool {
	return s.AccountId != 0
}

func (s *Session) receiveLoop() {
	defer s.close()

	for {
		mType, bytes, err := s.conn.ReadMessage()
		if err != nil {
			break
		}

		if mType != websocket.BinaryMessage {
			continue
		}

		frame, err := protocol.BuildFrame(bytes)
		if err != nil {
			break
		}

		for _, packet := range frame.Packets() {
			isNoneBlocking, ok := protocol.IsNoneBlockingPacketType(uint8(packet.Id))
			if !ok {
				continue
			}

			receivedPacket := ReceivedPacket{
				Packet:  packet,
				Session: s,
			}

			if isNoneBlocking {
				s.manager.NonBlocking <- receivedPacket
				continue
			}

			s.manager.Blocking <- receivedPacket
		}
	}
}

func (s *Session) sendLoop() {
	defer s.close()

	for frame := range s.send {
		bytes, err := frame.Bytes()
		if err != nil {
			break
		}

		err = s.conn.WriteMessage(websocket.BinaryMessage, bytes)
		if err != nil {
			break
		}
	}
}

func (s *Session) flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.frameBuf.Len() == 0 {
		return
	}

	s.send <- *s.frameBuf
	s.frameBuf.Empty()
}

func (s *Session) close() {
	s.closeOnce.Do(func() {
		s.conn.Close()
		close(s.send)
		s.manager.sessions.Delete(s.id)
	})
}
