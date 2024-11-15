package ui

import (
	"math/rand"
	"sort"
	"strconv"
)

type mao struct {
	Card []string
	Wan  []string
	Tong []string
	Tiao []string
	Word []string
}

func (m *mao) addCard() {
	kind := []string{"w", "t", "l"}

	for i := 0; i < 4; i++ {
		for _, v := range kind {
			for i := 1; i < 10; i++ {
				word := v + strconv.Itoa(i)
				m.Card = append(m.Card, word)
			}
		}
		m.Card = append(m.Card, "bai")
		m.Card = append(m.Card, "zhong")
		m.Card = append(m.Card, "fa")
		m.Card = append(m.Card, "dong")
		m.Card = append(m.Card, "nan")
		m.Card = append(m.Card, "xi")
		m.Card = append(m.Card, "bei")
	}

	rand.Shuffle(len(m.Card), func(i, j int) { m.Card[i], m.Card[j] = m.Card[j], m.Card[i] })

}

func (m *mao) splitCard() {
	for _, v := range m.Card {

		switch v[0] {
		case 'w':
			m.Wan = append(m.Wan, v)
		case 't':
			m.Tong = append(m.Tong, v)
		case 'l':
			m.Tiao = append(m.Tiao, v)
		default:
			m.Word = append(m.Word, v)
		}
	}

}

func (m *mao) SortCard() {
	sort.Slice(m.Wan, func(i, j int) bool { return m.Wan[i][1] < m.Wan[j][1] })
	sort.Slice(m.Tiao, func(i, j int) bool { return m.Tiao[i][1] < m.Tiao[j][1] })
	sort.Slice(m.Tong, func(i, j int) bool { return m.Tong[i][1] < m.Tong[j][1] })
	sort.Slice(m.Word, func(i, j int) bool { return m.Word[i] < m.Word[j] })
	// 清空 Card 切片
	m.Card = nil
	// 合併切片
	m.Card = append(m.Card, m.Wan...)
	m.Card = append(m.Card, m.Tong...)
	m.Card = append(m.Card, m.Tiao...)
	m.Card = append(m.Card, m.Word...)

}

func (c *mao) removeCard(remove string) {
	for i, v := range c.Card {
		if v == remove {
			c.Card = append(c.Card[:i], c.Card[i+1:]...)
			break
		}
	}

}
