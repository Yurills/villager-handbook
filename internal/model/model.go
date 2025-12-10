package model

import "fmt"

type Role int

const (
	Villager Role = iota
	Seer
	Werewolf
)

type Interaction struct {
	Actor  int
	Target int
	Type   string // e.g., "accuse", "claim", "fact"
	Result string // e.g., "Villager", "Seer", "Werewolf"
}

type PlayerInfo struct {
	VillagerCount int
	SearCount     int
	WarewolfCount int
	TotalPlayer   int
}

func (r Role) String() string {
	switch r {
	case Villager:
		return "Villager"
	case Seer:
		return "Seer"
	case Werewolf:
		return "Werewolf"
	default:
		return fmt.Sprintf("Role(%d)", r)
	}
}

type PlayerStat struct {
	ID                int
	RoleProbabilities map[Role]float64
}

func (ps PlayerStat) String() string {
	return fmt.Sprintf("Player %d: %v", ps.ID, ps.RoleProbabilities)
}

type LookaheadResult struct {
	ID      int
	Entropy float64
}
