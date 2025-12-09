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
		if actorRole == model.Seer {
			if targetRole == model.Werewolf {
				return WeightTruth //seer tells truth
			}
			if targetRole == model.Villager || targetRole == model.Seer {
				return WeightSeerLying //seer would lie as villager sometimes
			}

		}
		if actorRole == model.Villager {
			if targetRole == model.Werewolf {
				return WeightTruth //villager tells truth
			}
			if targetRole == model.Villager || targetRole == model.Seer {
				return WeightVillagerLying //villager lie rarely
			}

		}
	}
	if interaction.Type == "accuse" && interaction.Result == "Seer" {
		if actorRole == model.Werewolf {
			if targetRole == model.Seer {
				return WeightTruth //werewolf tells truth
			}
			if targetRole == model.Werewolf || targetRole == model.Villager {
				return WeightWolfLying //werewolf lies often
			}

		}
		if actorRole == model.Seer {
			if targetRole == model.Seer {
				return WeightTruth //seer tells truth
			}
			if targetRole == model.Werewolf || targetRole == model.Villager {
				return WeightSeerLying //seer would lie as villager sometimes
			}

		}
		if actorRole == model.Villager {
			if targetRole == model.Seer {
				return WeightTruth //villager tells truth
			}
			if targetRole == model.Werewolf || targetRole == model.Villager {
				return WeightVillagerLying //villager lie rarely
			}

		}
	}
	if interaction.Type == "accuse" && interaction.Result == "Villager" {
		if actorRole == model.Werewolf {
			if targetRole == model.Villager {
				return WeightTruth //werewolf tells truth
			}
			if targetRole == model.Werewolf || targetRole == model.Seer {
				return WeightWolfLying //werewolf lies often
			}

		}
		if actorRole == model.Seer {
			if targetRole == model.Villager {
				return WeightTruth //seer tells truth
			}
			if targetRole == model.Werewolf || targetRole == model.Seer {
				return WeightSeerLying //seer would lie as villager sometimes
			}

		}
		if actorRole == model.Villager {
			if targetRole == model.Villager {
				return WeightTruth //villager tells truth
			}
			if targetRole == model.Werewolf || targetRole == model.Seer {
				return WeightVillagerLying //villager lie rarely
			}

		}
	}

	if interaction.Type == "claim" && interaction.Result == "Werewolf" {
		return WeightImpossible //nobody would claim to be werewolf
	}
	if interaction.Type == "claim" && interaction.Result == "Seer" {
		if actorRole == model.Werewolf {
			return WeightWolfLying //werewolf lies often
		}
		if actorRole == model.Seer {
			return WeightTruth //seer tell truth
		}
		if actorRole == model.Villager {
			return WeightVillagerLying //villager bait werewolf as seer
		}
	}
	if interaction.Type == "claim" && interaction.Result == "Villager" {
		if actorRole == model.Werewolf {
			return WeightWolfLying //werewolf lies often
		}
		if actorRole == model.Seer {
			return WeightSeerLying ////seer would lie as villager sometimes
		}
		if actorRole == model.Villager {
			return WeightTruth //villager tell truth
		}
	}
	return WeightImpossible //unknown interaction type
}
