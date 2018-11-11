package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/uitools"
	"github.com/therecipe/qt/widgets"
       "github.com/McKael/madon"
	"fmt"
)

// TODO: WRITE IT TO PARAM TO DISK

func submitOAuth(gClient *madon.Client) *widgets.QWidget {
	var widget = widgets.NewQWidget(nil, 0)

	var loader = uitools.NewQUiLoader(nil)
	var file = core.NewQFile2(":/qml/submitOAuth.ui")

	file.Open(core.QIODevice__ReadOnly)
	var formWidget = loader.Load(file, widget)
	file.Close()

	var (
		ui_oAuthTextInput = widgets.NewQLineEditFromPointer(widget.FindChild("oAuth", core.Qt__FindChildrenRecursively).Pointer())
		ui_auth = widgets.NewQPushButtonFromPointer(widget.FindChild("submitOAuth", core.Qt__FindChildrenRecursively).Pointer())
	)

	ui_auth.ConnectClicked(func(_ bool) {
		token := ui_oAuthTextInput.Text()
		if token != "" {
		   // err outputted here
		   err := oAuth2ExchangeCode(token, gClient)
		   fmt.Println(err)
		
		} else if token == "" {
			widgets.QMessageBox_Information(nil, "I have logged in on this address", "Your token cannot be blank", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		}
	})

	var layout = widgets.NewQVBoxLayout()
	layout.AddWidget(formWidget, 0, 0)
	widget.SetLayout(layout)

	widget.SetWindowTitle("Enter Instance")

	return widget
}
