package main

import (
	"sort"
)

type Tile struct {
	Suit  string // 花色，如 "萬:wan(w)", "筒:tong(t)", "索:tiao(l), "字:word"
	Value int    // 數字 (1~9) // 東:1 南:2 西:3 北:4 白:5 發:6 中:7
}

type Hand struct {
	Tiles   []Tile
	WinTile Tile
}

type Yaku struct {
	Name string
	Han  int
}

type HandResult struct {
	Yaku      []Yaku
	Han       int
	Fu        int
	IsWinning bool
}

// ************************************************胡牌判定************************************************
func isWinningHand(hand Hand) bool {
	// 檢查是否符合 4 面子 + 1 將的結構
	return checkForMentsuAndPair(hand)
}

// 檢查是否有 4 個面子 + 1 將的輔助函數
func checkForMentsuAndPair(hand Hand) bool {
	// 這裡需要一個完整的面子拆解和配對的判定邏輯
	// 因篇幅問題，這裡不詳細實現。
	// 假設符合條件，返回 true；否則返回 false
	tiles := hand.Tiles

	if len(tiles) != 14 {
		return false
	}

	sort.Slice(tiles, func(i, j int) bool {
		if tiles[i].Suit == tiles[j].Suit {
			return tiles[i].Value < tiles[j].Value
		}
		return tiles[i].Suit < tiles[j].Suit
	})

	for i := 0; i < len(tiles)-2; i++ {
		if tiles[i] == tiles[i+1] {
			remainingTiles := append([]Tile{}, tiles[:i]...)
			remainingTiles = append(remainingTiles, tiles[i+2:]...)

			if canFormMentsu(remainingTiles) {
				return true
			}
		}
	}
	return false // 替換為實際邏輯
}

func canFormMentsu(tiles []Tile) bool {
	if len(tiles) == 0 {
		return true // 沒有剩餘牌，表示成功拆解成面子
	}
	if len(tiles) < 3 {
		return false // 剩餘牌不足以形成面子
	}

	// 如果是字牌，只能形成刻子
	if isHonorTile(tiles[0]) {
		if len(tiles) >= 3 && tiles[0] == tiles[1] && tiles[1] == tiles[2] {
			return canFormMentsu(tiles[3:]) // 若形成刻子則繼續遞迴
		}
		return false // 無法形成有效面子
	}

	// 對數牌嘗試用刻子（例如 "5筒、5筒、5筒"）拆解
	if len(tiles) >= 3 && tiles[0] == tiles[1] && tiles[1] == tiles[2] {
		if canFormMentsu(tiles[3:]) {
			return true
		}
	}

	// 對數牌嘗試用順子（例如 "1萬、2萬、3萬"）拆解
	for i := 1; i < len(tiles)-1; i++ {
		for j := i + 1; j < len(tiles); j++ {
			if isSequential(tiles[0], tiles[i], tiles[j]) {
				remainingTiles := removeTiles(tiles, []Tile{tiles[0], tiles[i], tiles[j]})
				if canFormMentsu(remainingTiles) {
					return true
				}
			}
		}
	}

	return false
}

// 判斷三張牌是否為順子（同花色且連續數字）
func isSequential(a, b, c Tile) bool {
	return a.Suit == b.Suit && b.Suit == c.Suit &&
		a.Value+1 == b.Value && b.Value+1 == c.Value
}

// 判斷是否為字牌
func isHonorTile(tile Tile) bool {
	return tile.Suit == "字" // 字牌的 Suit 設定為 "字"
}

func isWanTile(tile Tile) bool {
	return tile.Suit == "萬"
}

func isTongTile(tile Tile) bool {
	return tile.Suit == "筒"
}

func isTiaoTile(tile Tile) bool {
	return tile.Suit == "條"
}

// 移除指定的牌
func removeTiles(tiles, toRemove []Tile) []Tile {
	result := append([]Tile{}, tiles...)
	for _, removeTile := range toRemove {
		for i, tile := range result {
			if tile == removeTile {
				result = append(result[:i], result[i+1:]...)
				break
			}
		}
	}
	return result
}

// ************************************************飜數判定************************************************
func calculateHan(hand Hand) ([]Yaku, int) {
	yakuList := []Yaku{}

	// 舉例：檢查平和（門清門前和牌無加符）
	if isPinfu(hand) {
		yakuList = append(yakuList, Yaku{Name: "平和", Han: 1})
	}

	// 檢查其他役種，可以依據需要添加更多邏輯
	// 如一盃口、立直等...

	// 計算總飜數
	han := 0
	for _, yaku := range yakuList {
		han += yaku.Han
	}
	return yakuList, han
}

// 示例：檢查平和
func isPinfu(hand Hand) bool {
	// 實際檢查平和的判斷邏輯
	return true // 替換為實際邏輯
}

// ************************************************符數判定************************************************
func calculateFu(hand Hand) int {
	fu := 20 // 胡牌基礎符數為 20 符

	// 示例：根據不同條件調整符數
	// 如果手牌中包含暗刻或刻子，增加符數
	if hasAnko(hand) {
		fu += 4 // 每一暗刻增加 4 符
	}

	return fu
}

// 示例：判斷是否有暗刻
func hasAnko(hand Hand) bool {
	return false // 替換為實際邏輯
}

// ************************************************判定結果************************************************
func calculateHandResult(hand Hand) HandResult {
	isWinning := isWinningHand(hand)

	if !isWinning {
		return HandResult{
			IsWinning: false,
		}
	}

	yaku, han := calculateHan(hand)
	fu := calculateFu(hand)

	return HandResult{
		Yaku:      yaku,
		Han:       han,
		Fu:        fu,
		IsWinning: true,
	}
}
