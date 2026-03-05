package combat

type PrimaryStats struct {
	Vitality     float64
	intelligence float64
	Willpower    float64
	Dexterity    float64
	Spirit       float64
}

type Health struct {
	Max          float64
	Current      float64
	Regeneration float64
}

type Mana struct {
	Max          float64
	Current      float64
	Regeneration float64
}

type DefensiveStats struct {
	MagicResistance float64
	Evasion         float64
}

type OffensiveStats struct {
	SpellPower     float64
	CriticalChance float64
	Accuracy       float64
}

type DamageTaken struct {
	Amount float64
}

type HealTaken struct {
	Amount float64
}
