package main

import (
	"log"
	"math/rand"
	"strconv"
)

type Player struct {
	player_in
	Ma       mao
	TingCard bool
	Pong     bool
	Eat      bool
	Gang     bool
}

type Room struct {
	Players []Player
	Cardset mao
	Room_ID int
	private bool
}

var roomlist = make(map[int]*Room)

func (r *Room) Addplayer(player player_in) {
	playerinfo := Player{player, mao{}, false, false, false, false}
	r.Players = append(r.Players, playerinfo)
}

func makeRoom(room_id int, private bool) {

	room := &(Room{Players: []Player{}, Room_ID: room_id, private: private})
	roomlist[room_id] = room

}

func (r *Room) go_in_room(player *player_in, room_id int) bool {

	_, exist := roomlist[room_id]
	if exist {
		if len(roomlist[room_id].Players) < 4 {

			player.Room_ID = room_id
			roomlist[room_id].Addplayer(*player)
			log.Println("Player " + player.ID + " join room " + strconv.Itoa(room_id))
			player.conn.Write([]byte("True Room " + strconv.Itoa(room_id)))
			return true
		} else {
			player.conn.Write([]byte("False Room is full"))
			log.Printf("Room %d is full, %s can't get in\n", room_id, player.ID)
			return false
		}
	} else {
		player.conn.Write([]byte("False Room not exist"))
		log.Printf("Room %d not exist, %s can't get in\n", room_id, player.ID)
		return false
	}

}

func Room_finder(player *player_in) int {
	for id, r := range roomlist {
		if len(r.Players) < 4 && r.private == false {
			r.go_in_room(player, id)
			return id
		}
	}

	room_id := rand.Intn(99) + 1
	_, exist := roomlist[room_id]
	for exist {
		room_id = rand.Intn(99) + 1
		_, exist = roomlist[room_id]
	}
	makeRoom(room_id, false)
	room := roomlist[room_id]
	room.go_in_room(player, room_id)
	return room_id
}

func (r *Room) leave_room(player *player_in) {
	for i, p := range r.Players {
		if p.ID == player.ID {
			r.Players = append(r.Players[:i], r.Players[i+1:]...)
			player.Room_ID = -1
			log.Println("Player " + player.ID + " leave room " + strconv.Itoa(r.Room_ID))
			player.conn.Write([]byte("True Leave room"))
			break
		}
	}
}
