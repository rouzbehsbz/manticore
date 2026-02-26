package protocol

import "google.golang.org/protobuf/proto"

const (
	LoginRequestPacketId uint8 = iota
	LoginResponsePacketId
	RegisterRequestPacketId
	RegisterResponsePacketId
)

var PacketRegistry = map[uint8]proto.Message{
	LoginRequestPacketId:     &LoginRequest{},
	LoginResponsePacketId:    &LoginResponse{},
	RegisterRequestPacketId:  &RegisterRequest{},
	RegisterResponsePacketId: &RegisterResponse{},
}
