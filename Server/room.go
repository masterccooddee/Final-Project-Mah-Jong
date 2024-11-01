package main

import (
	"log"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

type Player struct {
	player_in
	Ma       mao
	TingCard bool
	Pong     bool
	Eat      bool
	Gang     bool
	Position int
}

type Room struct {
	Players  []*Player
	Cardset  mao
	Room_ID  int
	private  bool
	running  bool
	recvchan chan string
}

var roomlist = make(map[int]*Room)

func (r *Room) Addplayer(player player_in) {
	playerinfo := Player{player_in: player, Ma: mao{}}
	r.Players = append(r.Players, &playerinfo)
}

func makeRoom(room_id int, private bool) {

	room := &(Room{Room_ID: room_id, private: private})
	roomlist[room_id] = room
	RROM := roomlist[room_id]
	RROM.recvchan = make(chan string, 2)
	go RROM.startgame()

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

	// 優先找人多的房間
	var rooms []*Room
	for _, r := range roomlist {
		if len(r.Players) < 4 && r.private == false {
			rooms = append(rooms, r)
		}
	}

	sort.Slice(rooms, func(i, j int) bool { return len(rooms[i].Players) > len(rooms[j].Players) })

	for _, r := range rooms {
		if len(r.Players) < 4 && r.private == false {
			r.go_in_room(player, r.Room_ID)
			return r.Room_ID
		}
	}

	// 沒有空房間，創建新房間
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

var cleaning bool

func RoomCleaner() {
	for {
		select {
		case <-time.After(50 * time.Second):
			cleaning = true
			log.Println("Start to clean room")
			for id, r := range roomlist {
				if len(r.Players) == 0 {
					delete(roomlist, id)
					log.Println("Room", id, "is empty, delete")
				}

			}
			log.Println("Room clean finish")
			cleaning = false
		}
	}
}
