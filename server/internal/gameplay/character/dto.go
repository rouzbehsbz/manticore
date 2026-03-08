package character

import (
	"math"

	"github.com/rouzbehsbz/manticore/server/internal/infra/db/sources"
	"github.com/rouzbehsbz/manticore/server/pkg/util"
	"github.com/rouzbehsbz/zurvan"
)

func CharactersMap(w *zurvan.World) (*util.SyncMap[uint32, sources.Character], bool) {
	return zurvan.Resource[*util.SyncMap[uint32, sources.Character]](w)
}

func CharacterEntityMap(w *zurvan.World) (*util.SyncMap[uint32, zurvan.Entity], bool) {
	return zurvan.Resource[*util.SyncMap[uint32, zurvan.Entity]](w)
}

type Level struct {
	Value             int
	Xp                int
	NextLevelXpNeeded int
}

func NewLevel(value, xp int) Level {
	return Level{
		Value:             value,
		Xp:                xp,
		NextLevelXpNeeded: int(float64(100) * math.Pow(float64(value), 1.5)),
	}
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
