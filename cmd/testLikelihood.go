package main

import (
	"fmt"

	"github.com/Yurills/villager-handbook/internal/engine"

	"github.com/Yurills/villager-handbook/internal/core"
	"github.com/Yurills/villager-handbook/internal/model"
)

func main() {
	core.InputPlayer()
	playerInfo := core.GetPlayerInfo()
	players := []int{0, 1, 2, 3, 4}
	rules := engine.GameRule{
		NumVillagers:  playerInfo.VillagerCount,
		NumSeers:      playerInfo.SearCount,
		NumWerewolves: playerInfo.WarewolfCount,
	}
	e := engine.NewEngine(players, rules)

	for i, world := range e.Worlds {
		//Test fact
		interaction := model.Interaction{
			Actor:  0,
			Target: 2,
			Type:   "accuse",
			Result: "Werewolf",
		}

		likelihood := engine.GetLikelihoodWeight(world.Roles, interaction)
		fmt.Printf("World %d: %+v Likelihood: %f\n", i, world.Roles, likelihood)
	}
}
