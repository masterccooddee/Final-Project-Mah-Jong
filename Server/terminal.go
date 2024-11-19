package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var textView *tview.TextView
var list *tview.List
var server_info *tview.TextView
var now_page string
var ip string

var zmqlog *tview.TextView
var zmqloger *log.Logger

func makeMenu(app *tview.Application, pages *tview.Pages) {

	list = tview.NewList().
		AddItem("Room List", "Room info", 'r', func() {
			pages.SwitchToPage("ROOMLIST")
			now_page = "ROOMLIST"
		}).
		AddItem("Player List", "Player info", 'p', func() {
			pages.SwitchToPage("playerList")
			now_page = "playerList"
		}).
		AddItem("Close Server", "", 'q', func() {
			modal := tview.NewModal().
				SetText("確定退出？").
				AddButtons([]string{"Quit", "Cancel"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					if buttonLabel == "Quit" {
						cancel()
						app.Stop()
					} else {
						pages.RemovePage("quitModalPage")
					}
				})
			pages.AddPage("quitModalPage", modal, false, true)
		})

	list.SetFocusFunc(func() {
		textView.ScrollToEnd()
		zmqlog.ScrollToEnd()
	})
	frame := tview.NewFrame(list)
	frame.SetBorder(true).SetTitle("Menu")
	frame.SetBorders(0, 0, 1, 1, 1, 1)
	frame.AddText("Welcome to the Server", true, tview.AlignCenter, tcell.ColorWhite)

	textView = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true)
	textView.SetBorder(true).SetTitle("Lounge Log")

	// 獲取伺服器的 IP 地址
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("Error getting IP addresses: %v\n", err)
	} else {
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					ip = ipNet.IP.String()
					break
				}
			}
		}
	}

	online_player := strconv.Itoa(len(playerlist))
	server_info_str := "Server IP: " + "[yellow]" + ip + "  " + "[reset]Lounge Port: [yellow]8080" + "  " + "[reset]Room Port: [yellow]7125" + "  " + "[reset]Online Player: [yellow]" + online_player + "[reset]"

	textView.SetChangedFunc(func() {
		if now_page == "menu" {
			online_player = strconv.Itoa(len(playerlist))
			server_info_str = "Server IP: " + "[yellow]" + ip + "  " + "[reset]Lounge Port: [yellow]8080" + "  " + "[reset]Room Port: [yellow]7125" + "  " + "[reset]Online Player: [yellow]" + online_player + "[reset]"
			server_info.SetText(server_info_str)
			app.Draw()
		}
	})

	server_info = tview.NewTextView().
		SetText(server_info_str).
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true).
		SetTextAlign(tview.AlignCenter)
	server_info.SetBorder(true).SetTitle("Server Info")

	zmqlog = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true)
	zmqlog.SetBorder(true).SetTitle("ZMQ Log")

	zmqlog.SetChangedFunc(func() {
		app.Draw()
	})

	flex_infoAndLog := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(server_info, 3, 1, false).
		AddItem(textView, 0, 1, false).
		AddItem(zmqlog, 0, 1, false)

	flex := tview.NewFlex().
		AddItem(frame, 0, 1, true).
		AddItem(flex_infoAndLog, 0, 4, false)

	information := tview.NewTextView().
		SetText("[white]ESC: MENU   F9:ROOM LIST   F10: PLAYER LIST   ←↑→↓: CHOOSE   ENTER: SELECT").SetDynamicColors(true).SetRegions(true).SetWordWrap(true)
	information.SetTextStyle(tcell.StyleDefault.Bold(true))
	information.SetBackgroundColor(tcell.NewRGBColor(50, 88, 134))

	flex2 := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(flex, 0, 1, true).
			AddItem(information, 1, 1, false), 0, 1, true)

	zmqloger = log.New(zmqlog, "", log.Ldate|log.Ltime)
	pages.AddPage("menu", flex2, true, true)

}

var roomslice []int
var Roomlist_tview *tview.List

