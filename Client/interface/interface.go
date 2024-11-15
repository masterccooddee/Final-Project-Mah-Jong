package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
)

func ShowError(err error, w fyne.Window) {
	dialog.ShowError(err, w)
}

var LoginSuccess = false

func UI() {
	go interconnect()

	/* context := context.WithoutCancel(context.Background())
	go func() {
		for {
			if LoginSuccess {
				dealer = zmq4.NewDealer(context, zmq4.WithID(zmq4.SocketIdentity(ID)))
				defer dealer.Close()

				err := dealer.Dial("tcp://localhost:7125")
				if err != nil {
					fmt.Println("Error connecting dealer:", err)
					return
				}
				break
			}
		}
	}() */

	a := app.New()
	a.Settings().SetTheme(myCustomTheme())

	w := a.NewWindow("Mahjong")
	w.Resize(fyne.NewSize(1024, 600))
	w.SetContent(makeGUI())
	w.CenterOnScreen()

	x := a.NewWindow("Login or Register")
	x.Resize(fyne.NewSize(300, 100))
	x.CenterOnScreen()

	x.SetContent(LORinterface(&x, &w))
	x.Show()

	// DEALER 接收消息
	/* 	go func() {
		for {
			//fmt.Println("RoomID:", RoomID)
			if RoomID != "" {
				msg, err := dealer.Recv()
				if err != nil {
					fmt.Println("Error receiving message:", err)
					break
				}
				fmt.Println("Received message:", string(msg.Frames[0]))
				msg, _ = dealer.Recv()
				var pos Position
				json.Unmarshal(msg.Frames[0], &pos)
				fmt.Println(pos.Pos)
				fmt.Println(pos.Pos[ID])
			}
		}
	}() */

	a.Run()

}
