package gameplay

import "github.com/rouzbehsbz/manticore/server/pkg/network/session"

type PacketHandler interface {
	Handle(receivedPacket session.ReceivedPacket)
}

type Dispatcher struct {
	handlers map[uint8]PacketHandler
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[uint8]PacketHandler),
	}
}

func (d *Dispatcher) Register(packetId uint8, handler PacketHandler) {
	d.handlers[packetId] = handler
}

func (d *Dispatcher) Dispatch(receivedPacket session.ReceivedPacket) bool {
	handler, ok := d.handlers[uint8(receivedPacket.Packet.Id)]
	if !ok {
		return false
	}

	handler.Handle(receivedPacket)
	return true
}
