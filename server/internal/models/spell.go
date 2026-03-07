package models

import "time"

type Rarity int

const (
	CommonRarity Rarity = iota
	RareRarity
	EpicRarity
	LegendaryRarity
)

type EffectType int

const (
	DamageEffect EffectType = iota
	HealEffect
	DamageOverTimeEffect
	HealOverTimeEffect
)

type SpellEffect struct {
	Type        EffectType
	Amount      float64
	Coefficient float64
	Duration    time.Duration
}

type Spell struct {
	Id       int
	Rarity   Rarity
	ManaCost float64
	CastTime time.Duration
	Cooldown time.Duration
	Effect   SpellEffect
}