func makeRoomListPage(app *tview.Application, pages *tview.Pages) {
	log.SetOutput(textView)

	Roomlist_tview = tview.NewList()

	Roomlist_tview.ShowSecondaryText(false)

	frame := tview.NewFrame(Roomlist_tview)
	frame.SetBorder(true).SetTitle("Room List")
	frame.SetBorders(0, 0, 1, 1, 1, 1)

	roomtextview := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true)
	roomtextview.SetBorder(true).SetTitle("Room Info")

	Roomlist_tview.SetSelectedFunc(func(i int, roomtext string, _ string, _ rune) {
		room_id, _ := strconv.Atoi(roomtext[7:9])
		showroominfo(room_id, roomtextview)
	})

	Roomlist_tview.SetFocusFunc(func() {
		if len(roomslice) > 0 {
			showroominfo(roomslice[Roomlist_tview.GetCurrentItem()], roomtextview)
		}
	})

	go func() {
		for {
			select {
			case <-time.After(2 * time.Second):
				currnetIndex := Roomlist_tview.GetCurrentItem()
				if now_page == "ROOMLIST" || now_page == "playerList" {
					Roomlist_tview.Clear()
					roomslice = make([]int, 0, len(roomlist))
					for k := range roomlist {
						roomslice = append(roomslice, k)
					}
					sort.Slice(roomslice, func(i, j int) bool { return roomslice[i] < roomslice[j] })
					for _, r := range roomslice {
						title := fmt.Sprintf("Room - %02d", r)
						if roomlist[r].running {
							title += "   [green]•Running[reset]   "
						} else {
							title += "   [#DAA520]•Waiting[reset]   "
						}

						title += "[blue]" + strconv.Itoa(len(roomlist[r].Players)) + "/4"
						Roomlist_tview.AddItem(title, "", 0, nil)
					}
					if len(roomslice) > 0 {
						room_id := roomslice[currnetIndex]
						showroominfo(room_id, roomtextview)
					}
					Roomlist_tview.SetCurrentItem(currnetIndex)
					app.Draw()
				}
			case <-ctx.Done():
				return

			}

		}
	}()

	roomLog := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true)
	roomLog.SetBorder(true).SetTitle("Room Log")

	flex := tview.NewFlex().
		AddItem(frame, 30, 1, true).
		AddItem(roomtextview, 0, 1, false)

	information := tview.NewTextView().
		SetText("[white]ESC: MENU   F9:ROOM LIST   F10: PLAYER LIST   ←↑→↓: CHOOSE   ENTER: SELECT").SetDynamicColors(true).SetRegions(true).SetWordWrap(true)
	information.SetTextStyle(tcell.StyleDefault.Bold(true))
	information.SetBackgroundColor(tcell.NewRGBColor(50, 88, 134))

	flex2 := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(flex, 0, 1, true).
			AddItem(information, 1, 1, false), 0, 1, true)

	// 創建一個新的頁面並添加表格
	pages.AddPage("ROOMLIST", flex2, true, false)

}

