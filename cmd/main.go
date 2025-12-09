package main

import (
	"fmt"

	"github.com/Yurills/villager-handbook/internal/engine"
	"github.com/Yurills/villager-handbook/internal/model"
)

func main() {
	// core.InputPlayer()
	// playerInfo := core.GetPlayerInfo()
	players := []int{0, 1, 2, 3, 4, 5}
	rules := engine.GameRule{
		NumVillagers:  3,
		NumSeers:      1,
		NumWerewolves: 2,
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
	interaction := model.Interaction{
		Actor:  3,
		Target: 3,
		Type:   "fact",
		Result: "Villager",
	}
	e.ProcessMove(interaction)

	interaction = model.Interaction{
		Actor:  0,
		Target: 1,
		Type:   "accuse",
		Result: "Werewolf",
	}

	e.ProcessMove(interaction)

	interaction = model.Interaction{
		Actor:  1,
		Target: 0,
		Type:   "accuse",
		Result: "Werewolf",
	}
	e.ProcessMove(interaction)

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

	interaction = model.Interaction{
		Actor:  0,
		Target: 0,
		Type:   "fact",
		Result: "Werewolf",
	}

	e.ProcessMove(interaction)

	interaction = model.Interaction{
		Actor:  2,
		Target: 1,
		Type:   "claim",
		Result: "Seer",
	}
	e.ProcessMove(interaction)

	interaction = model.Interaction{
		Actor:  5,
		Target: 2,
		Type:   "claim",
		Result: "Seer",
	}
	e.ProcessMove(interaction)

	for i, world := range e.Worlds {
		if (world.Weight) > 0 {
			fmt.Printf("After Move - World %d (Weight: %.4f): %+v\n", i, world.Weight, world.Roles)
		} else {
			fmt.Printf("After Move - World %d (Impossible): %+v\n", i, world.Roles)
		}
	}

	stats = e.GetStats()
	for _, stat := range stats {
		fmt.Println(stat)
	}

}
