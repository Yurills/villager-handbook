package engine

import (
	"math"

	"github.com/Yurills/villager-handbook/internal/model"
)

func (e *Engine) calculateEntropy(playerStat []model.PlayerStat) float64 {
	totalEntropy := 0.0

	for _, ps := range playerStat {
		playerEntropy := 0.0
		for _, prob := range ps.RoleProbabilities {
			if prob > 0 {
				playerEntropy += -prob * math.Log2(prob)
			}
		}
		totalEntropy += playerEntropy
	}
	return totalEntropy
}
