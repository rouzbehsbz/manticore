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
	zurvan.Query2[Health, Mana](w, func(e zurvan.Entity, h *Health, m *Mana) {
		sec := dt.Seconds()

		h.Current += h.Regeneration * sec
		m.Current += m.Regeneration * sec

		m.Current = min(m.Current, m.Max)
		h.Current = min(h.Current, h.Max)
	})
}

type TakeDamageSystem struct{}

func (t *TakeDamageSystem) Update(w *zurvan.World, dt time.Duration) {
	zurvan.Query3[Health, DefensiveStats, DamageTaken](w, func(e zurvan.Entity, h *Health, ds *DefensiveStats, dt *DamageTaken) {
		if rand.Float64() > ds.Evasion {
			return
		}

		reduction := MagicResistanceHardCap * (ds.MagicResistance / (ds.MagicResistance + MagicResistanceConst))
		finalDamage := dt.Amount * (1 - reduction)

		h.Current -= finalDamage
		h.Current = min(h.Current, h.Max)
	})
}

type TakeHealSystem struct{}

func (t *TakeHealSystem) Update(w *zurvan.World, dt time.Duration) {
	zurvan.Query2[Health, HealTaken](w, func(e zurvan.Entity, h *Health, ht *HealTaken) {
		h.Current += ht.Amount
		h.Current = min(h.Current, h.Max)
	})
}
