package engine

import "github.com/Yurills/villager-handbook/internal/model"

const (
	WeightTruth         = 1.0
	WeightWolfLying     = 0.8 //werewolf lie often
	WeightWolfBussing   = 0.2 //werewolf accuse fellow werewolf sometimes (distancing strategy)
	WeightSeerLying     = 0.6 //seer would lie as villager sometimes
	WeightVillagerLying = 0.1 //villager bait werewolf as seer (not oftern)
	WeightImpossible    = 0.0 //default for impossible cases
	NoWeight            = 1.0 //statement doesn't take into account (e.g. villager accuse werewolf).
)

// by world
func GetLikelihoodWeight(currentWorld []model.Role, interaction model.Interaction) float64 {
	//1. Handle Hard Facts (Game Engine Events)
	//Logic: if the world doesn't match the hard fact, it's impossible.
	if interaction.Type == "fact" {
		actualRole := currentWorld[interaction.Target]
		// Convert string Result to Enum (Assuming interaction.Result matches model.Role)
		// You might need a helper here to convert string "Werewolf" -> model.Werewolf
		if roleMatches(actualRole, interaction.Result) {
			return WeightTruth
		}
		return WeightImpossible
	}

	actorRole := currentWorld[interaction.Actor]
	targetRole := currentWorld[interaction.Target]

	// 2. Handle Accusations ("I think Target is X")
	if interaction.Type == "accuse" {
		// CONTEXT: Actor says "Target is [Result]"
		accusedRole := interaction.Result // e.g. "Werewolf"

		// CASE A: Accusing someone of being a WEREWOLF
		if accusedRole == "Werewolf" {
			switch actorRole {
			case model.Werewolf:
				if targetRole == model.Werewolf {
					return WeightWolfBussing //wolves do accuse each other
				}
				//wolf accusing Non-wolf
				return WeightWolfLying //wolves lie often

			case model.Seer:
				if targetRole == model.Werewolf {
					return WeightTruth //seer found a wolf
				}
				return WeightImpossible //seer wouldn't lie and call a Good person Evil

			case model.Villager:
				return NoWeight //villagers guess randomly
			}
		}
	}

	//3. Handle Role Claim ("I am X")
	if interaction.Type == "claim" {
		claimedRole := interaction.Result // e.g. "Seer"

		// CASE A: Claiming Werewolf (Nobody does this usually)
		if claimedRole == "Werewolf" {
			return WeightImpossible
		}

		// CASE B: Claiming Seer
		if claimedRole == "Seer" {
			switch actorRole {
			case model.Werewolf:
				return WeightWolfLying //wolf lying as seer
			case model.Seer:
				return WeightTruth //seer telling truth
			case model.Villager:
				return WeightVillagerLying //villager bluffing as seer
			}
		}
		// CASE C: Claiming Villager
		if claimedRole == "Villager" {
			switch actorRole {
			case model.Werewolf:
				return WeightWolfLying //wolf lying as villager
			case model.Seer:
				return WeightSeerLying //seer hiding
			case model.Villager:
				return WeightTruth //villager telling truth
			}
		}

	}
	return WeightImpossible //unknown
}

func roleMatches(r model.Role, roleStr string) bool {
	switch roleStr {
	case "Villager":
		return r == model.Villager
	case "Seer":
		return r == model.Seer
	case "Werewolf":
		return r == model.Werewolf
	default:
		return false
	}
}
