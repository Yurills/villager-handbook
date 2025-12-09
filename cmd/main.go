package main

import (
	"fmt"

	"github.com/Yurills/villager-handbook/internal/engine"
	"github.com/Yurills/villager-handbook/internal/model"
)

<<<<<<< HEAD
func main2() {
	core.InputPlayer()
	playerInfo := core.GetPlayerInfo()
=======
func main() {
>>>>>>> 1248d97 (merged)
	players := []int{0, 1, 2, 3, 4}
	rules := engine.GameRule{
		NumVillagers:  3,
		NumSeers:      1,
		NumWerewolves: 1,
	}
	e := engine.NewEngine(players, rules)

	for i, world := range e.Worlds {

		// This line will now trigger the Role.String() method for each role in the map.
<<<<<<< HEAD
		fmt.Printf("World %d: %+v Weight: %f\n", i, world.Roles, world.Weight)
=======
		if (world.Weight) > 0 {
			fmt.Printf("World %d (Weight: %.4f): %+v\n", i, world.Weight, world.Roles)
		} else {
			fmt.Printf("World %d (Impossible): %+v\n", i, world.Roles)
		}
	}

	interaction := model.Interaction{
		Actor:  0,
		Target: 1,
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
>>>>>>> 1248d97 (merged)
	}

}
