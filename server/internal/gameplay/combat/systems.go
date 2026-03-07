package combat

import (
	"math/rand"
	"time"

	"github.com/rouzbehsbz/manticore/server/internal/gameplay/character"
	"github.com/rouzbehsbz/manticore/server/internal/models"
	"github.com/rouzbehsbz/zurvan"
)

type CastSpellSystem struct{}

func (c *CastSpellSystem) Update(w *zurvan.World, dt time.Duration) {
	events := zurvan.OnEvent[CastSpellEvent](w)

	for _, event := range events {
		c := zurvan.QueryOne1[CastingSpell](w, event.Caster)
		if c != nil {
			continue
		}

		m := zurvan.QueryOne1[character.Mana](w, event.Caster)

		if event.Spell.ManaCost > m.Current {
			continue
		}

		m.Current -= event.Spell.ManaCost

		w.PushCommands(
			zurvan.NewSetComponentsCommand(event.Caster,
				CastingSpell{
					Target:        event.Target,
					Spell:         event.Spell,
					RemainingTime: event.Spell.CastTime,
				},
			),
		)
	}
}

type CastingSpellSystem struct{}

func (c *CastingSpellSystem) Update(w *zurvan.World, dt time.Duration) {
	events := zurvan.OnEvent[CancelCastSpellEvent](w)

	for _, event := range events {
		w.PushCommands(
			zurvan.NewDeleteComponentsCommand(event.Caster,
				CastingSpell{},
			),
		)
	}

	zurvan.QueryMany1[CastingSpell](w, func(e []zurvan.Entity, cs []CastingSpell) {
		for i, _ := range e {
			cs[i].RemainingTime -= dt

			if cs[i].RemainingTime <= 0 {
				w.PushCommands(
					zurvan.NewDeleteComponentsCommand(e[i],
						CastingSpell{},
					),
				)

				// set specific spell cooldown and all spells cooldown

				w.EmitEvents(
					FireSpellEvent{
						Caster:      e[i],
						Target:      cs[i].Target,
						SpellEffect: cs[i].Spell.Effect,
					},
				)
			}
		}
	})
}

type FireSpellSystem struct{}

func (f *FireSpellSystem) Update(w *zurvan.World, dt time.Duration) {
	events := zurvan.OnEvent[FireSpellEvent](w)

	// maybe need to generate the spell entity in the map and
	// calculate the distance before firing

	for _, event := range events {
		os := zurvan.QueryOne1[character.OffensiveStats](w, event.Caster)

		accuracyPerc := 0.5 * os.Accuracy / (os.Accuracy + 100)
		if rand.Float64() > accuracyPerc {
			continue
		}

		amount := event.SpellEffect.Amount + (os.SpellPower * event.SpellEffect.Coefficient)

		criticalPerc := 0.5 * os.CriticalRating / (os.CriticalRating + 100)
		if rand.Float64() <= criticalPerc {
			amount *= 1.5
		}

		switch event.SpellEffect.Type {
		case models.DamageEffect:
			w.EmitEvents(
				TakeDamageEvent{
					Target: event.Target,
					Amount: amount,
				},
			)

		case models.HealEffect:
			w.EmitEvents(
				TakeHealEvent{
					Target: event.Target,
					Amount: amount,
				},
			)

		case models.DamageOverTimeEffect:
			w.PushCommands(
				zurvan.NewSetComponentsCommand(
					event.Target,
					TakingOverTime{
						SpellEffectType: event.SpellEffect.Type,
						Amount:          amount,
						RemainingTime:   event.SpellEffect.Duration,
					},
				),
			)

		case models.HealOverTimeEffect:
			w.PushCommands(
				zurvan.NewSetComponentsCommand(
					event.Target,
					TakingOverTime{
						SpellEffectType: event.SpellEffect.Type,
						Amount:          amount,
						RemainingTime:   event.SpellEffect.Duration,
					},
				),
			)
		}
	}
}

type TakeOverTimeSystem struct{}

func (t *TakeOverTimeSystem) Update(w *zurvan.World, dt time.Duration) {
	zurvan.QueryMany1[TakingOverTime](w, func(e []zurvan.Entity, tot []TakingOverTime) {
		for i, _ := range e {
			switch tot[i].SpellEffectType {
			case models.DamageOverTimeEffect:
				w.EmitEvents(
					TakeDamageEvent{
						Target: e[i],
						Amount: tot[i].Amount,
					},
				)

			case models.HealOverTimeEffect:
				w.EmitEvents(
					TakeHealEvent{
						Target: e[i],
						Amount: tot[i].Amount,
					},
				)
			}

			tot[i].RemainingTime -= dt // need to revise ? i think its wrong based on every tick

			if tot[i].RemainingTime <= 0 {
				w.PushCommands(
					zurvan.NewDeleteComponentsCommand(e[i],
						TakingOverTime{},
					),
				)
			}
		}
	})
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
	events := zurvan.OnEvent[TakeDamageEvent](w)

	for _, event := range events {
		h, ds := zurvan.QueryOne2[character.Health, character.DefensiveStats](w, event.Target)

		evasionPerc := 0.5 * ds.Evasion / (ds.Evasion + 100)
		if rand.Float64() > evasionPerc {
			continue
		}

		reduction := 0.5 * (ds.MagicResistance / (ds.MagicResistance + 200))
		finalDamage := event.Amount * (1 - reduction)

		h.Current -= finalDamage
		h.Current = min(h.Current, h.Max)

		w.EmitEvents(
			DeathEvent{
				Entity: event.Target,
			},
		)
	}
}

type TakeHealSystem struct{}

func (t *TakeHealSystem) Update(w *zurvan.World, dt time.Duration) {
	events := zurvan.OnEvent[TakeHealEvent](w)

	for _, event := range events {
		h := zurvan.QueryOne1[character.Health](w, event.Target)

		h.Current += event.Amount
		h.Current = min(h.Current, h.Max)
	}
}
