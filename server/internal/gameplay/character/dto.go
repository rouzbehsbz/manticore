package character

import (
	"github.com/rouzbehsbz/manticore/server/pkg/util"
	"github.com/rouzbehsbz/zurvan"
)

func Characters(w *zurvan.World) (*util.SyncMap[uint32, zurvan.Entity], bool) {
	return zurvan.Resource[*util.SyncMap[uint32, zurvan.Entity]](w)
}

type Level struct {
	Value             int
	Xp                int
	NextLevelXpNeeded int
}

type PrimaryStats struct {
	Vitality     float64
	Intelligence float64
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
	CriticalRating float64
	Accuracy       float64
}

type JoinsWorldEvent struct {
	Character zurvan.Entity
	Id        uint32
}

type LevelUpEvent struct {
	Entity zurvan.Entity
}

type XpGainedEvent struct {
	Entity zurvan.Entity
	Amount int
}

type RecalculateStatsEvent struct {
	Entity zurvan.Entity
}
