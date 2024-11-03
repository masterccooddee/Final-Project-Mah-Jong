package main

import "fmt"

func main() {

	//startserver()
	// makeRoom(1, false)
	player := Player{}
	player.Ma.Card = []string{"w1", "w1", "w1", "w2", "w3", "w4", "w5", "w6", "w7", "w8", "w9", "w9", "w9"}
	player.Ma.splitCard()
	player.Pong = make(map[string]struct{})
	player.Pong["w1"] = struct{}{}
	player.Position = 3
	now = 2
	//comb := canChi(&player, "w6")
	fmt.Println(canChi(&player, "w2"))

}
