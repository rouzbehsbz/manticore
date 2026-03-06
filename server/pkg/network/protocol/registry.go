package protocol

import (
	"google.golang.org/protobuf/proto"
)

const (
	ErrorResponsePacketId uint8 = iota
	LoginRequestPacketId
	LoginResponsePacketId
	RegisterRequestPacketId
	RegisterResponsePacketId
	MyCharactersListRequestPacketId
	MyCharactersListResponsePacketId
)

const (
	BlockingPacketType uint8 = iota
	NoneBlockingPacketType
)

var PacketRegistry = map[uint8]proto.Message{
	ErrorResponsePacketId:            &ErrorResponse{},
	LoginRequestPacketId:             &LoginRequest{},
	LoginResponsePacketId:            &LoginResponse{},
	RegisterRequestPacketId:          &RegisterRequest{},
	RegisterResponsePacketId:         &RegisterResponse{},
	MyCharactersListRequestPacketId:  &MyCharactersListRequest{},
	MyCharactersListResponsePacketId: &MyCharactersListResponse{},
}

func BuildErrorResponsePacket(msg string) *Packet {
	return &Packet{
		Id: uint32(ErrorResponsePacketId),
		Payload: &Packet_ErrorResponse{
			ErrorResponse: &ErrorResponse{
				Msg: msg,
			},
		},
	}
}

func BuildRegisterResponsePacket() *Packet {
	return &Packet{
		Id: uint32(RegisterResponsePacketId),
		Payload: &Packet_RegisterResponse{
			RegisterResponse: &RegisterResponse{},
		},
	}
}

func BuildLoginResponsePacket() *Packet {
	return &Packet{
		Id: uint32(LoginResponsePacketId),
		Payload: &Packet_LoginResponse{
			LoginResponse: &LoginResponse{},
		},
	}
}

func BuildMyCharactersListResponsePacket(characters []*MyCharacter) *Packet {
	return &Packet{
		Id: uint32(MyCharactersListResponsePacketId),
		Payload: &Packet_MyCharactersListResponse{
			MyCharactersListResponse: &MyCharactersListResponse{
				Characters: characters,
			},
		},
	}
}

var PacketTypeRegistry = map[uint8]uint8{
	ErrorResponsePacketId:            NoneBlockingPacketType,
	LoginRequestPacketId:             BlockingPacketType,
	LoginResponsePacketId:            BlockingPacketType,
	RegisterRequestPacketId:          BlockingPacketType,
	RegisterResponsePacketId:         BlockingPacketType,
	MyCharactersListRequestPacketId:  BlockingPacketType,
	MyCharactersListResponsePacketId: NoneBlockingPacketType,
}

func IsNoneBlockingPacketType(packetId uint8) (bool, bool) {
	pType, ok := PacketTypeRegistry[packetId]
	if !ok {
		return false, false
	}

	return pType == NoneBlockingPacketType, true
}
