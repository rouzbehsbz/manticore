package combat

import (
	"math/rand"
	"time"

	"github.com/rouzbehsbz/zurvan"
)

type StatCalculationSystem struct{}

func (s *StatCalculationSystem) Update(w *zurvan.World, dt time.Duration) {

}

type RegenerationSystem struct{}

func (r *RegenerationSystem) Update(w *zurvan.World, dt time.Duration) {
	sec := dt.Seconds()

	zurvan.QueryMany2[Health, Mana](w, func(e []zurvan.Entity, h []Health, m []Mana) {
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
		h, ds := zurvan.QueryOne2[Health, DefensiveStats](w, event.Entity)

		if rand.Float64() > ds.Evasion {
			return
		}

		reduction := MagicResistanceHardCap * (ds.MagicResistance / (ds.MagicResistance + MagicResistanceConst))
		finalDamage := event.Amount * (1 - reduction)

		h.Current -= finalDamage
		h.Current = min(h.Current, h.Max)
	}
}

type TakeHealSystem struct{}

func (t *TakeHealSystem) Update(w *zurvan.World, dt time.Duration) {
	events := zurvan.OnEvent[HealTakenEvent](w)

	for _, event := range events {
		h := zurvan.QueryOne1[Health](w, event.Entity)

		h.Current += event.Amount
		h.Current = min(h.Current, h.Max)
	}
}
