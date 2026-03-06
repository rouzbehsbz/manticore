package models

import (
	"github.com/rouzbehsbz/manticore/server/internal/infra/db/sources"
	"github.com/rouzbehsbz/manticore/server/pkg/network/protocol"
)

func MapCharacterToMyCharacterPacket(character sources.Character) *protocol.MyCharacter {
	return &protocol.MyCharacter{
		Id:       uint32(character.ID),
		Nickname: character.Nickname,
		Level:    uint32(character.Level),
	}
}
