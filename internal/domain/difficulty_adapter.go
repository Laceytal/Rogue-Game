package domain

import (
	"time"
)

type DifficultyAdapter struct {
	TotalDamageReceived  float64 `json:"total_damage_received"`
	TotalDamageDealт     float64 `json:"total_damage_dealt"`
	TotalHealthUsed      int     `json:"total_health_used"`
	TotalFightsWon       int     `json:"total_fights_won"`
	TotalFightsLost      int     `json:"total_fights_lost"`
	AverageHealthPercent float64 `json:"average_health_percent"`
	DeathsInLevel        int     `json:"deaths_in_level"`

	LastLevelStartTime  time.Time     `json:"-"`
	LevelCompletionTime time.Duration `json:"level_completion_time"`

	MonsterStatMultiplier  float64 `json:"monster_stat_multiplier"`
	MonsterCountMultiplier float64 `json:"monster_count_multiplier"`
	HealthItemBonus        float64 `json:"health_item_bonus"`
	ConsumableBonus        float64 `json:"consumable_bonus"`

	EasyThreshold float64 `json:"easy_threshold"`
	HardThreshold float64 `json:"hard_threshold"`

	RecentHealthLoss float64 `json:"recent_health_loss"`
	RecentMoves      int     `json:"recent_moves"`
}

func NewDifficultyAdapter() *DifficultyAdapter {
	return &DifficultyAdapter{
		MonsterStatMultiplier:  1.0,
		MonsterCountMultiplier: 1.0,
		HealthItemBonus:        1.0,
		ConsumableBonus:        1.0,
		EasyThreshold:          0.3,
		HardThreshold:          0.7,
		LastLevelStartTime:     time.Now(),
	}
}

func (da *DifficultyAdapter) analyzePerformance(player *Player) float64 {
	difficulty := 0.0
	factors := 0

	healthPercent := player.BaseStats.Health / float64(player.RegenLimit)
	if healthPercent > 0.8 {
		difficulty -= 0.2
	} else if healthPercent < 0.3 {
		difficulty += 0.3
	}
	factors++

	if da.TotalFightsWon+da.TotalFightsLost > 0 {
		winRate := float64(da.TotalFightsWon) / float64(da.TotalFightsWon+da.TotalFightsLost)
		if winRate > 0.85 {
			difficulty -= 0.25
		} else if winRate < 0.4 {
			difficulty += 0.25
		}
		factors++
	}

	if da.TotalDamageReceived > 0 {
		damageRatio := da.TotalDamageDealт / da.TotalDamageReceived
		if damageRatio > 3.0 {
			difficulty -= 0.15
		} else if damageRatio < 1.0 {
			difficulty += 0.2
		}
		factors++
	}

	if da.TotalHealthUsed > 10 {
		difficulty += 0.15
	} else if da.TotalHealthUsed < 3 {
		difficulty -= 0.1
	}
	factors++

	if da.DeathsInLevel > 0 {
		difficulty += float64(da.DeathsInLevel) * 0.2
		factors++
	}

	if factors > 0 {
		difficulty = difficulty / float64(factors)
	}

	if difficulty < -1.0 {
		difficulty = -1.0
	}
	if difficulty > 1.0 {
		difficulty = 1.0
	}

	return difficulty
}

func (da *DifficultyAdapter) AdjustDifficulty(player *Player) {
	performanceScore := da.analyzePerformance(player)

	adjustmentRate := 0.1

	if performanceScore < -da.EasyThreshold {
		da.MonsterStatMultiplier += adjustmentRate
		da.MonsterCountMultiplier += adjustmentRate * 0.5
		da.HealthItemBonus -= adjustmentRate * 0.5
		da.ConsumableBonus -= adjustmentRate * 0.3
	} else if performanceScore > da.HardThreshold {
		da.MonsterStatMultiplier -= adjustmentRate
		da.MonsterCountMultiplier -= adjustmentRate * 0.5
		da.HealthItemBonus += adjustmentRate * 0.7
		da.ConsumableBonus += adjustmentRate * 0.4
	}

	da.MonsterStatMultiplier = clamp(da.MonsterStatMultiplier, 0.5, 2.0)
	da.MonsterCountMultiplier = clamp(da.MonsterCountMultiplier, 0.3, 1.5)
	da.HealthItemBonus = clamp(da.HealthItemBonus, 1.0, 3.0)
	da.ConsumableBonus = clamp(da.ConsumableBonus, 0.7, 2.0)
}

func (da *DifficultyAdapter) ResetLevelStats() {
	da.LevelCompletionTime = time.Since(da.LastLevelStartTime)
	da.LastLevelStartTime = time.Now()
	da.DeathsInLevel = 0
	da.RecentHealthLoss = 0
	da.RecentMoves = 0
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