func showroominfo(room_id int, roomtextview *tview.TextView) {
	room := roomlist[room_id]
	roomtextview.Clear()
	roomtextview.SetTitle("Room " + strconv.Itoa(room_id))
	roomtextview.SetBorderAttributes(tcell.AttrBold)
	roomtextview.SetBorderPadding(1, 0, 1, 0)
	roomtextview.SetTitleAlign(tview.AlignCenter)
	roomtextview.SetTitleColor(tcell.ColorYellow)
	roomtextview.SetTextColor(tcell.ColorWhite)
	roomtextview.SetDynamicColors(true)
	roomtextview.SetRegions(true)
	roomtextview.SetWordWrap(true)
	roomtextview.SetScrollable(true)

	roomtextview.Write([]byte("[blue]Room ID:[reset] " + strconv.Itoa(room_id) + "\n"))
	roomtextview.Write([]byte("[blue]Players:[reset] "))
	for i, p := range room.Players {
		switch i {
		case 0:
			roomtextview.Write([]byte("(東) " + p.ID + ": [yellow]" + strconv.Itoa(p.Point) + "[reset]  "))
		case 1:
			roomtextview.Write([]byte("(南) " + p.ID + ": [yellow]" + strconv.Itoa(p.Point) + "[reset]  "))
		case 2:
			roomtextview.Write([]byte("(西) " + p.ID + ": [yellow]" + strconv.Itoa(p.Point) + "[reset]  "))
		case 3:
			roomtextview.Write([]byte("(北) " + p.ID + ": [yellow]" + strconv.Itoa(p.Point) + "[reset]  "))
		}

	}
	fmt.Fprint(roomtextview, "\n\n")
	if room.running {

		switch room.wind {
		case 0:
			fmt.Fprintf(roomtextview, "[blue]局數：[reset]東%d局 %d本場\n\n", room.round, room.bunround)
		case 1:
			fmt.Fprintf(roomtextview, "[blue]局數：[reset]南%d局 %d本場\n\n", room.round, room.bunround)
		case 2:
			fmt.Fprintf(roomtextview, "[blue]局數：[reset]西%d局 %d本場\n\n", room.round, room.bunround)
		case 3:
			fmt.Fprintf(roomtextview, "[blue]局數：[reset]北%d局 %d本場\n\n", room.round, room.bunround)
		}
		fmt.Fprintf(roomtextview, "[blue]牌譜：[reset]\n%v\n", room.Cardset.Card)
		fmt.Fprintf(roomtextview, "[blue]剩餘牌數：[reset]%d\n\n", room.lastcard)
		fmt.Fprint(roomtextview, "[blue]各家手牌：[reset]\n")
		fmt.Fprintf(roomtextview, "東家：\n%v\n", room.Players[0].Ma.Card)
		fmt.Fprintf(roomtextview, "南家：\n%v\n", room.Players[1].Ma.Card)
		fmt.Fprintf(roomtextview, "西家：\n%v\n", room.Players[2].Ma.Card)
		fmt.Fprintf(roomtextview, "北家：\n%v\n", room.Players[3].Ma.Card)
	}

}

var Playerlist *tview.Table

func makePlayerListPage(app *tview.Application, pages *tview.Pages) {

	Playerlist = tview.NewTable()
	Playerlist.SetTitle("Player List")
	Playerlist.SetBorders(true)
	Playerlist.SetBorderPadding(1, 1, 1, 1)
	Playerlist.SetSelectable(true, true)

	Playerlist.SetCell(0, 0, &tview.TableCell{Text: "ID", Align: tview.AlignCenter, Color: tcell.ColorYellow, Expansion: 1})
	Playerlist.SetCell(0, 1, &tview.TableCell{Text: "IP", Align: tview.AlignCenter, Color: tcell.ColorYellow, Expansion: 1})
	Playerlist.SetCell(0, 2, &tview.TableCell{Text: "Room", Align: tview.AlignCenter, Color: tcell.ColorYellow, Expansion: 1})
	Playerlist.SetCell(0, 3, &tview.TableCell{Text: "Status", Align: tview.AlignCenter, Color: tcell.ColorYellow, Expansion: 1})

	frame := tview.NewFrame(Playerlist)
	frame.SetBorders(1, 0, 0, 0, 4, 4)
	frame.AddText("Player List", true, tview.AlignCenter, tcell.ColorWhite)
	frame.AddText("Online Player: [#18E230]"+strconv.Itoa(len(playerlist))+"[reset]", true, tview.AlignRight, tcell.ColorWhite)

	information := tview.NewTextView().
		SetText("[white]ESC: MENU   F9:ROOM LIST   F10: PLAYER LIST   ←↑→↓: CHOOSE   ENTER: SELECT").SetDynamicColors(true).SetRegions(true).SetWordWrap(true)
	information.SetTextStyle(tcell.StyleDefault.Bold(true))
	information.SetBackgroundColor(tcell.NewRGBColor(50, 88, 134))

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(frame, 0, 1, true).
			AddItem(information, 1, 1, false), 0, 1, true)

	go func() {
		for {
			select {
			case <-time.After(5 * time.Second):
				if now_page == "playerList" {
					showplaylist(Playerlist, frame)
					app.Draw()
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	Playerlist.SetFocusFunc(func() {
		showplaylist(Playerlist, frame)
	})

	selectedfunc := func(row, column int) {
		if column == 2 {
			room_id := Playerlist.GetCell(row, column).Text
			if room_id == "None" {
				return
			}
			room_id_int, _ := strconv.Atoi(room_id)

			for i := range roomslice {
				if roomslice[i] == room_id_int {
					pages.SwitchToPage("ROOMLIST")
					Roomlist_tview.SetCurrentItem(i)
					app.SetFocus(Roomlist_tview)
					now_page = "ROOMLIST"
					break
				}
			}

		}
	}

	Playerlist.SetSelectedFunc(selectedfunc)

	var lastClickTime time.Time

	Playerlist.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if event.Buttons() == tcell.Button1 {
			now := time.Now()
			if now.Sub(lastClickTime) < 500*time.Millisecond {
				// 雙擊事件處理
				row, column := Playerlist.GetSelection()
				selectedfunc(row, column)
			}
			lastClickTime = now
		}
		return action, event
	})

	Playerlist.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack))

	pages.AddPage("playerList", flex, true, false)

}

