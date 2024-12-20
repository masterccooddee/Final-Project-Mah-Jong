package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-zeromq/zmq4"
)

var router zmq4.Socket
var selftouch bool

var (
	order = map[string]int{
		"Hu":   1,
		"Gang": 2,
		"Pong": 3,
		"Chi":  4,
	}
)

type Position struct {
	Pos map[string]int
	Ma  mao
}

var zmqmu sync.Mutex
var pos_history []int

func sendtoplayer(msg string, ID string) {
	msgout := zmq4.NewMsgFrom([]byte(ID), []byte(msg))
	zmqmu.Lock()
	zmqloger.Printf("Send to [#FFA500]%s[reset] (Room [yellow]%d[reset]): %s", ID, playerlist[ID].Room_ID, msg)
	router.SendMulti(msgout)
	zmqmu.Unlock()
}

func (r *Room) sendtoall(msg string) {
	for _, p := range r.Players {
		sendtoplayer(msg, p.ID)
	}
}

// 除了now player之外的玩家能不能鳴牌
func (r *Room) MingCard(player *Player, card string, now int) (count int) { //count 有幾個人能鳴牌
	// 判斷是否有人能鳴牌
	var po, g, c, h bool
	var msgcomb string

	for _, p := range r.Players {
		if p.ID == player.ID {
			continue
		}
		// 判斷能否碰
		if canPong(p, card) {
			// 處理碰牌邏輯
			msgcomb += "Pong " + card + ","
			po = true
		}
		// 判斷能否槓
		if canGang(p, card) {
			// 處理槓牌邏輯
			msgcomb += "Gang " + card + ","
			g = true
		}
		// 判斷能否吃 (只有下家可以吃)
		comb := canChi(p, card)
		if isNextPlayer(p, now) && comb != nil {
			// 處理吃牌邏輯

			// 傳出去是數字，代表第幾種組合，0: card,x,x 1: x,card,x 2: x,x,card
			var chi string
			for _, c := range comb {
				chi += c + " "
			}
			chi = chi[:len(chi)-1]

			msgcomb += "Chi " + chi + ","
			c = true
		}

		hucard := p.checkcard()
		hucard.Card = append(hucard.Card, card)
		hucard.splitCard()
		if isWinningHand(MaoToHand(&hucard)) {
			msgcomb += "Hu " + card + ","
			h = true

		}

		if po || g || c || h {
			count++
			msgcomb = strings.TrimRight(msgcomb, ",")
			sendtoplayer(msgcomb, p.ID)
		} else {
			sendtoplayer("MING FROM OTHER NO", p.ID)
		}
		po, g, c, h = false, false, false, false
		msgcomb = ""

	}
	return count
}

func isNextPlayer(player *Player, now int) bool {
	// 判斷是否為下家
	// 假設玩家順序存儲在 r.Players 中
	return player.Position == (now+1)%4

}

var windong bool

func (r *Room) endgame(now int) {
	var getpoint int
	if r.Players[now].Position == 0 {
		if selftouch {
			getpoint = 300*r.bunround + 3000
			r.Players[now].Point += getpoint
			for i, _ := range r.Players {
				if i != now {
					r.Players[i].Point -= getpoint / 3
				}
			}
		} else {
			getpoint = 300*r.bunround + 1000
			r.Players[now].Point += getpoint
			r.Players[pos_history[0]].Point -= getpoint
		}

		r.bunround++
		r.round--
		windong = true
	} else {
		if selftouch {
			getpoint = 300*r.bunround + 1500
			r.Players[now].Point += getpoint
			for i, _ := range r.Players {
				if i != now {
					if i == 0 {
						r.Players[i].Point -= getpoint / 2
					} else {
						r.Players[i].Point -= getpoint / 4
					}
				}
			}
		} else {
			getpoint = 300*r.bunround + 300
			r.Players[now].Point += getpoint
			r.Players[pos_history[0]].Point -= getpoint
		}
		r.bunround = 0
	}
	pointstr := fmt.Sprintf("%d %d %d %d", r.Players[0].Point, r.Players[1].Point, r.Players[2].Point, r.Players[3].Point)
	r.sendtoall(pointstr)
}

func makeFromSlice(sl []string) []string {
	result := make([]string, len(sl))
	copy(result, sl)
	return result
}

