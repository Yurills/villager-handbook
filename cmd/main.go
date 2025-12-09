package main

import (
	"fmt"

	"github.com/Yurills/villager-handbook/internal/engine"

	"github.com/Yurills/villager-handbook/internal/core"
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

		// This line will now trigger the Role.String() method for each role in the map.
		fmt.Printf("World %d: %+v\n", i, world.Roles)
	}

}
