package main

import (
	"context"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

type Player struct {
	*player_in
	Ma       mao
	TingCard bool
	Pong     map[string]struct{}
	Chi      map[string]struct{}
	Gang     map[string]struct{}

	clean    bool //門清
	Position int  //0.東 1.南 2.西 3.北
	Point    int
}

type Room struct {
	Players         []*Player
	Cardset         mao
	Room_ID         int
	private         bool
	running         bool        //遊戲是否開始
	recvchan        chan string //接收玩家訊息
	lastcard        int         //剩幾張
	round           int         //第幾局 如：東1
	wind            int         // 0.東 1.南 2.西 3.北
	bunround        int         //本場
	gang            [4]int      //槓的次數
	now             int         //當前玩家
	real_player_nun int         //真實玩家數
}

var roomlist = make(map[int]*Room)

func (r *Room) Addplayer(player *player_in) {
	playerinfo := Player{player_in: player, Ma: mao{}, Point: 25000, clean: true, Pong: make(map[string]struct{}), Chi: make(map[string]struct{}), Gang: make(map[string]struct{})}
	r.Players = append(r.Players, &playerinfo)
}

func makeRoom(room_id int, private bool) {

	room := new(Room)
	room.Room_ID = room_id
	room.private = private
	room.round = 1
	room.wind = 0
	room.real_player_nun = 4
	roomlist[room_id] = room
	RROM := roomlist[room_id]
	RROM.recvchan = make(chan string, 2)
	go RROM.startgame(ctx)

}

func (r *Room) go_in_room(player *player_in, room_id int) bool {

	_, exist := roomlist[room_id]
	if exist {
		if len(roomlist[room_id].Players) < 4 {

			player.Room_ID = room_id
			roomlist[room_id].Addplayer(player)
			log.Println("[#FFA500]" + player.ID + "[reset] join room [yellow]" + strconv.Itoa(room_id) + "[reset]")
			player.conn.Write([]byte("True Room " + strconv.Itoa(room_id)))
			return true
		} else {
			player.conn.Write([]byte("False Room is full"))
			log.Printf("[red]ERROR:[reset] Room %d is full, %s can't get in\n", room_id, player.ID)
			return false
		}
	} else {
		player.conn.Write([]byte("False Room not exist"))
		log.Printf("[red]ERROR:[reset] Room %d not exist, %s can't get in\n", room_id, player.ID)
		return false
	}

}

func Room_finder(player *player_in) int {

	// 優先找人多的房間
	var rooms []*Room
	for _, r := range roomlist {
		if len(r.Players) < 4 && r.private == false && r.running == false {
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
			if r.running {
				r.Players[i].Ma = mao{}
				r.Players[i].Pong = make(map[string]struct{})
				r.Players[i].Chi = make(map[string]struct{})
				r.Players[i].Gang = make(map[string]struct{})
				r.real_player_nun--

			} else {
				r.Players = append(r.Players[:i], r.Players[i+1:]...)
			}
			player.Room_ID = -1
			log.Println("[#FFA500]" + player.ID + " [reset]leave room [yellow]" + strconv.Itoa(r.Room_ID) + "[reset]")
			player.conn.Write([]byte("True Leave room"))
			break
		}
	}
}

var cleaning bool

func RoomCleaner(ctx context.Context) {
	for {
		select {
		case <-time.After(50 * time.Second):
			cleaning = true
			log.Println("Start to clean room")
			for id, r := range roomlist {
				if len(r.Players) == 0 || r.real_player_nun == 0 {
					delete(roomlist, id)
					log.Println("Room[#FFA500]", id, "[reset]is empty, delete")
				}

			}
			log.Println("Room clean finish")
			cleaning = false
		case <-ctx.Done():
			return
		}

	}
}
