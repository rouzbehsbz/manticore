package character

import (
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
		lvl.NextLevelXpNeeded = 0
	}
}
