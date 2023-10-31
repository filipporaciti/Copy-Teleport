package graphic


import (
	"time"
	"strings"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"Copy-Teleport/connect"
)



var usernameEntry = widget.NewEntry()
var passwordEntry = widget.NewPasswordEntry()
var devicesList *widget.List



func GetConnectContainer() *fyne.Container {

	usernameEntry.OnChanged = func(text string) {
		connect.Username = text
	}
	passwordEntry.OnChanged = func(text string) {
		connect.Password  = text
	}


	devicesList = widget.NewList(
		func() int {
			return len(connect.Values)
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			input := widget.NewPasswordEntry()
			button := widget.NewButton("Connect", nil)
			// return container.NewHBox(input, button, label)
			return container.NewBorder(
				nil, 
				nil,
				label, 
				button,
				input,
			)
		},
		func(index int, obj fyne.CanvasObject) {
			input := obj.(*fyne.Container).Objects[0].(*widget.Entry)
			label := obj.(*fyne.Container).Objects[1].(*widget.Label)
			button := obj.(*fyne.Container).Objects[2].(*widget.Button)

			label.SetText(connect.Values[index].Username + ": " + connect.Values[index].Ip_address)
			button.OnTapped = func() {
				if data, _ := connect.SendConnectionRequest(connect.Values[index].Ip_address, input.Text); !data {
					SpawnPopUp("Connection falied")
				}
			}
        },
	)
	refreshButton := widget.NewButton("Refresh", func () {

		// qui aggiorno la device List
		if data, err := connect.DiscoverDevices(); !data && !(strings.Contains(err.Error(), "i/o timeout") || strings.Contains(err.Error(), "host is down") || strings.Contains(err.Error(), "no route to host")) {
			fmt.Println("----",err.Error())
			SpawnPopUp("Discover error")
		} else {
			RefreshConnectItems()
		}
	})

	manualIpButton := widget.NewButton("Add device", func() {
		var popup *widget.PopUp

		entry := widget.NewEntry()
		entry.SetPlaceHolder("Ip address")

		button := widget.NewLabel("Add device")
		button.TextStyle = fyne.TextStyle{Bold: true}

		errorLabel := widget.NewLabel("")

		content := container.NewBorder(
	            container.NewVBox(
	            	container.NewCenter(button),
		            entry,
		            errorLabel,
		        ),
	            widget.NewButton("Add", func() {
	            	b, err := connect.SendOneBeaconRequest(entry.Text)
	            	if !b {
	            		errorLabel.SetText("Errore: " + err.Error())
	            		time.Sleep(time.Millisecond * 3000)
	            	}
	            	popup.Hide()
	            }),
	            nil,
	            nil,
	    )

	    popup = widget.NewModalPopUp(content, window.Canvas())
	    popup.Resize(fyne.NewSize(200, 150))
	    popup.Show()
	})
	

	usernameContainer := container.NewBorder(
			nil,
			nil,
			container.New(
				layout.NewHBoxLayout(),
				widget.NewLabel(""),
				widget.NewLabel("Username\t"),
			),
			widget.NewLabel("   "),
			usernameEntry,
		)


	passwordContainer := container.NewBorder(
			nil,
			nil,
			container.New(
				layout.NewHBoxLayout(),
				widget.NewLabel(""),
				widget.NewLabel("Password\t"),
			),
			
			container.New(
				layout.NewHBoxLayout(),
				widget.NewButton("Random", func() {
					connect.SetRandomPassword(10)
					RefreshConnectItems()
					}),
				widget.NewLabel("   "),
			),
			passwordEntry,
		)

	aviableLabel := widget.NewLabel("Avaiable devices")
	aviableLabel.TextStyle = fyne.TextStyle{Bold: true}


	out := container.NewBorder(

		container.New(
			layout.NewVBoxLayout(),
			usernameContainer,
			passwordContainer,
			widget.NewLabel(""),
			container.NewCenter(aviableLabel),
		),
		container.NewBorder(
			nil,
			nil,
			nil,
			manualIpButton,
			refreshButton,
		),
		nil,
		nil,
		devicesList,

	)


	connect.SetRandomUsername()
	connect.SetRandomPassword(10)
	RefreshConnectItems()
	return out
	
}

func SpawnPopUp(labelName string) {
	var popup *widget.PopUp
	content := container.NewVBox(
            widget.NewLabel(labelName),
            widget.NewButton("Ok", func() {
            	popup.Hide()
            }),
    )
    popup = widget.NewModalPopUp(content, window.Canvas())
    popup.Show()
}

func RefreshConnectItems() {
	usernameEntry.SetText(connect.Username)
	passwordEntry.SetText(connect.Password)

	usernameEntry.Refresh()
	passwordEntry.Refresh()
	devicesList.Refresh()
}