func (r *Room) startgame(ctx context.Context) {

	select {
	case <-ctx.Done():
		return
	default:
	start:
		var now int //當前玩家
		var baocard []string
		var hGang bool
		var hPong bool
		var hChi bool
		roundmax := 5
		//確認是否有4個玩家
		for len(r.Players) != 4 {
			if _, exist := roomlist[r.Room_ID]; exist == false {
				return
			}
			time.Sleep(1 * time.Second)
		}

		//通知所有玩家遊戲開始
		r.running = true
		r.sendtoall("Game start")

		//隨機選座位、發牌
		r.Cardset.addCard()
		position := make(map[string]int)
		rand.Shuffle(len(r.Players), func(i, j int) { r.Players[i], r.Players[j] = r.Players[j], r.Players[i] })

		var cli_info Position
		for i, p := range r.Players {
			p.Position = i
			position[p.ID] = p.Position
			p.Ma.Card = r.Cardset.Card[:13]
			//log.Println(p.Ma.Card)
			r.Cardset.Card = r.Cardset.Card[13:]
			p.Ma.splitCard()

		}

		for _, p := range r.Players {
			//打包座位、手牌並發送給玩家
			cli_info.Pos = position
			cli_info.Ma = p.Ma
			cli_pos, _ := json.Marshal(cli_info)
			sendtoplayer(string(cli_pos), p.ID)
		}

		now = 0
		r.now = now
		pos_history = append(pos_history, now)

		num := len(r.Cardset.Card)
		r.lastcard = num - 14
		baocard = r.Cardset.Card[num-14:]
		log.Println(baocard)

		for r.round < roundmax { //到東4局結束

			for r.lastcard > 0 { //直到剩0張牌
				//發一張牌
				if !hPong && !hChi {
					if !hGang {
						//newcard = r.Cardset.Card[0]
						newcardset := append(makeFromSlice(r.Players[now].Ma.Card), r.Cardset.Card[0])
						r.Players[now].Ma.Card = newcardset
						r.Cardset.Card = r.Cardset.Card[1:]

					} else { //槓拿嶺上牌
						hGang = false
						r.Players[now].Ma.Card = append(r.Players[now].Ma.Card, r.Cardset.Card[len(r.Cardset.Card)-1])
						r.Cardset.Card = r.Cardset.Card[:len(r.Cardset.Card)-1]
					}
					r.lastcard--

					//通知玩家抽到的牌
					sendtoplayer(r.Players[now].Ma.Card[len(r.Players[now].Ma.Card)-1], r.Players[now].ID)
					//通知其他玩家有人抽牌

					for i, p := range r.Players {
						if i != now {
							sendtoplayer("Draw", p.ID)
						}
					}

					var selfmsg string
					//有沒有辦法胡牌、槓牌

					hucard := r.Players[now].checkcard()
					if isWinningHand(MaoToHand(&hucard)) {
						//胡牌
						selfmsg += "Hu " + r.Players[now].Ma.Card[len(r.Players[now].Ma.Card)-1] + ","

					}
					if canGangself(r.Players[now], r.Players[now].Ma.Card[len(r.Players[now].Ma.Card)-1]) || r.Players[now].HasPong(r.Players[now].Ma.Card[len(r.Players[now].Ma.Card)-1]) {
						//槓牌
						selfmsg += "Gang " + r.Players[now].Ma.Card[len(r.Players[now].Ma.Card)-1] + ","

					}

					if selfmsg != "" {
						selfmsg = strings.TrimRight(selfmsg, ",")
						sendtoplayer(selfmsg, r.Players[now].ID)

						getcard := strings.TrimSpace(<-r.recvchan)
						getslice := strings.Split(getcard, " ")
						if getslice[1] == "Cancel" {
							goto nottaken
						}
						for r.Players[now].ID != getslice[0] {
							getcard = strings.TrimSpace(<-r.recvchan)
							getslice = strings.Split(getcard, " ")
						}
						getcard = getslice[1]
						if getcard == "Gang" {
							//槓牌
							if canGang(r.Players[now], r.Players[now].Ma.Card[len(r.Players[now].Ma.Card)-1]) {
								for i := 0; i < 3; i++ {
									r.Players[now].Ma.removeCard(r.Players[now].Ma.Card[len(r.Players[now].Ma.Card)-1])
								}
							} else {
								delete(r.Players[now].Pong, r.Players[now].Ma.Card[len(r.Players[now].Ma.Card)-1])
							}
							r.Players[now].Ma.splitCard()
							r.Players[now].Gang[r.Players[now].Ma.Card[len(r.Players[now].Ma.Card)-1]] = struct{}{}
							hGang = true

							r.sendtoall("Gang " + r.Players[now].Ma.Card[len(r.Players[now].Ma.Card)-1] + " " + r.Players[now].ID)

							continue
						} else if getcard == "Hu" {
							r.sendtoall("Hu " + r.Players[now].Ma.Card[len(r.Players[now].Ma.Card)-1] + " " + r.Players[now].ID)
							selftouch = true
							r.endgame(now)
							selftouch = false
							goto nextround
						}

					} else {
						for _, p := range r.Players {
							//if p.ID != r.Players[now].ID {
							sendtoplayer("MING FROM SELF NO", p.ID)
							//}
						}
					}

				}
				//接收玩家出的牌
			nottaken:
				outcard := strings.TrimSpace(<-r.recvchan)
				outcardslice := strings.Split(outcard, " ")
				for r.Players[now].ID != outcardslice[0] {
					outcard = strings.TrimSpace(<-r.recvchan)
					outcardslice = strings.Split(outcard, " ")
				}
				outcard = outcardslice[1]
				log.Println("Player", r.Players[now].ID, "discard", outcard)
				showcardmsg := strconv.Itoa(now) + " " + outcard
				for _, p := range r.Players {
					if p.ID != r.Players[now].ID {
						sendtoplayer(showcardmsg, p.ID)
					}
				}
				//r.sendtoall(showcardmsg)
				hChi = false
				hPong = false
				r.Players[now].Ma.removeCard(outcard)
				r.Players[now].Ma.splitCard()

				//判斷有沒有人能鳴牌 (順序：胡、槓、碰、吃)
				var ming []string
				cnt := r.MingCard(r.Players[now], outcard, now)

			loop:
				for i := 0; i < (cnt); i++ {
					//接收玩家鳴牌
					getcard := strings.TrimSpace(<-r.recvchan)
					getcardslice := strings.Split(getcard, " ")
					if getcardslice[1] == "Cancel" {
						continue
					} else {
						getcard = getcardslice[1] + " " + getcardslice[2] + " " + getcardslice[3]
					}

					ming = append(ming, getcard)
					cnt--
					if cnt == 0 {
						break loop
					}

					//有人鳴牌、倒數1秒可以鳴牌
					timer := time.NewTimer(1 * time.Second)

					for {
						select {
						case getcard = <-r.recvchan:
							getcardslice := strings.Split(getcard, " ")
							if getcardslice[1] == "Cancel" {
								continue
							} else {
								getcard = getcardslice[1] + " " + getcardslice[2] + " " + getcardslice[3]
							}

							ming = append(ming, getcard)
							cnt--
							if cnt == 0 {
								break loop
							}
						case <-timer.C:
							break loop
						}
					}
				}

				//清空channel
			clean:
				for {
					select {
					case <-r.recvchan:
					default:
						break clean
					}
				}
				//格式： 動作 位置 牌 ex: Pong 1 w1, Chi 1 0 (0是種類)
				//回傳： 動作 牌 playerID  ex: Pong w1 hehehe, Chi 0 hehehe (0是種類)
				if ming != nil {

					sort.Slice(ming, func(i, j int) bool {
						msgI := strings.Split(ming[i], " ")
						msgJ := strings.Split(ming[j], " ")
						return order[msgI[0]] < order[msgJ[0]]
					})

					msg := strings.Split(ming[0], " ")
					pos, _ := strconv.Atoi(msg[1])
					if msg[0] == "Hu" {
						//胡牌
						r.sendtoall("Hu " + r.Players[now].Ma.Card[len(r.Players[now].Ma.Card)-1] + " " + r.Players[now].ID)
						r.endgame(pos)
						goto nextround

					}
					if msg[0] == "Gang" {
						//槓牌

						for i := 0; i < 3; i++ {
							r.Players[pos].Ma.removeCard(outcard)
						}
						r.Players[pos].Ma.splitCard()
						r.Players[pos].clean = false
						r.Players[pos].Gang[outcard] = struct{}{}
						hGang = true

						now = pos
						r.now = now
						pos_history = pos_history[1:]
						pos_history = append(pos_history, now)

						for _, p := range r.Players {
							//if p.ID != r.Players[now].ID {
							sendtoplayer("Gang "+outcard+" "+r.Players[now].ID, p.ID)
							//}
						}
						continue
					}

					if msg[0] == "Pong" {
						//碰牌
						hPong = true
						for i := 0; i < 2; i++ {
							r.Players[pos].Ma.removeCard(outcard)
						}
						r.Players[pos].Ma.splitCard()
						r.Players[pos].clean = false
						r.Players[pos].Pong[outcard] = struct{}{}

						now = pos
						r.now = now
						pos_history = pos_history[1:]
						pos_history = append(pos_history, now)

						for _, p := range r.Players {
							//if p.ID != r.Players[now].ID {
							sendtoplayer("Pong "+outcard+" "+r.Players[now].ID, p.ID)
							//}
						}
						continue
					}

					if msg[0] == "Chi" {
						hChi = true
						r.Players[pos].clean = false

						num, _ := strconv.Atoi(string(outcard[1]))
						switch msg[2] {
						case "0":
							r.Players[pos].Chi[outcard+" "+string(outcard[0])+strconv.Itoa(num+1)+" "+string(outcard[0])+strconv.Itoa(num+2)] = struct{}{}
							r.Players[pos].Ma.removeCard(string(outcard[0]) + strconv.Itoa(num+1))
							r.Players[pos].Ma.removeCard(string(outcard[0]) + strconv.Itoa(num+2))
						case "1":
							r.Players[pos].Chi[string(outcard[0])+strconv.Itoa(num-1)+" "+outcard+" "+string(outcard[0])+strconv.Itoa(num+1)] = struct{}{}
							r.Players[pos].Ma.removeCard(string(outcard[0]) + strconv.Itoa(num-1))
							r.Players[pos].Ma.removeCard(string(outcard[0]) + strconv.Itoa(num+1))
						case "2":
							r.Players[pos].Chi[string(outcard[0])+strconv.Itoa(num-2)+" "+string(outcard[0])+strconv.Itoa(num-1)+" "+outcard] = struct{}{}
							r.Players[pos].Ma.removeCard(string(outcard[0]) + strconv.Itoa(num-2))
							r.Players[pos].Ma.removeCard(string(outcard[0]) + strconv.Itoa(num-1))
						}
						r.Players[pos].Ma.splitCard()

						now = pos
						r.now = now
						pos_history = pos_history[1:]
						pos_history = append(pos_history, now)

						for _, p := range r.Players {
							//if p.ID != r.Players[now].ID {
							sendtoplayer("Chi "+msg[2]+" "+r.Players[now].ID, p.ID)
							//}
						}
						continue
					}

					//處理順序：胡、槓、碰、吃
				} else {

					r.sendtoall("WAITING MING NO")

				}
				ming = nil

				now = (now + 1) % 4
				r.now = now
				pos_history = pos_history[1:]
				pos_history = append(pos_history, now)
			}
		nextround:
			r.round++
			if r.round == roundmax && r.bunround == 0 {
				goto endgame
			}
			now = 0
			r.now = now
			pos_history = nil
			pos_history = append(pos_history, now)
			position := make(map[string]int)
			if !windong {
				r.Players = append(r.Players[1:], r.Players[0])
			} else {
				windong = false
			}
			r.Cardset = mao{}
			r.Cardset.addCard()
			r.Cardset.splitCard()
			for i, p := range r.Players {
				p.Position = i
				p.clean = true
				p.Chi = make(map[string]struct{})
				p.Pong = make(map[string]struct{})
				p.Gang = make(map[string]struct{})
				p.TingCard = false
				p.Ma.Card = r.Cardset.Card[:13]
				p.Ma.splitCard()
				r.Cardset.Card = r.Cardset.Card[13:]
			}
			num = len(r.Cardset.Card)
			r.lastcard = num - 14
			baocard = r.Cardset.Card[num-14:]

			r.sendtoall("Next round")
			for _, p := range r.Players {
				position[p.ID] = p.Position
			}

			for _, p := range r.Players {
				//打包座位、手牌並發送給玩家
				cli_info.Pos = position
				cli_info.Ma = p.Ma
				cli_pos, _ := json.Marshal(cli_info)
				sendtoplayer(string(cli_pos), p.ID)
			}

		}
	endgame:
		r.sendtoall("Game over")
		for _, p := range r.Players {
			r.leave_room(p.player_in)
		}
		r.running = false
		r.Players = nil
		r.Cardset = mao{}
		r.round = 1
		r.wind = 0
		r.bunround = 0
		r.gang = [4]int{}
		r.Players = nil
		if r.private {
			delete(roomlist, r.Room_ID)
		} else {
			time.Sleep(2 * time.Second)
			goto start
		}

	}

}
