package combat

import (
	"github.com/rouzbehsbz/manticore/server/internal/common"
	"github.com/rouzbehsbz/manticore/server/internal/gameplay/character"
	"github.com/rouzbehsbz/manticore/server/pkg/network/protocol"
	"github.com/rouzbehsbz/manticore/server/pkg/network/session"
	"github.com/rouzbehsbz/zurvan"
)

type CastSpellHandler struct {
	world *zurvan.World
}

func NewCastSpellHandler(world *zurvan.World) *CastSpellHandler {
	return &CastSpellHandler{
		world: world,
	}
}

func (c *CastSpellHandler) Handle(rp session.ReceivedPacket) {
	if !rp.Session.IsAuthenticated() || !rp.Session.IsCharacterSelected() {
		common.ErrorRes(rp.Session, common.UnauthorizedErrorMsg)
		return
	}

	characters, _ := character.CharacterEntityMap(c.world)
	spells, _ := SpellsMap(c.world)

	payload := rp.Packet.Payload.(*protocol.Packet_CastSpellReq)

	spellId := payload.CastSpellReq.SpellId
	targetId := payload.CastSpellReq.TargetId

	spell, ok := spells.Get(spellId)
	if !ok {
		common.ErrorRes(rp.Session, "Invalid spell.")
		return
	}

	target, ok := characters.Get(targetId)
	if !ok {
		common.ErrorRes(rp.Session, "Invalid target.")
		return
	}

	caster, ok := characters.Get(rp.Session.CharacterId)
	if !ok {
		common.ErrorRes(rp.Session, "Invalid caster.")
		return
	}

	c.world.EmitEvents(
		CastSpellEvent{
			Caster: caster,
			Target: target,
			Spell:  spell,
		},
	)
}
