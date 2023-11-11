package graphic

import (
	"Copy-Teleport/devices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	
	"strconv"
)

func GetDevicesContainer() *fyne.Container {

	devicesList = widget.NewList(
		func() int {
			return len(devices.Values)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(strconv.Itoa(i+1) + ") " + devices.Values[i].Username + ": " + devices.Values[i].Ip_address)
		})



	clearButton := widget.NewButton("Clear", func () {
		devices.Values = make([]devices.DevicesElement, 0)
		RefreshHistoryItems()
	})


	deviceLabel := widget.NewLabel("Devices")
	deviceLabel.TextStyle = fyne.TextStyle{Bold: true}

	out := container.NewBorder(
		container.NewCenter(deviceLabel),
		clearButton,
		nil,
		nil,
		devicesList,

	)
	return out

}