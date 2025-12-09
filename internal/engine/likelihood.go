package engine

import "github.com/Yurills/villager-handbook/internal/model"

const (
	WeightTruth         = 1.0
	WeightWolfLying     = 0.8 //werewolf lie often
	WeightSeerLying     = 0.6 //seer would lie as villager sometimes
	WeightVillagerLying = 0.1 //villager bait werewolf as seer (not oftern)
	WeightImpossible    = 0.0 //default for impossible cases
)

// by world
func GetLikelihoodWeight(interaction model.Interaction) float64 {
	if interaction.Type == "accuse" && interaction.Result == "Werewolf" {
		if model.Role(interaction.Actor) == model.Werewolf {
			if model.Role(interaction.Target) == model.Werewolf {
				return WeightImpossible //werewolf would not accuse fellow werewolf
			}

		}
	}
	return WeightImpossible //unknown interaction type
}
