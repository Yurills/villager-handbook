package engine

import (
	"slices"
	"sync"

	"github.com/Yurills/villager-handbook/internal/model"
)

type World struct {
	Roles  []model.Role
	Weight float64
}

type Engine struct {
	Players []int //state of players
	Worlds  []World
	mu      sync.RWMutex
}

type GameRule struct {
	NumVillagers  int
	NumSeers      int
	NumWerewolves int
}

func NewEngine(players []int, rules GameRule) *Engine {
	e := &Engine{
		Players: players,
	}
	e.generateWorlds(rules)
	return e
}

func (e *Engine) ProcessMove(move model.Interaction) {
	totalWeight := 0.0

	if move.Type == "fact" && move.Result != "good" && move.Result != "evil" {
		e.SetPlayerDead(move.Actor)
	}

	for i := range e.Worlds {
		likelihood := GetLikelihoodWeight(e.Worlds[i].Roles, move)
		e.Worlds[i].Weight *= likelihood
		totalWeight += e.Worlds[i].Weight
	}

	//normalize weights
	if totalWeight > 0 {
		for i := range e.Worlds {
			e.Worlds[i].Weight /= totalWeight
		}
	} else {
		//all worlds are impossible, reset weights uniformly
		uniformWeight := 1.0 / float64(len(e.Worlds))
		for i := range e.Worlds {
			e.Worlds[i].Weight = uniformWeight
		}
	}
}

func (e *Engine) GetStats() []model.PlayerStat {

	//Outer Key: Player ID (int)
	//Inner Key: Role
	//Value: Accumulated weight

	playerStats := []model.PlayerStat{}

	for i := range e.Players {
		stat := e.GetPlayerRoleProbabilities(i)
		playerStats = append(playerStats, stat)
	}

	return playerStats

}

func (e *Engine) GetPlayerRoleProbabilities(playerID int) model.PlayerStat {
	roleProbabilities := make(map[model.Role]float64)

	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, world := range e.Worlds {
		if world.Weight <= 0 {
			continue
		}

		role := world.Roles[playerID]
		roleProbabilities[role] += world.Weight
	}

	return model.PlayerStat{
		ID:                playerID,
		RoleProbabilities: roleProbabilities,
	}
}

func (e *Engine) generateWorlds(rules GameRule) {
	// create deck consisting of all roles
	deck := []model.Role{}
	for i := 0; i < rules.NumVillagers; i++ {
		deck = append(deck, model.Villager)
	}
	for i := 0; i < rules.NumSeers; i++ {
		deck = append(deck, model.Seer)
	}
	for i := 0; i < rules.NumWerewolves; i++ {
		deck = append(deck, model.Werewolf)
	}

	slices.Sort(deck)

	if len(deck) != len(e.Players) {
		panic("number of roles does not match number of players")
	}

	used := make([]bool, len(deck))
	currentAssignment := make([]model.Role, 0, len(deck))

	e.Worlds = []World{}

	var backtrack func()
	backtrack = func() {
		// return clause
		if len(currentAssignment) == len(deck) {
			// append current assignment as a new world
			copyOfCurrent := make([]model.Role, len(currentAssignment))
			copy(copyOfCurrent, currentAssignment)
			e.Worlds = append(e.Worlds, World{ //append world result
				Roles:  copyOfCurrent,
				Weight: 1.0, //initial weight
			})
			return
		}

		for i := 0; i < len(deck); i++ {
			if used[i] {
				continue
			}
			//no duplicate world
			if i > 0 && deck[i] == deck[i-1] && !used[i-1] {
				continue
			}

			used[i] = true
			currentAssignment = append(currentAssignment, deck[i])

			backtrack()

			used[i] = false
			currentAssignment = currentAssignment[:len(currentAssignment)-1]

		}

	}

	backtrack()
	
	for i := range e.Worlds {
		e.Worlds[i].Weight /= float64(len(e.Worlds))
	}
	
}

func (e *Engine) fork() *Engine { //use for cloning worlds to do expectimax
	e.mu.RLock()
	defer e.mu.RUnlock()

	newWorlds := make([]World, len(e.Worlds))
	copy(newWorlds, e.Worlds)

	return &Engine{
		Players: e.Players,
		Worlds:  newWorlds,
	}
}

func (e *Engine) SetPlayerDead(playerID int) {
	e.Players[playerID] = -1
}
