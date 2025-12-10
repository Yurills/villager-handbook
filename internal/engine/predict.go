package engine

import (
	"fmt"

	"github.com/Yurills/villager-handbook/internal/model"
)

var tempWorlds []World

// deep copy a single world
func cloneWorld(w World) World {
	newRoles := make([]model.Role, len(w.Roles))
	copy(newRoles, w.Roles)

	return World{
		Roles:  newRoles,
		Weight: w.Weight,
	}
}

var allRoles = []model.Role{
	model.Seer,
	model.Werewolf,
	model.Villager,
}

// predict by branching each world into multiple worlds based on possible facts
func (e *Engine) PredictMove() {
	var newWorlds []World

	for _, world := range e.Worlds {

		for _, target := range e.Players {

			for _, r := range allRoles {

				move := model.Interaction{
					Actor:  -1,
					Target: target,
					Type:   "fact",
					Result: r.String(), // match your model.Role type
				}

				// Clone world (deep copy)
				cloned := cloneWorld(world)

				// Apply likelihood
				likelihood := GetLikelihoodWeight(cloned.Roles, move)
				cloned.Weight *= likelihood

				// Store branched world
				newWorlds = append(newWorlds, cloned)
			}
		}
	}

	// normalize all new branching worlds
	totalWeight := 0.0
	for i := range newWorlds {
		totalWeight += newWorlds[i].Weight
	}

	if totalWeight > 0 {
		for i := range newWorlds {
			newWorlds[i].Weight /= totalWeight
		}
	} else {
		uniform := 1.0 / float64(len(newWorlds))
		for i := range newWorlds {
			newWorlds[i].Weight = uniform
		}
	}

	tempWorlds = newWorlds
}

func (e *Engine) GetPredictStat() []model.PlayerStat {

	//Outer Key: Player ID (int)
	//Inner Key: Role
	//Value: Accumulated weight
	playerRoleTotals := make(map[int]map[model.Role]float64)

	for _, id := range e.Players {
		playerRoleTotals[id] = make(map[model.Role]float64)
	}

	for _, world := range tempWorlds {
		if world.Weight <= 0 {
			continue
		}

		for playerID, role := range world.Roles {
			playerRoleTotals[playerID][role] += world.Weight
		}
	}

	var playerStats []model.PlayerStat

	for _, id := range e.Players {
		stats := model.PlayerStat{
			ID:                id,
			RoleProbabilities: playerRoleTotals[id],
		}
		playerStats = append(playerStats, stats)
	}
	return playerStats
}

// provides a recommendation based on werewolf probabilities
func (e *Engine) GetRecommend(stats []model.PlayerStat) string {

    if len(stats) == 0 {
        return "No statistics available."
    }

    type scored struct {
        ID     int
        WWProb float64
    }

    var list []scored

    // Build list
    for _, ps := range stats {
        list = append(list, scored{
            ID:     ps.ID,
            WWProb: ps.RoleProbabilities[model.Werewolf],
        })
    }

    // Find highest werewolf probability
    highest := list[0].WWProb
    for _, s := range list {
        if s.WWProb > highest {
            highest = s.WWProb
        }
    }

    // Build recommendation output
    result := "Recommend voting these players:\n"
    for _, s := range list {
        if s.WWProb == highest {
            result += fmt.Sprintf(
                "- Player %d (Werewolf likelihood: %.2f %%)\n",
                s.ID,
                s.WWProb*100,
            )
        }
    }

    return result
}
 
