package character

import (
	"context"
	"math"
	"time"

	"github.com/rouzbehsbz/manticore/server/internal/infra/db"
	"github.com/rouzbehsbz/zurvan"
)

type LoadSystem struct {
	db *db.Db
}

func NewLoadSystem(db *db.Db) *LoadSystem {
	return &LoadSystem{
		db: db,
	}
}

func (l *LoadSystem) Update(w *zurvan.World, dt time.Duration) {
	characters, _ := CharactersMap(w)

	rawCharacters, err := l.db.Q.GetAllCharacters(context.Background())
	if err != nil {
		panic(err.Error())
	}

	for _, character := range rawCharacters {
		characters.Insert(uint32(character.ID), character)
	}
}

type JoinWorldSystem struct{}

func (j *JoinWorldSystem) Update(w *zurvan.World, dt time.Duration) {
	characters, _ := CharactersMap(w)

	events := zurvan.OnEvent[JoinsWorldEvent](w)

	for _, event := range events {
		char, _ := characters.Get(event.Id)

		w.PushCommands(
			zurvan.NewSetComponentsCommand(event.Character,
				NewLevel(int(char.Level), int(char.Xp)),
				PrimaryStats{
					Vitality:     float64(char.Vitality),
					Intelligence: float64(char.Intelligence),
					Willpower:    float64(char.Willpower),
					Dexterity:    float64(char.Dexterity),
					Spirit:       float64(char.Spirit),
				},
				Health{},
				Mana{},
				DefensiveStats{},
				OffensiveStats{},
			),
		)

		w.EmitEvents(
			RecalculateStatsEvent{
				Entity: event.Character,
			},
		)
	}
}

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
		lvl.NextLevelXpNeeded = int(float64(100) * math.Pow(float64(lvl.Value), 1.5))

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
