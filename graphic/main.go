package graphic


import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	// "fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2"

	// "fmt"

	// "Copy-Teleport/connect"
	// "Copy-Teleport/history"
)

var window fyne.Window

func GetGUI() {
	myApp := app.New()
	window = myApp.NewWindow("Copy-Teleport")

	connectContainer := GetConnectContainer()
	historyContainer := GetHistoryContainer()
	devicesContainer := GetDevicesContainer()

	tabs := container.NewAppTabs(
		container.NewTabItem("Connect", connectContainer),
		container.NewTabItem("Devices", devicesContainer),
		container.NewTabItem("History", historyContainer),
		container.NewTabItem("Settings", widget.NewLabel("Setting work in progress!!!")),
	)

	tabs.SetTabLocation(container.TabLocationLeading)
	window.SetContent(tabs)
	window.Resize(fyne.NewSize(600, 400))
	window.ShowAndRun()

}