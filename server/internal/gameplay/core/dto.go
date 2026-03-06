package core

import (
	"github.com/rouzbehsbz/manticore/server/pkg/network"
	"github.com/rouzbehsbz/manticore/server/pkg/network/session"
	"github.com/rouzbehsbz/zurvan"
)

func sessionManager(w *zurvan.World) (*session.SessionManager, bool) {
	return zurvan.Resource[*session.SessionManager](w)
}

func nonBlockingPackets(w *zurvan.World) (<-chan session.ReceivedPacket, bool) {
	return zurvan.Resource[<-chan session.ReceivedPacket](w)
}

func dispatcher(w *zurvan.World) (*network.Dispatcher, bool) {
	return zurvan.Resource[*network.Dispatcher](w)
}
