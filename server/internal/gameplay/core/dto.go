package core

import (
	"github.com/rouzbehsbz/manticore/server/pkg/network"
	"github.com/rouzbehsbz/manticore/server/pkg/network/session"
	"github.com/rouzbehsbz/zurvan"
)

func SessionManager(w *zurvan.World) *session.SessionManager {
	return zurvan.Resource[*session.SessionManager](w)
}

func NonBlockingPackets(w *zurvan.World) <-chan session.ReceivedPacket {
	return zurvan.Resource[<-chan session.ReceivedPacket](w)
}

func Dispatcher(w *zurvan.World) *network.Dispatcher {
	return zurvan.Resource[*network.Dispatcher](w)
}
