package character

import (
	"math"
	"time"

	"github.com/rouzbehsbz/zurvan"
)

type ExperienceSystem struct{}

func (e *ExperienceSystem) Update(w *zurvan.World, dt time.Duration) {
	events := zurvan.OnEvent[XpGainedEvent](w)

	for _, event := range events {
		lvl := zurvan.QueryOne1[Level](w, event.Entity)

		lvl.Xp += event.Amount

		if lvl.Xp >= lvl.NextLevelXpNeeded {
			w.EmitEvents(
				LevelUpEvent{
					Entity: event.Entity,
				},
			)
		}
	}
}

type LevelUpSystem struct{}

func (l *LevelUpSystem) Update(w *zurvan.World, dt time.Duration) {
	events := zurvan.OnEvent[LevelUpEvent](w)

	for _, event := range events {
		lvl, ps := zurvan.QueryOne2[Level, PrimaryStats](w, event.Entity)

		ps.Vitality += 1
		ps.Intelligence += 1
		ps.Dexterity += 1
		ps.Spirit += 1
		ps.Willpower += 1

		lvl.Value += 1
		lvl.Xp = lvl.Xp - lvl.NextLevelXpNeeded
		lvl.NextLevelXpNeeded = lvl.Value * 100

		w.EmitEvents(
			RecalculateStatsEvent{
				Entity: event.Entity,
			},
		)
	}
}

type StatCalculationSystem struct{}

func (s *StatCalculationSystem) Update(w *zurvan.World, dt time.Duration) {
	events := zurvan.OnEvent[RecalculateStatsEvent](w)

	for _, event := range events {
		h, m, o, d, p := zurvan.QueryOne5[Health, Mana, OffensiveStats, DefensiveStats, PrimaryStats](w, event.Entity)

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
