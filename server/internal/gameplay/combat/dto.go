package combat

import (
	"time"

	"github.com/rouzbehsbz/manticore/server/internal/models"
	"github.com/rouzbehsbz/manticore/server/pkg/util"
	"github.com/rouzbehsbz/zurvan"
)

func SpellsMap(w *zurvan.World) (*util.SyncMap[uint32, models.Spell], bool) {
	return zurvan.Resource[*util.SyncMap[uint32, models.Spell]](w)
}

type TakingOverTime struct {
	SpellEffectType models.EffectType
	Amount          float64
	RemainingTime   time.Duration
}

type CastingSpell struct {
	Target        zurvan.Entity
	Spell         models.Spell
	RemainingTime time.Duration
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
	Caster      zurvan.Entity
	Target      zurvan.Entity
	SpellEffect models.SpellEffect
}

type TakeDamageEvent struct {
	Target zurvan.Entity
	Amount float64
}

type TakeHealEvent struct {
	Target zurvan.Entity
	Amount float64
}

type DeathEvent struct {
	Entity zurvan.Entity
}
