package main

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/go-zeromq/zmq4"
)

var router zmq4.Socket

type Position struct {
	Pos map[string]int
}

var zmqmu sync.Mutex

func sendtoplayer(msg string, ID string) {
	msgout := zmq4.NewMsgFrom([]byte(ID), []byte(msg))
	zmqmu.Lock()
	router.SendMulti(msgout)
	zmqmu.Unlock()
}

func (r *Room) sendtoall(msg string) {
	for _, p := range r.Players {
		sendtoplayer(msg, p.ID)
	}
}

func (r *Room) startgame() {

	//確認是否有4個玩家
	for len(r.Players) != 4 {
		if _, exist := roomlist[r.Room_ID]; exist == false {
			return
		}
		log.Println("Room", r.Room_ID, "is not full", len(r.Players))
		time.Sleep(1 * time.Second)
	}

	//通知所有玩家遊戲開始
	r.running = true
	r.sendtoall("Game start")

	//隨機選座位
	position := make(map[string]int)
	rand.Shuffle(len(r.Players), func(i, j int) { r.Players[i], r.Players[j] = r.Players[j], r.Players[i] })
	for i, p := range r.Players {
		p.Position = i
		position[p.ID] = p.Position
	}
	//通知所有玩家座位
	cli_pos, _ := json.Marshal(Position{Pos: position})
	r.sendtoall(string(cli_pos))

	for {
		command := <-r.recvchan
		log.Println(command)
	}
}
