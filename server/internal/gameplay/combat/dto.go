package combat

import (
	"time"

	"github.com/rouzbehsbz/manticore/server/internal/models"
	"github.com/rouzbehsbz/manticore/server/pkg/util"
	"github.com/rouzbehsbz/zurvan"
)

func Spells(w *zurvan.World) (*util.SyncMap[uint32, models.Spell], bool) {
	return zurvan.Resource[*util.SyncMap[uint32, models.Spell]](w)
}

type CastSpellEvent struct {
	Caster zurvan.Entity
	Target zurvan.Entity
	Spell  models.Spell
}

type CancelCastSpellEvent struct {
	Caster zurvan.Entity
}

type FireSpellEvent struct {
	Caster zurvan.Entity
	Target zurvan.Entity
	Spell  models.Spell
}

type TakeDamageEvent struct {
	Target zurvan.Entity
	Amount float64
}

type TakeHealEvent struct {
	Target zurvan.Entity
	Amount float64
}

type TakingOverTime struct {
	SpellEffect   models.SpellEffect
	RemainingTime time.Duration
}

type CastingSpell struct {
	Target        zurvan.Entity
	Spell         models.Spell
	RemainingTime time.Duration
}
