package core

import (
	"github.com/Yurills/villager-handbook/internal/model"
	"fmt"
)

var player model.PlayerInfo

func InputPlayer(){
	fmt.Println("Input number of Villager: ")
	fmt.Scan(&player.VillagerCount)
	fmt.Println("Input number of Sear: ")
	fmt.Scan(&player.SearCount)
	fmt.Println("Input number of Warewolf: ")
	fmt.Scan(&player.WarewolfCount)
	player.TotalPlayer = player.VillagerCount + player.SearCount + player.WarewolfCount
	fmt.Println("Total Players:", player.TotalPlayer)
}

func GetPlayerInfo() model.PlayerInfo{
	return player
}

func UpdateRound(){
	
}

