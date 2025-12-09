package engine

import "github.com/Yurills/villager-handbook/internal/model"

const (
	WeightTruth         = 1.0
	WeightWolfLying     = 0.8 //werewolf lie often
	WeightSeerLying     = 0.6 //seer would lie as villager sometimes
	WeightVillagerLying = 0.1 //villager bait werewolf as seer (not oftern)
	WeightImpossible    = 0.0 //default for impossible cases
	NoWeight            = 1.0 //statement doesn't take into account.
)

// by world
func GetLikelihoodWeight(currentWorld []model.Role, interaction model.Interaction) float64 {
	if interaction.Type == "accuse" && interaction.Result == "Werewolf" {
		if currentWorld[interaction.Actor] == model.Werewolf {
			if currentWorld[interaction.Target] == model.Werewolf {
				return WeightImpossible //werewolf would not accuse fellow werewolf
			}
			if currentWorld[interaction.Target] == model.Villager || currentWorld[interaction.Target] == model.Seer {
				return WeightWolfLying //werewolf lies often
			}

		}
		if currentWorld[interaction.Actor] == model.Seer {
			if currentWorld[interaction.Target] == model.Werewolf {
				return WeightTruth //seer tells truth
			}
			if currentWorld[interaction.Target] == model.Villager || currentWorld[interaction.Target] == model.Seer {
				return WeightImpossible //seer would not accuse villager or seer
			}

		}
		if currentWorld[interaction.Actor] == model.Villager {
			if currentWorld[interaction.Target] == model.Werewolf {
				return NoWeight //villager might accuse werewolf or not
			}
			if currentWorld[interaction.Target] == model.Villager || currentWorld[interaction.Target] == model.Seer {
				return NoWeight //villager might accuse villager or seer or not
			}

		}
	}
	if interaction.Type == "accuse" && interaction.Result == "Seer" {
		if currentWorld[interaction.Actor] == model.Werewolf {
			if currentWorld[interaction.Target] == model.Seer || currentWorld[interaction.Target] == model.Villager {
				return NoWeight //means nothing, werewolf might accuse anyone
			}
			if currentWorld[interaction.Target] == model.Werewolf {
				return WeightWolfLying //werewolf protects fellow werewolf
			}

		}
		if currentWorld[interaction.Actor] == model.Seer {
			if currentWorld[interaction.Target] == model.Seer {
				return WeightTruth //seer tells truth
			}
			if currentWorld[interaction.Target] == model.Werewolf || currentWorld[interaction.Target] == model.Villager {
				return WeightSeerLying //seer would lie as villager sometimes
			}

		}
		if currentWorld[interaction.Actor] == model.Villager {
			if currentWorld[interaction.Target] == model.Seer {
				return NoWeight //villager tells truth
			}
			if currentWorld[interaction.Target] == model.Werewolf || currentWorld[interaction.Target] == model.Villager {
				return NoWeight //villager might accuse werewolf or not
			}

		}
	}
	if interaction.Type == "accuse" && interaction.Result == "Villager" {
		if currentWorld[interaction.Actor] == model.Werewolf {
			if currentWorld[interaction.Target] == model.Villager || currentWorld[interaction.Target] == model.Seer {
				return NoWeight //werewolf tells truth
			}
			if currentWorld[interaction.Target] == model.Werewolf {
				return WeightWolfLying //werewolf lies often
			}

		}
		if currentWorld[interaction.Actor] == model.Seer {
			if currentWorld[interaction.Target] == model.Villager {
				return WeightTruth //seer tells truth
			}
			if currentWorld[interaction.Target] == model.Werewolf || currentWorld[interaction.Target] == model.Seer {
				return WeightImpossible //seer would lie as villager sometimes
			}

			if currentWorld[interaction.Actor] == model.Villager {
				if currentWorld[interaction.Target] == model.Villager {
					return NoWeight //villager tells truth
				}
				if currentWorld[interaction.Target] == model.Werewolf || currentWorld[interaction.Target] == model.Seer {
					return NoWeight //villager might accuse werewolf or not
				}

			}
		}
	}

	if interaction.Type == "claim" && interaction.Result == "Werewolf" {
		return WeightImpossible //nobody would claim to be werewolf
	}
	if interaction.Type == "claim" && interaction.Result == "Seer" {
		if currentWorld[interaction.Actor] == model.Werewolf {
			return WeightWolfLying //werewolf lies often
		}
		if currentWorld[interaction.Actor] == model.Seer {
			return WeightTruth //seer tell truth
		}
		if currentWorld[interaction.Actor] == model.Villager {
			return WeightVillagerLying //villager bait werewolf as seer
		}
	}
	if interaction.Type == "claim" && interaction.Result == "Villager" {
		if currentWorld[interaction.Actor] == model.Werewolf {
			return WeightWolfLying //werewolf lies often
		}
		if currentWorld[interaction.Actor] == model.Seer {
			return WeightSeerLying ////seer would lie as villager sometimes
		}
		if currentWorld[interaction.Actor] == model.Villager {
			return NoWeight //villager tell truth
		}

	}

	if interaction.Type == "fact" && interaction.Result == "Werewolf" {
		if currentWorld[interaction.Target] == model.Werewolf {
			return WeightTruth
		}
		if currentWorld[interaction.Target] == model.Seer || currentWorld[interaction.Target] == model.Villager {
			return WeightImpossible
		}
	}
	if interaction.Type == "fact" && interaction.Result == "Seer" {
		if currentWorld[interaction.Target] == model.Seer {
			return WeightTruth
		}
		if currentWorld[interaction.Target] == model.Werewolf || currentWorld[interaction.Target] == model.Villager {
			return WeightImpossible
		}
	}
	if interaction.Type == "fact" && interaction.Result == "Villager" {
		if currentWorld[interaction.Target] == model.Villager {
			return WeightTruth
		}
		if currentWorld[interaction.Target] == model.Seer || currentWorld[interaction.Target] == model.Werewolf {
			return WeightImpossible
		}

	}
	return WeightImpossible //unknown interaction type
}
