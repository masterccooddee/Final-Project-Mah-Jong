package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const sideWidth = 80

type fysionLayout struct {
	top, top_bar, left_bar, right_bar, content, mingbuttonlist fyne.CanvasObject
	bottom_bar                                                 [14]fyne.CanvasObject
	//throwlist                                                  [4]fyne.CanvasObject
	dividers [5]fyne.CanvasObject
}

func NewFysionLayout(top, top_bar, left_bar, right_bar, content, mingbuttonlist fyne.CanvasObject, bottom_bar [14]fyne.CanvasObject, dividers [5]fyne.CanvasObject) fyne.Layout {
	return &fysionLayout{top: top, top_bar: top_bar, left_bar: left_bar, right_bar: right_bar, bottom_bar: bottom_bar, content: content, mingbuttonlist: mingbuttonlist, dividers: dividers}
}

func (l *fysionLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {

	topHeight := l.top.MinSize().Height
	l.top.Resize(fyne.NewSize(size.Width, topHeight))

	l.top_bar.Move(fyne.NewPos(sideWidth, topHeight))
	l.top_bar.Resize(fyne.NewSize(size.Width-sideWidth*2, sideWidth))

	l.left_bar.Move(fyne.NewPos(0, topHeight))
	l.left_bar.Resize(fyne.NewSize(size.Width, size.Height-topHeight))

	l.right_bar.Move(fyne.NewPos(size.Width-sideWidth, topHeight))
	l.right_bar.Resize(fyne.NewSize(sideWidth, size.Height-topHeight))
	for c := 0; c < 14; c++ {
		//fmt.Println("c:", c)
		if c == 0 {
			bottom_bar[0].Move(fyne.NewPos(sideWidth, GUI.Size().Height-sideWidth))
			bottom_bar[0].Resize(fyne.NewSize((GUI.Size().Width-sideWidth*2)/13, sideWidth))
		} else if c == 13-mingcardamount {
			bottom_bar[c].Move(fyne.NewPos(sideWidth+(GUI.Size().Width-sideWidth*2-150*GUI.Size().Width/1024)/13*(float32)(14-mingcardamount*3), GUI.Size().Height-sideWidth))
			bottom_bar[c].Resize(fyne.NewSize((GUI.Size().Width-sideWidth*2)/13, sideWidth))
		} else {
			bottom_bar[c].Move(fyne.NewPos(sideWidth+(GUI.Size().Width-sideWidth*2-150*GUI.Size().Width/1024)/13*((float32)(c)), GUI.Size().Height-sideWidth))
			bottom_bar[c].Resize(fyne.NewSize((GUI.Size().Width-sideWidth*2)/13, sideWidth))
		}
	}

	l.content.Move(fyne.NewPos(size.Width/3+sideWidth/3, (topHeight*2+sideWidth+size.Height)/3))
	l.content.Resize(fyne.NewSize((size.Width-sideWidth*2)/3, (size.Height-topHeight-sideWidth*2)/3))

	l.mingbuttonlist.Move(fyne.NewPos(size.Width-sideWidth-600, (topHeight*2+sideWidth+size.Height)/3*1.8-size.Height/10))
	l.mingbuttonlist.Resize(fyne.NewSize(600, 50))

	/* l.throwlist[0].Move(fyne.NewPos(size.Width/3+sideWidth/3, (topHeight*2+sideWidth+size.Height)/3*2))
	l.throwlist[0].Resize(fyne.NewSize((size.Width-sideWidth*2)/3, (size.Height-topHeight-sideWidth*2)/3))
	*/
	/* l.throwlist[1].Move(fyne.NewPos(size.Width/3+sideWidth/3, (topHeight*2+sideWidth+size.Height)/3*2))
	l.throwlist[1].Resize(fyne.NewSize((size.Width-sideWidth*2)/3, (size.Height-topHeight-sideWidth*2)/3))

	l.throwlist[2].Move(fyne.NewPos(size.Width/3+sideWidth/3, (topHeight*2+sideWidth+size.Height)/3*2))
	l.throwlist[2].Resize(fyne.NewSize((size.Width-sideWidth*2)/3, (size.Height-topHeight-sideWidth*2)/3))

	l.throwlist[3].Move(fyne.NewPos(size.Width/3+sideWidth/3, (topHeight*2+sideWidth+size.Height)/3*2))
	l.throwlist[3].Resize(fyne.NewSize((size.Width-sideWidth*2)/3, (size.Height-topHeight-sideWidth*2)/3)) */

	dividerThickness := theme.SeparatorThicknessSize()
	l.dividers[0].Move(fyne.NewPos(0, topHeight))
	l.dividers[0].Resize(fyne.NewSize(size.Width, dividerThickness))

	l.dividers[1].Move(fyne.NewPos(sideWidth, topHeight))
	l.dividers[1].Resize(fyne.NewSize(dividerThickness, size.Height-topHeight))

	l.dividers[2].Move(fyne.NewPos(size.Width-sideWidth, topHeight))
	l.dividers[2].Resize(fyne.NewSize(dividerThickness, size.Height-topHeight))

	l.dividers[3].Move(fyne.NewPos(sideWidth, topHeight+sideWidth))
	l.dividers[3].Resize(fyne.NewSize(size.Width-sideWidth*2, dividerThickness))

	l.dividers[4].Move(fyne.NewPos(sideWidth, size.Height-sideWidth))
	l.dividers[4].Resize(fyne.NewSize(size.Width-sideWidth*2, dividerThickness))

	//fmt.Println("Size:", size)
	tpapped_time = 0
	tap_item = nil
}

func (l *fysionLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	borders := fyne.NewSize(
		sideWidth*5,
		l.top.MinSize().Height,
	)
	return borders.AddWidthHeight(100, 100)

}