func showplaylist(Playerlist *tview.Table, frame *tview.Frame) {

	Playerlist.SetSelectable(true, true)

	playerslice := make([]string, 0, len(playerlist))
	for k := range playerlist {
		playerslice = append(playerslice, k)
	}
	sort.Strings(playerslice)

	i := 1
	for _, pp := range playerslice {

		p := playerlist[pp]
		Playerlist.SetCell(i, 0, &tview.TableCell{Text: p.ID, Align: tview.AlignCenter, Expansion: 1})
		Playerlist.SetCell(i, 1, &tview.TableCell{Text: p.conn.RemoteAddr().String(), Align: tview.AlignCenter, Expansion: 1})
		if p.Room_ID == -1 {
			Playerlist.SetCell(i, 2, &tview.TableCell{Text: "None", Align: tview.AlignCenter, Expansion: 1})
		} else {
			Playerlist.SetCell(i, 2, &tview.TableCell{Text: strconv.Itoa(p.Room_ID), Align: tview.AlignCenter, Expansion: 1})
		}
		Playerlist.SetCell(i, 3, &tview.TableCell{Text: "[#18E230]•Online [reset]", Align: tview.AlignCenter, Expansion: 1})
		i++
	}

	f, err := os.Open("playerlist.txt")
	if err != nil {
		log.Println("Error reading:", err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, _, c := r.ReadLine()
		if c == io.EOF {
			break
		}
		ID_l := strings.TrimSpace(string(line))
		_, exist := playerlist[ID_l]

		if !exist {
			Playerlist.SetCell(i, 0, &tview.TableCell{Text: ID_l, Align: tview.AlignCenter, Expansion: 1})
			Playerlist.SetCell(i, 1, &tview.TableCell{Text: "None", Align: tview.AlignCenter, Expansion: 1})
			Playerlist.SetCell(i, 2, &tview.TableCell{Text: "None", Align: tview.AlignCenter, Expansion: 1})
			Playerlist.SetCell(i, 3, &tview.TableCell{Text: "[#FF6347]•Offline[reset]", Align: tview.AlignCenter, Expansion: 1})
			i++
		}

	}

	frame.Clear()
	frame.AddText("Player List", true, tview.AlignCenter, tcell.ColorWhite)
	frame.AddText("\rOnline Player: [#18E230]"+strconv.Itoa(len(playerlist))+"[reset]", true, tview.AlignRight, tcell.ColorWhite)

}

func ter() {
	go startserver()
	app := tview.NewApplication()

	app.EnableMouse(true)
	pages := tview.NewPages()

	makeMenu(app, pages)
	makeRoomListPage(app, pages)
	makePlayerListPage(app, pages)
	now_page = "menu"

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			pages.SwitchToPage("menu")
			now_page = "menu"
			online_player := strconv.Itoa(len(playerlist))
			server_info_str := "Server IP: " + "[yellow]" + ip + "  " + "[reset]Lounge Port: [yellow]8080" + "  " + "[reset]Room Port: [yellow]7125" + "  " + "[reset]Online Player: [yellow]" + online_player + "[reset]"
			server_info.SetText(server_info_str)
			app.SetFocus(list)
		case tcell.KeyF9:
			pages.SwitchToPage("ROOMLIST")
			now_page = "ROOMLIST"
			app.SetFocus(Roomlist_tview)
		case tcell.KeyF10:
			pages.SwitchToPage("playerList")
			now_page = "playerList"
			app.SetFocus(Playerlist)
		}
		return event
	})

	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}
