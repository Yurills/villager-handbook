package model

type Role int

const (
	Villager Role = iota
	Seer
	Werewolf
)

type Interaction struct {
	Actor  int
	Target int
	Type   string // e.g., "accuse", "claim"
	Result string // e.g., "Villager", "Seer", "Werewolf"
}
