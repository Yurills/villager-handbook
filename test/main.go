package main

import (
	"github.com/Yurills/villager-handbook/internal/engine"
	"github.com/Yurills/villager-handbook/internal/model"
)

func main() {
	rules := engine.GameRule{
		NumVillagers:  3,
		NumSeers:      1,
		NumWerewolves: 2,
	}
	players := []int{0, 1, 2, 3, 4, 5}
	e := engine.NewEngine(players, rules)

	e.ProcessMove(model.Interaction{
		Actor:  0,
		Target: 0,
		Type:   "claim",
		Result: "seer",
	})

	e.ProcessMove(model.Interaction{
		Actor:  0,
		Target: 2,
		Type:   "accuse",
		Result: "Villager",
	})
	e.ProcessMove(model.Interaction{
		Actor:  1,
		Target: 1,
		Type:   "claim",
		Result: "seer",
	})
	e.ProcessMove(model.Interaction{
		Actor:  1,
		Target: 3,
		Type:   "accuse",
		Result: "Villager",
	})

	prob := e.GetStats()
	results := e.LookaheadBestCandidate()

	for _, res := range prob {
		println(res.String())
	}

	for _, res := range results {
		println("Player", res.ID, "Expected Entropy:", res.Entropy)
	}

}
