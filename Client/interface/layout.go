package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const sideWidth = 100

type fysionLayout struct {
	top, top_bar, left_bar, right_bar, bottom_bar, content fyne.CanvasObject
	dividers                                               [5]fyne.CanvasObject
}

func NewFysionLayout(top, top_bar, left_bar, right_bar, bottom_bar, content fyne.CanvasObject, dividers [5]fyne.CanvasObject) fyne.Layout {
	return &fysionLayout{top: top, top_bar: top_bar, left_bar: left_bar, right_bar: right_bar, bottom_bar: bottom_bar, content: content, dividers: dividers}
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

	l.bottom_bar.Move(fyne.NewPos(sideWidth, size.Height-sideWidth))
	l.bottom_bar.Resize(fyne.NewSize(size.Width-sideWidth*2, sideWidth))

	l.content.Move(fyne.NewPos(sideWidth, topHeight+sideWidth))
	l.content.Resize(fyne.NewSize(size.Width-sideWidth*2, size.Height-topHeight-sideWidth*2))

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

	//log.Println("Size:", size)
}

func (l *fysionLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	borders := fyne.NewSize(
		sideWidth*5,
		l.top.MinSize().Height,
	)
	return borders.AddWidthHeight(100, 100)

}