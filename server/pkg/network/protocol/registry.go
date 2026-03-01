package protocol

import "google.golang.org/protobuf/proto"

const (
	LoginRequestPacketId uint8 = iota
	LoginResponsePacketId
	RegisterRequestPacketId
	RegisterResponsePacketId
)

const (
	BlockingPacketType uint8 = iota
	NoneBlockingPacketType
)

var PacketRegistry = map[uint8]proto.Message{
	LoginRequestPacketId:     &LoginRequest{},
	LoginResponsePacketId:    &LoginResponse{},
	RegisterRequestPacketId:  &RegisterRequest{},
	RegisterResponsePacketId: &RegisterResponse{},
}

var PacketTypeRegistry = map[uint8]uint8{
	LoginRequestPacketId:     BlockingPacketType,
	LoginResponsePacketId:    BlockingPacketType,
	RegisterRequestPacketId:  BlockingPacketType,
	RegisterResponsePacketId: BlockingPacketType,
}

func IsNoneBlockingPacketType(packetId uint8) (bool, bool) {
	pType, ok := PacketTypeRegistry[packetId]
	if !ok {
		return false, false
	}

	return pType == NoneBlockingPacketType, true
}
