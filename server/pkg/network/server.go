package network

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rouzbehsbz/manticore/server/pkg/network/session"
)

const (
	ReadBufferSize  = 1024
	WriteBufferSize = 1024
)

type Server struct {
	upgrader        websocket.Upgrader
	sessionsManager *session.SessionManager
}

func NewServer() *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  ReadBufferSize,
			WriteBufferSize: WriteBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		sessionsManager: session.NewSessionManager(),
	}
}

func (s *Server) Listen(addr string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/ws", s.accept)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) accept(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	s.sessionsManager.Insert(conn)
}
