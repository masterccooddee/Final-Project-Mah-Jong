package ui

import (
	"math/rand"
	"sort"
	"strconv"
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
	wanKeys := make([]int, 0, len(m.Wan))
	for k := range m.Wan {
		wanKeys = append(wanKeys, k)
	}
	sort.Ints(wanKeys)

	tiaoKeys := make([]int, 0, len(m.Tiao))
	for k := range m.Tiao {
		tiaoKeys = append(tiaoKeys, k)
	}
	sort.Ints(tiaoKeys)

	tongKeys := make([]int, 0, len(m.Tong))
	for k := range m.Tong {
		tongKeys = append(tongKeys, k)
	}
	sort.Ints(tongKeys)

	wordKeys := make([]string, 0, len(m.Word))
	for k := range m.Word {
		wordKeys = append(wordKeys, k)
	}
	sort.Strings(wordKeys)

	// 清空 Card 切片
	m.Card = nil
	// 合併切片
	for _, k := range wanKeys {
		m.Card = append(m.Card, "w"+strconv.Itoa(k))
	}
	for _, k := range tiaoKeys {
		m.Card = append(m.Card, "l"+strconv.Itoa(k))
	}
	for _, k := range tongKeys {
		m.Card = append(m.Card, "t"+strconv.Itoa(k))
	}
	m.Card = append(m.Card, wordKeys...)

}

func (c *mao) removeCard(remove string) {
	for i, v := range c.Card {
		if v == remove {
			c.Card = append(c.Card[:i], c.Card[i+1:]...)
			break
		}
	}

}
