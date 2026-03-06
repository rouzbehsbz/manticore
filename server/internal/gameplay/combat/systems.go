package combat

import (
	"math"
	"math/rand"
	"time"

	"github.com/rouzbehsbz/manticore/server/internal/gameplay/character"
	"github.com/rouzbehsbz/zurvan"
)

type StatCalculationSystem struct{}

func (s *StatCalculationSystem) Update(w *zurvan.World, dt time.Duration) {
	events := zurvan.OnEvent[RecalculateStatsEvent](w)

	for _, event := range events {
		h, m, o, d, p := zurvan.QueryOne5[character.Health, character.Mana, character.OffensiveStats, character.DefensiveStats, character.PrimaryStats](w, event.Entity)

		h.Max = 100 + p.Vitality*10 + math.Sqrt(p.Spirit)*5
		h.Regeneration = 0.1*math.Sqrt(p.Vitality) + 0.3*math.Sqrt(p.Spirit)

		m.Max = 50 + p.Intelligence*12 + p.Willpower*5
		m.Regeneration = 0.2*math.Sqrt(p.Willpower) + 0.2*math.Sqrt(p.Spirit)

		d.MagicResistance = 0.5 * p.Willpower
		d.Evasion = 0.5 * p.Dexterity

		o.SpellPower = p.Intelligence*2 + p.Willpower*0.5
		o.CriticalRating = p.Dexterity + 0.5*p.Spirit
		o.Accuracy = p.Dexterity
	}
}

type RegenerationSystem struct{}

func (r *RegenerationSystem) Update(w *zurvan.World, dt time.Duration) {
	sec := dt.Seconds()

	zurvan.QueryMany2[character.Health, character.Mana](w, func(e []zurvan.Entity, h []character.Health, m []character.Mana) {
		for i, _ := range e {
			h[i].Current += h[i].Regeneration * sec
			m[i].Current += m[i].Regeneration * sec

			m[i].Current = min(m[i].Current, m[i].Max)
			h[i].Current = min(h[i].Current, h[i].Max)
		}
	})
}

type TakeDamageSystem struct{}

func (t *TakeDamageSystem) Update(w *zurvan.World, dt time.Duration) {
	events := zurvan.OnEvent[DamageTakenEvent](w)

	for _, event := range events {
		h, ds := zurvan.QueryOne2[character.Health, character.DefensiveStats](w, event.Entity)

		evasionPerc := 0.5 * ds.Evasion / (ds.Evasion + 100)
		if rand.Float64() > evasionPerc {
			continue
		}

		reduction := 0.5 * (ds.MagicResistance / (ds.MagicResistance + 200))
		finalDamage := event.Amount * (1 - reduction)

		h.Current -= finalDamage
		h.Current = min(h.Current, h.Max)
	}
}

type TakeHealSystem struct{}

func (t *TakeHealSystem) Update(w *zurvan.World, dt time.Duration) {
	events := zurvan.OnEvent[HealTakenEvent](w)

	for _, event := range events {
		h := zurvan.QueryOne1[character.Health](w, event.Entity)

		h.Current += event.Amount
		h.Current = min(h.Current, h.Max)
	}
}
