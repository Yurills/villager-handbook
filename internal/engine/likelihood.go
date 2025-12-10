package engine

import (
	"strings"

	"github.com/Yurills/villager-handbook/internal/model"
)

const (
	// PROBABILITIES (P(Action | Role))
	// We try to keep these below 1.0 to respect probability theory,
	// but 1.0 is fine for "Baseline/Truth".

	WeightTruth    = 1.0
	WeightStandard = 1.0 // Baseline behavior (random guessing)

	// Wolf Strategies
	WeightWolfLying     = 0.8 // Standard Wolf play
	WeightWolfBussing   = 0.2 // Rare: Wolf accusing partner
	WeightWolfPocketing = 0.6 // Wolf defending a Villager to buy trust

	// Seer Strategies
	WeightSeerTruth   = 2.0 // High confidence: Seer accuses Wolf , x2.0 is rewarded (maybe adjust?)
	WeightSeerMistake = 0.1 // Seer accusing Good (Possible if guessing, but rare)
	WeightSeerHiding  = 0.5 // Seer claiming Villager

	// Logic Breakers
	WeightImpossible = 0.0
)

func GetLikelihoodWeight(currentWorld []model.Role, interaction model.Interaction) float64 {
	// 1. NORMALIZE INPUT
	result := strings.ToLower(interaction.Result)

	// 2. HARD FACTS (Game Engine / Death Reveals)
	if interaction.Type == "fact" {
		actualRole := currentWorld[interaction.Target]
		if roleMatches(actualRole, result) {
			return WeightTruth
		}
		return WeightImpossible
	}

	// Get roles for the behavioral checks
	actorRole := currentWorld[interaction.Actor]
	targetRole := currentWorld[interaction.Target]

	// 3. ACCUSATIONS ("I think Target is X")
	if interaction.Type == "accuse" {

		// --- CASE A: Accusing someone of being EVIL/WEREWOLF ---
		if result == "werewolf" || result == "evil" {
			switch actorRole {
			case model.Werewolf:
				if targetRole == model.Werewolf {
					return WeightWolfBussing // Wolf accusing Wolf (Low chance)
				}
				return WeightWolfLying // Wolf accusing Good (Standard play)

			case model.Seer:
				if targetRole == model.Werewolf {
					return WeightSeerTruth // Seer found a Wolf! (High weight)
				}
				// CRITICAL FIX: Do NOT use WeightImpossible here.
				// Seers can guess wrong before they scan.
				return WeightSeerMistake

			case model.Villager:
				return WeightStandard // Villagers guess randomly
			}
		}

		// --- CASE B: Accusing someone of being GOOD/VILLAGER ---
		if result == "villager" || result == "seer" || result == "good" {
			switch actorRole {
			case model.Werewolf:
				if targetRole == model.Werewolf {
					return WeightWolfLying // Wolf defending Wolf (Lying by omission)
				}
				// Wolf defending Good Player ("Pocketing")
				return WeightWolfPocketing

			case model.Seer:
				if targetRole != model.Werewolf {
					// Seer defending Good (likely scanned them)
					return WeightSeerTruth
				}
				return WeightSeerMistake // Seer defending Wolf (Very rare mistake)

			case model.Villager:
				return WeightStandard
			}
		}
	}

	// 4. ROLE CLAIMS ("I am X")
	if interaction.Type == "claim" {
		// Note: In "claim", the Target is usually the Actor themselves.

		// CASE A: Claiming Werewolf (Gamethrowing)
		if result == "werewolf" {
			return WeightImpossible
		}

		// CASE B: Claiming Seer
		if result == "seer" {
			switch actorRole {
			case model.Werewolf:
				return WeightWolfLying // Wolf fake claiming
			case model.Seer:
				return WeightTruth // Real Seer
			case model.Villager:
				return 0.1 // Villager fake claiming (Rare "Bait")
			}
		}

		// CASE C: Claiming Villager
		if result == "villager" {
			switch actorRole {
			case model.Werewolf:
				return WeightWolfLying // Wolf hiding
			case model.Seer:
				return WeightSeerHiding // Seer hiding
			case model.Villager:
				return WeightTruth // Villager being honest
			}
		}
	}

	// Default safety
	return WeightStandard
}

// Helper to match strings to Enum
func roleMatches(r model.Role, roleStr string) bool {
	switch roleStr {
	case "villager":
		return r == model.Villager
	case "good":
		return r == model.Villager || r == model.Seer // "Good" covers both
	case "seer":
		return r == model.Seer
	case "werewolf", "evil":
		return r == model.Werewolf
	default:
		return false
	}
}
