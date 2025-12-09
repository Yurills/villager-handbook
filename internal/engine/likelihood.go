package engine

import "github.com/Yurills/villager-handbook/internal/model"

const (
	WeightTruth         = 1.0
	WeightWolfLying     = 0.8 //werewolf lie often
	WeightSeerLying     = 0.6 //seer would lie as villager sometimes
	WeightVillagerLying = 0.1 //villager bait werewolf as seer (not oftern)
	WeightImpossible    = 0.0 //default for impossible cases
)


//by world
func GetLikelihoodWeight(actorRole model.Role, targetRole model.Role, interaction model.Interaction) float64 {
	if interaction.Type == "accuse" && interaction.Result == "Werewolf" {
		if actorRole == model.Werewolf {
			if targetRole == model.Werewolf {
				return WeightImpossible //werewolf would not accuse fellow werewolf
			}
			if targetRole == model.Villager || targetRole == model.Seer {
				return WeightWolfLying //werewolf lies often
			}


	}
	return WeightImpossible //unknown interaction type
}
