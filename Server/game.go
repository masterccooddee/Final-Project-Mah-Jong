package main

import (
	"math/rand/v2"
	"strconv"
	"sync"

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

// 除了now player之外的玩家能不能鳴牌
func (r *Room) MingCard(player *Player, card string) (count int) { //count 有幾個人能鳴牌
	// 判斷是否有人能鳴牌
	var po, g, c bool
	for _, p := range r.Players {
		po = false
		g = false
		c = false
		if p.ID == player.ID {
			continue
		}
		// 判斷能否碰
		if canPong(p, card) {
			// 處理碰牌邏輯
			po = true
		}
		// 判斷能否槓
		if canGang(p, card) {
			// 處理槓牌邏輯
			g = true
		}
		// 判斷能否吃 (只有下家可以吃)
		comb := canChi(p, card)
		if isNextPlayer(p) && comb != nil {
			// 處理吃牌邏輯
			c = true
		}

		if po || g || c {
			count++
		}
	}
	return count
}

func canPong(player *Player, card string) bool {
	count := 0
	for _, c := range player.Ma.Card {
		if c == card {
			count++
		}
	}
	return count >= 2
}

func canGang(player *Player, card string) bool {
	count := 0
	for _, c := range player.Ma.Card {
		if c == card {
			count++
		}
	}
	return count == 3
}

func canChi(player *Player, card string) (combinations [][]string) {
	// 假設牌是按順序存儲的
	// 需要判斷是否有連續的三張牌
	// 例如: card 是 3, 需要判斷是否有 1, 2 或 2, 4 或 4, 5

	if _, exist := player.Ma.Word[card]; exist {
		return nil
	}

	cardkind := string(card[0])
	cardvalue, _ := strconv.Atoi(card[1:])
	var combination []string

	// 1, 2
	if cardvalue > 2 {
		if player.HasCard(cardkind, cardvalue-2) && player.HasCard(cardkind, cardvalue-1) {
			combination = append(combination, cardkind+strconv.Itoa(cardvalue-2), cardkind+strconv.Itoa(cardvalue-1), card)
			combinations = append(combinations, combination)
			combination = nil
		}

	}
	// 2, 4
	if cardvalue > 1 && cardvalue < 9 {
		if player.HasCard(cardkind, cardvalue-1) && player.HasCard(cardkind, cardvalue+1) {
			combination = append(combination, cardkind+strconv.Itoa(cardvalue-1), card, cardkind+strconv.Itoa(cardvalue+1))
			combinations = append(combinations, combination)
			combination = nil
		}

	}
	// 4, 5
	if cardvalue < 8 {
		if player.HasCard(cardkind, cardvalue+1) && player.HasCard(cardkind, cardvalue+2) {
			combination = append(combination, card, cardkind+strconv.Itoa(cardvalue+1), cardkind+strconv.Itoa(cardvalue+2))
			combinations = append(combinations, combination)

		}

	}

	return combinations

}

func isNextPlayer(player *Player) bool {
	// 判斷是否為下家
	// 假設玩家順序存儲在 r.Players 中
	return player.Position == (now+1)%4

}

func (p *Player) HasPong(card string) bool {
	// 判斷是否有碰過該牌
	_, exist := p.Pong[card]
	return exist
}

func (p *Player) HasCard(cardkind string, cardValue int) bool {
	// 判斷是否有該牌

	switch cardkind {
	case "w":
		_, exist := p.Ma.Wan[cardValue]
		return exist
	case "t":
		_, exist := p.Ma.Tong[cardValue]
		return exist
	case "l":
		_, exist := p.Ma.Tiao[cardValue]
		return exist
	default:
		_, exist := p.Ma.Word[cardkind]
		return exist
	}

}

var now int //當前玩家

func (r *Room) startgame() {

	//確認是否有4個玩家
	// for len(r.Players) != 4 {
	// 	if _, exist := roomlist[r.Room_ID]; exist == false {
	// 		return
	// 	}
	// 	log.Println("Room", r.Room_ID, "is not full", len(r.Players))
	// 	time.Sleep(1 * time.Second)
	// }
	r.Addplayer(player_in{ID: "1", conn: nil})
	r.Addplayer(player_in{ID: "2", conn: nil})
	r.Addplayer(player_in{ID: "3", conn: nil})
	r.Addplayer(player_in{ID: "4", conn: nil})

	//通知所有玩家遊戲開始
	r.running = true
	//r.sendtoall("Game start")

	//隨機選座位
	//position := make(map[string]int)
	rand.Shuffle(len(r.Players), func(i, j int) { r.Players[i], r.Players[j] = r.Players[j], r.Players[i] })
	for i, p := range r.Players {
		p.Position = i
		//position[p.ID] = p.Position
	}
	//通知所有玩家座位
	//cli_pos, _ := json.Marshal(Position{Pos: position})
	//r.sendtoall(string(cli_pos))

	now = 0

	//發牌
	r.Cardset.addCard()
	for i := 0; i < 4; i++ {
		r.Players[i].Ma.Card = r.Cardset.Card[:13]
		r.Cardset.Card = r.Cardset.Card[13:]
		r.Players[i].Ma.splitCard()
	}
	num := len(r.Cardset.Card)
	r.lastcard = num - 14

	for r.round < 5 { //到東4局結束
		for r.lastcard > 0 { //直到剩0張牌

			//發一張牌
			r.Players[now].Ma.Card = append(r.Players[now].Ma.Card, r.Cardset.Card[0])
			r.Cardset.Card = r.Cardset.Card[1:]
			r.lastcard--

			//通知玩家抽到的牌
			sendtoplayer(r.Players[now].Ma.Card[len(r.Players[now].Ma.Card)-1], r.Players[now].ID)
			//通知其他玩家有人抽牌
			for i, p := range r.Players {
				if i != now {
					sendtoplayer("Draw", p.ID)
				}
			}
			//有沒有辦法胡牌、槓牌、碰牌

			//接收玩家出的牌
			getcard := <-r.recvchan
			//判斷有沒有人能鳴牌 (順序：胡、槓、碰、吃)
			r.MingCard(r.Players[now], getcard)

		}
	}

}
