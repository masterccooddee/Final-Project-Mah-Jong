package main

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
