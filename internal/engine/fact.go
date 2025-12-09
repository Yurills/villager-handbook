package engine

import {
	"github.com/Yurills/villager-handbook/internal/model"
	"github.com/Yurills/villager-handbook/internal/engine"
}

//Update state space by fact
func UpdateState(actorRole model.Role, targetRole model.Role, interaction model.Interaction) float64 {