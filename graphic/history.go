package graphic


import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	// "fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"Copy-Teleport/history"

	// "fyne.io/fyne/v2/canvas"
	// "image/color"

	 // "reflect"
	 // "fmt"
	 // "crypto/rand"
	 // "math/big"

)

var historyList *widget.List


func GetHistoryContainer() *fyne.Container{

	// Add("pippo", "ciao")	// test
	// Add("pippo1", "fffff")	// test

	historyList = widget.NewList(
		func() int {
			return len(history.Values)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText("Copy from " + history.Values[i].Username + ": " + history.Values[i].Value)
		})



	clearButton := widget.NewButton("Clear", func () {
		history.Values = make([]history.HistoryElement, 0)
		RefreshHistoryItems()
	})

	hystoryLabel := widget.NewLabel("History")
	hystoryLabel.TextStyle = fyne.TextStyle{Bold: true}

	out := container.NewBorder(
		container.NewCenter(hystoryLabel),
		clearButton,
		nil,
		nil,
		historyList,

	)
	return out
}



func RefreshHistoryItems(){
	historyList.Refresh()
}