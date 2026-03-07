package protocol

import (
	"google.golang.org/protobuf/proto"
)

const (
	ErrorResPacketId uint8 = iota
	LoginReqPacketId
	LoginResPacketId
	RegisterReqPacketId
	RegisterResPacketId
	MyCharactersListReqPacketId
	MyCharactersListResPacketId
	CastSpellReqPacketId
	CastSpellResPacketId
)

const (
	BlockingPacketType uint8 = iota
	NoneBlockingPacketType
)

var PacketRegistry = map[uint8]proto.Message{
	ErrorResPacketId:            &ErrorRes{},
	LoginReqPacketId:            &LoginReq{},
	LoginResPacketId:            &LoginRes{},
	RegisterReqPacketId:         &RegisterReq{},
	RegisterResPacketId:         &RegisterRes{},
	MyCharactersListReqPacketId: &MyCharactersListReq{},
	MyCharactersListResPacketId: &MyCharactersListRes{},
	CastSpellReqPacketId:        &CastSpellReq{},
	CastSpellResPacketId:        &CastSpellReq{},
}

func BuildErrorResPacket(msg string) *Packet {
	return &Packet{
		Id: uint32(ErrorResPacketId),
		Payload: &Packet_ErrorRes{
			ErrorRes: &ErrorRes{
				Msg: msg,
			},
		},
	}
}

func BuildRegisterResPacket() *Packet {
	return &Packet{
		Id: uint32(RegisterResPacketId),
		Payload: &Packet_RegisterRes{
			RegisterRes: &RegisterRes{},
		},
	}
}

func BuildLoginResPacket() *Packet {
	return &Packet{
		Id: uint32(LoginResPacketId),
		Payload: &Packet_LoginRes{
			LoginRes: &LoginRes{},
		},
	}
}

func BuildMyCharactersListResPacket(characters []*MyCharacter) *Packet {
	return &Packet{
		Id: uint32(MyCharactersListResPacketId),
		Payload: &Packet_MyCharactersListRes{
			MyCharactersListRes: &MyCharactersListRes{
				Characters: characters,
			},
		},
	}
}

var PacketTypeRegistry = map[uint8]uint8{
	ErrorResPacketId:            NoneBlockingPacketType,
	LoginReqPacketId:            BlockingPacketType,
	LoginResPacketId:            BlockingPacketType,
	RegisterReqPacketId:         BlockingPacketType,
	RegisterResPacketId:         BlockingPacketType,
	MyCharactersListReqPacketId: BlockingPacketType,
	MyCharactersListResPacketId: BlockingPacketType,
	CastSpellReqPacketId:        NoneBlockingPacketType,
	CastSpellResPacketId:        NoneBlockingPacketType,
}

func IsNoneBlockingPacketType(packetId uint8) (bool, bool) {
	pType, ok := PacketTypeRegistry[packetId]
	if !ok {
		return false, false
	}

	return pType == NoneBlockingPacketType, true
}
