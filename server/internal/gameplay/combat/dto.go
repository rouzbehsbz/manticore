package combat

import (
	"github.com/rouzbehsbz/zurvan"
)

type RecalculateStatsEvent struct {
	Entity zurvan.Entity
}

type DamageTakenEvent struct {
	Entity zurvan.Entity
	Amount float64
}

type HealTakenEvent struct {
	Entity zurvan.Entity
	Amount float64
}
