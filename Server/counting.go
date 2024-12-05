package main

import (
	"sort"
)

type Tile struct {
	Suit  string // 花色，如 "萬:wan(w)", "筒:tong(t)", "索:tiao(l), "字:word"
	Value int    // 數字 (1~9) // 東:1 南:2 西:3 北:4 白:5 發:6 中:7
}

type Hand struct {
	Tiles    []Tile
	WinTile  Tile
	Melded   bool
	Ready    bool
	Declared bool
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

var tile_name = map[string]Tile{
	"w1": {Suit: "萬", Value: 1},
	"w2": {Suit: "萬", Value: 2},
	"w3": {Suit: "萬", Value: 3},
	"w4": {Suit: "萬", Value: 4},
	"w5": {Suit: "萬", Value: 5},
	"w6": {Suit: "萬", Value: 6},
	"w7": {Suit: "萬", Value: 7},
	"w8": {Suit: "萬", Value: 8},
	"w9": {Suit: "萬", Value: 9},
	"t1": {Suit: "筒", Value: 1},
	"t2": {Suit: "筒", Value: 2},
	"t3": {Suit: "筒", Value: 3},
	"t4": {Suit: "筒", Value: 4},
	"t5": {Suit: "筒", Value: 5},
	"t6": {Suit: "筒", Value: 6},
	"t7": {Suit: "筒", Value: 7},
	"t8": {Suit: "筒", Value: 8},
	"t9": {Suit: "筒", Value: 9},
	"l1": {Suit: "條", Value: 1},
	"l2": {Suit: "條", Value: 2},
	"l3": {Suit: "條", Value: 3},
	"l4": {Suit: "條", Value: 4},
	"l5": {Suit: "條", Value: 5},
	"l6": {Suit: "條", Value: 6},
	"l7": {Suit: "條", Value: 7},
	"l8": {Suit: "條", Value: 8},
	"l9": {Suit: "條", Value: 9},
	"1":  {Suit: "字", Value: 1},
	"2":  {Suit: "字", Value: 2},
	"3":  {Suit: "字", Value: 3},
	"4":  {Suit: "字", Value: 4},
	"5":  {Suit: "字", Value: 5},
	"6":  {Suit: "字", Value: 6},
	"7":  {Suit: "字", Value: 7},
}

var mao_name = map[Tile]string{
	{Suit: "萬", Value: 1}: "w1",
	{Suit: "萬", Value: 2}: "w2",
	{Suit: "萬", Value: 3}: "w3",
	{Suit: "萬", Value: 4}: "w4",
	{Suit: "萬", Value: 5}: "w5",
	{Suit: "萬", Value: 6}: "w6",
	{Suit: "萬", Value: 7}: "w7",
	{Suit: "萬", Value: 8}: "w8",
	{Suit: "萬", Value: 9}: "w9",
	{Suit: "筒", Value: 1}: "t1",
	{Suit: "筒", Value: 2}: "t2",
	{Suit: "筒", Value: 3}: "t3",
	{Suit: "筒", Value: 4}: "t4",
	{Suit: "筒", Value: 5}: "t5",
	{Suit: "筒", Value: 6}: "t6",
	{Suit: "筒", Value: 7}: "t7",
	{Suit: "筒", Value: 8}: "t8",
	{Suit: "筒", Value: 9}: "t9",
	{Suit: "條", Value: 1}: "l1",
	{Suit: "條", Value: 2}: "l2",
	{Suit: "條", Value: 3}: "l3",
	{Suit: "條", Value: 4}: "l4",
	{Suit: "條", Value: 5}: "l5",
	{Suit: "條", Value: 6}: "l6",
	{Suit: "條", Value: 7}: "l7",
	{Suit: "條", Value: 8}: "l8",
	{Suit: "條", Value: 9}: "l9",
	{Suit: "字", Value: 1}: "1",
	{Suit: "字", Value: 2}: "2",
	{Suit: "字", Value: 3}: "3",
	{Suit: "字", Value: 4}: "4",
	{Suit: "字", Value: 5}: "5",
	{Suit: "字", Value: 6}: "6",
	{Suit: "字", Value: 7}: "7",
}

func MaoToHand(m *mao) Hand {
	var hand Hand
	for _, tile := range m.Card {
		hand.Tiles = append(hand.Tiles, tile_name[tile])
	}
	return hand
}

func HandToMao(hand *Hand) mao {
	var m mao
	for _, tile := range hand.Tiles {
		m.Card = append(m.Card, mao_name[tile])
	}
	return m
}

// ***********************************************聽牌判定************************************************
func checkTenpai(hand Hand) ([]Tile, bool) {
	var winningTiles []Tile
	var ready bool
	uniqueTiles := generateAllTiles() // 產生所有可能的麻將牌

	for _, tile := range uniqueTiles {
		hand.Tiles = append(hand.Tiles, tile)
		if isWinningHand(hand) {
			winningTiles = append(winningTiles, tile)
			ready = true
		}
		hand.Tiles = hand.Tiles[:len(hand.Tiles)-1] // 移除模擬的牌
	}

	return winningTiles, ready
}

// 產生所有可能的麻將牌
func generateAllTiles() []Tile {
	var tiles []Tile
	suits := []string{"萬", "筒", "條"}
	for _, suit := range suits {
		for value := 1; value <= 9; value++ {
			tiles = append(tiles, Tile{Suit: suit, Value: value})
		}
	}
	honorValues := []int{1, 2, 3, 4, 5, 6, 7} // 東南西北白發中
	for _, value := range honorValues {
		tiles = append(tiles, Tile{Suit: "字", Value: value})
	}
	return tiles
}

// ***********************************************胡牌判定***********************************************
func isWinningHand(hand Hand) bool {
	// 檢查是否符合 4 面子 + 1 將的結構
	return checkForMentsuAndPair(hand) || checkSevenPairs(hand.Tiles)
}

// 檢查7對子
func checkSevenPairs(tiles []Tile) bool {
	if len(tiles) != 14 {
		return false
	}

	pairs := 0
	for i := 0; i < len(tiles)-1; i += 2 {
		if tiles[i] == tiles[i+1] {
			pairs++
		} else {
			return false
		}
	}

	return pairs == 7
}

// 檢查特殊牌型

// 檢查特殊牌型（此處僅處理槓子）
func checkGang(hand Hand) int {
	pairs := 0
	if len(hand.Tiles) > 14 {
		for i := 0; i < len(hand.Tiles)-3; i++ {
			if hand.Tiles[i] == hand.Tiles[i+1] && hand.Tiles[i+1] == hand.Tiles[i+2] && hand.Tiles[i+2] == hand.Tiles[i+3] {
				pairs += 1
			}
		}
	}

	return pairs
}

// 檢查是否有 4 個面子 + 1 將的輔助函式
// 檢查是否可以拆分為 4 面子 + 1 將
func checkForMentsuAndPair(hand Hand) bool {
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

	return findMentsuAndPair(tiles, 0, 0)
}

// 遞迴查找是否可以拆分成 4 面子 + 1 將
func findMentsuAndPair(tiles []Tile, mentsuCount, pairCount int) bool {
	/* for _, tile := range tiles {
		fmt.Print(tile)
	}
	fmt.Println(mentsuCount, pairCount)
	fmt.Println() */
	if len(tiles) == 0 {
		//fmt.Println("========================================================")
		return mentsuCount == 4 && pairCount == 1
	}

	// 嘗試找將牌（雀頭）
	if pairCount == 0 && len(tiles) >= 2 && tiles[0] == tiles[1] {
		//fmt.Println("head: ", tiles[0], tiles[1])
		if findMentsuAndPair(tiles[2:], mentsuCount, pairCount+1) {
			return true
		}
	}

	// 嘗試找順子
	for i := 1; i < len(tiles)-1; i++ {
		for j := i + 1; j < len(tiles); j++ {
			if !isHonorTile(tiles[0]) {
				if len(tiles) >= 3 && isSequential(tiles[0], tiles[i], tiles[j]) {
					remain := removeTiles(tiles, []Tile{tiles[0], tiles[i], tiles[j]})
					if findMentsuAndPair(remain, mentsuCount+1, pairCount) {
						return true
					}
				}
			}
		}
	}
	/* if len(tiles) >= 3 && isSequential(tiles[0], tiles[1], tiles[2]) {
		fmt.Printf("sequence%d: ", mentsuCount)
		fmt.Println(tiles[0], tiles[1], tiles[2])
		if findMentsuAndPair(tiles[3:], mentsuCount+1, pairCount) {
			return true
		}
	} */

	// 嘗試找刻子
	if len(tiles) >= 3 && tiles[0] == tiles[1] && tiles[1] == tiles[2] {
		//fmt.Printf("triplet%d: ", mentsuCount)
		//fmt.Println(tiles[0], tiles[1], tiles[2])
		if findMentsuAndPair(tiles[3:], mentsuCount+1, pairCount) {
			return true
		}
	}

	return false
}

// 判斷是否為順子
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

// ***********************************************飜數判定***********************************************
//1番：riichi(立直) fully_conceal(門清) unbroken(一發) all_inside(斷么九) pinfu(平和) twin_seq(一盃口) honor(場風/自風/中發白)
//2番：double_riichi(天聽) seven_pair(七對子) full_straight(一氣通貫) mixed_seq(三色同順)
//3番：
//6番：
//役滿：
/* func calculateHan(hand Hand) ([]Yaku, int) {
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
*/
func isRiichi(player *Player, discarded Tile) bool {
	// 條件 1：不能有副露
	var Melded, Ready, Declared bool
	for _, tile := range player.Ma.Card {
		if player.HasChi(tile) || player.HasPong(tile) || player.HasGang(tile) {
			Melded = true
			break
		}
	}
	if Melded {
		return false
	}
	// 條件 2：必須處於聽牌狀態
	_, Ready = checkTenpai(MaoToHand(&player.Ma))
	if !Ready {
		return false
	}
	// 條件 3：立直必須在打出牌後
	if discarded == (Tile{}) {
		return false
	}
	// 條件 4：不能改變手牌
	if Declared {
		return false
	}
	// 所有條件滿足
	return true
}

// ***********************************************符數判定***********************************************
/*
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
} */

// ***********************************************判定結果***********************************************
/* func calculateHandResult(hand Hand) HandResult {
	var hand := MaoToHand(player.Ma)
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
} */

/* func main() {
	// 定義測試牌
	testTiles := []Tile{
		{Suit: "萬", Value: 4},
		{Suit: "萬", Value: 4},
		{Suit: "萬", Value: 4},
		{Suit: "萬", Value: 5},
		{Suit: "萬", Value: 6},
		{Suit: "筒", Value: 1},
		{Suit: "筒", Value: 2},
		{Suit: "筒", Value: 3},
		{Suit: "條", Value: 4},
		{Suit: "條", Value: 5},
		{Suit: "條", Value: 6},
		{Suit: "條", Value: 5},
		{Suit: "條", Value: 6},
		{Suit: "條", Value: 7},
	}

	hand := Hand{Tiles: testTiles}
	isWinningHand(hand)
	if isWinningHand(hand) {
		println("這是一手胡牌！")
	} else {
		println("這不是一手胡牌。")
	}
}
*/
