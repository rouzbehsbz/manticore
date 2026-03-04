package session

import (
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/rouzbehsbz/manticore/server/pkg/util"
)

const (
	BlockingBufferSize     = 256
	NoneBlockingBufferSize = 256
)

type SessionManager struct {
	sessions  *util.SyncMap[uint32, *Session]
	idCounter atomic.Uint32

	Blocking    chan ReceivedPacket
	NonBlocking chan ReceivedPacket
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions:    util.NewSyncMap[uint32, *Session](),
		idCounter:   atomic.Uint32{},
		Blocking:    make(chan ReceivedPacket, BlockingBufferSize),
		NonBlocking: make(chan ReceivedPacket, NoneBlockingBufferSize),
	}
}

func (s *SessionManager) Insert(conn *websocket.Conn) {
	id := s.idCounter.Add(1)
	session := NewSession(id, conn, s)

	s.sessions.Insert(id, session)
}

func (s *SessionManager) FlushAll() {
	s.sessions.Iter(func(k uint32, v *Session) {
		v.flush()
	})
}
