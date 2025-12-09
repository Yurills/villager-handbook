package main

import (
	"fmt"

	"github.com/Yurills/villager-handbook/internal/engine"
	"github.com/Yurills/villager-handbook/internal/model"
	"github.com/Yurills/villager-handbook/internal/core"
	// tea "github.com/charmbracelet/bubbletea"
)

func main2() {
	core.InputPlayer()
	playerInfo := core.GetPlayerInfo()
	
	players := make([]int, playerInfo.TotalPlayer)
	for i := 0; i < playerInfo.TotalPlayer; i++ {
		players[i] = i
	}


	// players := []int{0, 1, 2, 3, 4, 5}
	rules := engine.GameRule{
		NumVillagers:  playerInfo.VillagerCount,
		NumSeers:      playerInfo.SearCount,
		NumWerewolves: playerInfo.TotalPlayer,
	}
	e := engine.NewEngine(players, rules)

	// for i, world := range e.Worlds {

	// 	// This line will now trigger the Role.String() method for each role in the map.
	// 	if (world.Weight) > 0 {
	// 		fmt.Printf("World %d (Weight: %.4f): %+v\n", i, world.Weight, world.Roles)
	// 	} else {
	// 		fmt.Printf("World %d (Impossible): %+v\n", i, world.Roles)
	// 	}
	// }

	// this wait for input to proceed this is for AddEvent
	var interaction model.Interaction
	fmt.Print("Add actor: ")
	fmt.Scan(interaction.Actor)
	fmt.Print("Add target: ")
	fmt.Scan(interaction.Target)
	fmt.Print("Add type: ")
	fmt.Scan(interaction.Type)
	fmt.Print("Add result: ")
	fmt.Scan(interaction.Result)

	e.ProcessMove(interaction)

	// interaction = model.Interaction{
	// 	Actor:  0,
	// 	Target: 1,
	// 	Type:   "accuse",
	// 	Result: "Werewolf",
	// }

	// e.ProcessMove(interaction)

	// interaction = model.Interaction{
	// 	Actor:  1,
	// 	Target: 0,
	// 	Type:   "accuse",
	// 	Result: "Werewolf",
	// }
	// e.ProcessMove(interaction)
	// This is show when click Show Recommend Move
	for i, world := range e.Worlds {
		if (world.Weight) > 0 {
			fmt.Printf("After Move - World %d (Weight: %.4f): %+v\n", i, world.Weight, world.Roles)
		} else {
			fmt.Printf("After Move - World %d (Impossible): %+v\n", i, world.Roles)
		}
	}
	stats := e.GetStats()
	for _, stat := range stats {
		fmt.Println(stat)
	}

	// stats := e.GetStats()
	// for _, stat := range stats {
	// 	fmt.Println(stat)
	// }

	// interaction = model.Interaction{
	// 	Actor:  0,
	// 	Target: 0,
	// 	Type:   "fact",
	// 	Result: "Werewolf",
	// }

	// e.ProcessMove(interaction)

	// interaction = model.Interaction{
	// 	Actor:  2,
	// 	Target: 1,
	// 	Type:   "claim",
	// 	Result: "Seer",
	// }
	// e.ProcessMove(interaction)

	// interaction = model.Interaction{
	// 	Actor:  5,
	// 	Target: 2,
	// 	Type:   "claim",
	// 	Result: "Seer",
	// }
	// e.ProcessMove(interaction)

	// for i, world := range e.Worlds {
	// 	if (world.Weight) > 0 {
	// 		fmt.Printf("After Move - World %d (Weight: %.4f): %+v\n", i, world.Weight, world.Roles)
	// 	} else {
	// 		fmt.Printf("After Move - World %d (Impossible): %+v\n", i, world.Roles)
	// 	}
	// }

	

}
