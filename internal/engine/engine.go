package engine

import (
	"slices"

	"github.com/Yurills/villager-handbook/internal/model"
)

type World struct {
	Roles  []model.Role
	Weight float64
}

type Engine struct {
	Players []int //state of players
	Worlds  []World
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
}
