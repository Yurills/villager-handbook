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
			if model.Role(interaction.Target) == model.Villager || model.Role(interaction.Target) == model.Seer {
				return WeightWolfLying //werewolf lies often
			}

		}
		if model.Role(interaction.Actor) == model.Seer {
			if model.Role(interaction.Target) == model.Werewolf {
				return WeightTruth //seer tells truth
			}
			if model.Role(interaction.Target) == model.Villager || model.Role(interaction.Target) == model.Seer {
				return WeightSeerLying //seer would lie as villager sometimes
			}

		}
		if model.Role(interaction.Actor) == model.Villager {
			if model.Role(interaction.Target) == model.Werewolf {
				return WeightTruth //villager tells truth
			}
			if model.Role(interaction.Target) == model.Villager || model.Role(interaction.Target) == model.Seer {
				return WeightVillagerLying //villager lie rarely
			}

		}
	}
	if interaction.Type == "accuse" && interaction.Result == "Seer" {
		if model.Role(interaction.Actor) == model.Werewolf {
			if model.Role(interaction.Target) == model.Seer {
				return WeightTruth //werewolf tells truth
			}
			if model.Role(interaction.Target) == model.Werewolf || model.Role(interaction.Target) == model.Villager {
				return WeightWolfLying //werewolf lies often
			}

		}
		if model.Role(interaction.Actor) == model.Seer {
			if model.Role(interaction.Target) == model.Seer {
				return WeightTruth //seer tells truth
			}
			if model.Role(interaction.Target) == model.Werewolf || model.Role(interaction.Target) == model.Villager {
				return WeightSeerLying //seer would lie as villager sometimes
			}

		}
		if model.Role(interaction.Actor) == model.Villager {
			if model.Role(interaction.Target) == model.Seer {
				return WeightTruth //villager tells truth
			}
			if model.Role(interaction.Target) == model.Werewolf || model.Role(interaction.Target) == model.Villager {
				return WeightVillagerLying //villager lie rarely
			}

		}
	}
	if interaction.Type == "accuse" && interaction.Result == "Villager" {
		if model.Role(interaction.Actor) == model.Werewolf {
			if model.Role(interaction.Target) == model.Villager {
				return WeightTruth //werewolf tells truth
			}
			if model.Role(interaction.Target) == model.Werewolf || model.Role(interaction.Target) == model.Seer {
				return WeightWolfLying //werewolf lies often
			}

		}
		if model.Role(interaction.Actor) == model.Seer {
			if model.Role(interaction.Target) == model.Villager {
				return WeightTruth //seer tells truth
			}
			if model.Role(interaction.Target) == model.Werewolf || model.Role(interaction.Target) == model.Seer {
				return WeightSeerLying //seer would lie as villager sometimes
			}

		}
		if model.Role(interaction.Actor) == model.Villager {
			if model.Role(interaction.Target) == model.Villager {
				return WeightTruth //villager tells truth
			}
			if model.Role(interaction.Target) == model.Werewolf || model.Role(interaction.Target) == model.Seer {
				return WeightVillagerLying //villager lie rarely
			}

		}
	}

	if interaction.Type == "claim" && interaction.Result == "Werewolf" {
		return WeightImpossible //nobody would claim to be werewolf
	}
	if interaction.Type == "claim" && interaction.Result == "Seer" {
		if model.Role(interaction.Actor) == model.Werewolf {
			return WeightWolfLying //werewolf lies often
		}
		if model.Role(interaction.Actor) == model.Seer {
			return WeightTruth //seer tell truth
		}
		if model.Role(interaction.Actor) == model.Villager {
			return WeightVillagerLying //villager bait werewolf as seer
		}
	}
	if interaction.Type == "claim" && interaction.Result == "Villager" {
		if model.Role(interaction.Actor) == model.Werewolf {
			return WeightWolfLying //werewolf lies often
		}
		if model.Role(interaction.Actor) == model.Seer {
			return WeightSeerLying ////seer would lie as villager sometimes
		}
		if model.Role(interaction.Actor) == model.Villager {
			return WeightTruth //villager tell truth
		}

	}

	if interaction.Type == "fact" && interaction.Result == "Werewolf" {
		if model.Role(interaction.Target) == model.Werewolf {
			return WeightTruth
		}
		if model.Role(interaction.Target) == model.Seer || model.Role(interaction.Target) == model.Villager {
			return WeightImpossible
		}
	}
	if interaction.Type == "fact" && interaction.Result == "Seer" {
		if model.Role(interaction.Target) == model.Seer {
			return WeightTruth
		}
		if model.Role(interaction.Target) == model.Werewolf || model.Role(interaction.Target) == model.Villager {
			return WeightImpossible
		}
	}
	if interaction.Type == "fact" && interaction.Result == "Villager" {
		if model.Role(interaction.Target) == model.Villager {
			return WeightTruth
		}
		if model.Role(interaction.Target) == model.Seer || model.Role(interaction.Target) == model.Werewolf {
			return WeightImpossible
		}

	}
	return WeightImpossible //unknown interaction type
}
