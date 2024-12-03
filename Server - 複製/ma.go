package main

import (
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

type mao struct {
	Card []string
	//牌的set
	Wan  map[int]struct{}
	Tong map[int]struct{}
	Tiao map[int]struct{}
	Word map[string]struct{}
}

// 東 : 1 南 : 2 西 : 3 北 : 4 白 : 5 發 : 6  中: 7
func (m *mao) addCard() {
	kind := []string{"w", "t", "l"}

	for i := 0; i < 4; i++ {
		for _, v := range kind {
			for i := 1; i < 10; i++ {
				word := v + strconv.Itoa(i)
				m.Card = append(m.Card, word)
			}
		}
		m.Card = append(m.Card, "1") //東
		m.Card = append(m.Card, "2") //南
		m.Card = append(m.Card, "3") //西
		m.Card = append(m.Card, "4") //北
		m.Card = append(m.Card, "5") //白
		m.Card = append(m.Card, "6") //發
		m.Card = append(m.Card, "7") //中
	}

	rand.Shuffle(len(m.Card), func(i, j int) { m.Card[i], m.Card[j] = m.Card[j], m.Card[i] })

}

func (m *mao) splitCard() {
	m.Wan = make(map[int]struct{})
	m.Tong = make(map[int]struct{})
	m.Tiao = make(map[int]struct{})
	m.Word = make(map[string]struct{})
	for _, v := range m.Card {

		switch v[0] {
		case 'w':
			num, _ := strconv.Atoi(v[1:])
			m.Wan[num] = struct{}{}
		case 't':
			num, _ := strconv.Atoi(v[1:])
			m.Tong[num] = struct{}{}
		case 'l':
			num, _ := strconv.Atoi(v[1:])
			m.Tiao[num] = struct{}{}
		default:
			m.Word[v] = struct{}{}
		}
	}

}

func (m *mao) SortCard() {

	var wan []string
	var tiao []string
	var tong []string
	var word []string

	for _, k := range m.Card {
		switch k[0] {
		case 'w':
			wan = append(wan, k)
		case 't':
			tong = append(tong, k)
		case 'l':
			tiao = append(tiao, k)
		default:
			word = append(word, k)
		}
	}
	sort.Slice(wan, func(i, j int) bool { return wan[i][1] < wan[j][1] })
	sort.Slice(tiao, func(i, j int) bool { return tiao[i][1] < tiao[j][1] })
	sort.Slice(tong, func(i, j int) bool { return tong[i][1] < tong[j][1] })
	sort.Slice(word, func(i, j int) bool { return word[i] < word[j] })
	// 清空 Card 切片
	m.Card = nil
	// 合併切片
	m.Card = append(m.Card, wan...)
	m.Card = append(m.Card, tong...)
	m.Card = append(m.Card, tiao...)
	m.Card = append(m.Card, word...)

}

func (c *mao) removeCard(remove string) {
	for i, v := range c.Card {
		if v == remove {
			c.Card = append(c.Card[:i], c.Card[i+1:]...)
			break
		}
	}

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

func canGangself(player *Player, card string) bool {
	count := 0
	for _, c := range player.Ma.Card {
		if c == card {
			count++
		}
	}
	return count == 4
}

func canChi(player *Player, card string) (combinations []string) {
	// 假設牌是按順序存儲的
	// 需要判斷是否有連續的三張牌
	// 例如: card 是 3, 需要判斷是否有 1, 2 或 2, 4 或 4, 5

	if _, exist := player.Ma.Word[card]; exist {
		return nil
	}

	cardkind := string(card[0])
	cardvalue, _ := strconv.Atoi(card[1:])

	// 4, 5
	if cardvalue < 8 {
		if player.HasCard(cardkind, cardvalue+1) && player.HasCard(cardkind, cardvalue+2) {
			combinations = append(combinations, "0")

		}

	}

	// 2, 4
	if cardvalue > 1 && cardvalue < 9 {
		if player.HasCard(cardkind, cardvalue-1) && player.HasCard(cardkind, cardvalue+1) {
			combinations = append(combinations, "1")
		}

	}

	// 1, 2
	if cardvalue > 2 {
		if player.HasCard(cardkind, cardvalue-2) && player.HasCard(cardkind, cardvalue-1) {
			combinations = append(combinations, "2")
		}

	}
	return combinations

}

func (p *Player) HasPong(card string) bool {
	// 判斷是否有碰過該牌
	_, exist := p.Pong[card]
	return exist
}

func (p *Player) HasChi(card string) bool {
	// 判斷是否有吃過該牌
	_, exist := p.Chi[card]
	return exist
}

func (p *Player) HasGang(card string) bool {
	// 判斷是否有槓過該牌
	_, exist := p.Gang[card]
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

func (p *Player) checkcard() (ma mao) {
	ma = p.Ma
	for k, _ := range p.Pong {
		for i := 0; i < 3; i++ {
			ma.Card = append(ma.Card, k)
		}
	}
	for k, _ := range p.Chi {
		cc := strings.Split(k, " ")
		for _, v := range cc {
			ma.Card = append(ma.Card, v)
		}
	}
	for k, _ := range p.Gang {
		for i := 0; i < 4; i++ {
			ma.Card = append(ma.Card, k)
		}
	}
	ma.SortCard()
	ma.splitCard()
	return
}
