package engine

import "github.com/Yurills/villager-handbook/internal/model"

//Update state space by fact
func UpdateState(actorRole model.Role, targetRole model.Role, interaction model.Interaction) float64 {