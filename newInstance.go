package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/uitools"
	"github.com/therecipe/qt/widgets"
)

func NewInstance() *widgets.QWidget {
	var widget = widgets.NewQWidget(nil, 0)

	var loader = uitools.NewQUiLoader(nil)
	var file = core.NewQFile2(":/qml/calculatorform.ui")

	file.Open(core.QIODevice__ReadOnly)
	var formWidget = loader.Load(file, widget)
	file.Close()

	var (
		ui_urlInput = widgets.NewQLineEditFromPointer(widget.FindChild("url", core.Qt__FindChildrenRecursively).Pointer())
		ui_auth = widgets.NewQPushButtonFromPointer(widget.FindChild("pushButton", core.Qt__FindChildrenRecursively).Pointer())
		ui_basic = widgets.NewQCheckBoxFromPointer(widget.FindChild("checkBox", core.Qt__FindChildrenRecursively).Pointer())
	)

	ui_auth.ConnectClicked(func(_ bool) {
		gClient := getClient(ui_urlInput.Text())
		if !(ui_basic.IsChecked()) {
		    url := getAuthOAuth(ui_urlInput.Text(), gClient)
		    if url == "" {
			widgets.QMessageBox_Information(nil, "I have logged in on this address", "This is not a valid instance", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		    } else {
			widgets.QMessageBox_Information(nil, "I have logged in on this address", url, widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
			submitOAuth(gClient).Show()
			widget.Close()
		    }
		} else {
		    login(gClient).Show()
		    widget.Close()
		}
	})

	var layout = widgets.NewQVBoxLayout()
	layout.AddWidget(formWidget, 0, 0)
	widget.SetLayout(layout)

	widget.SetWindowTitle("Enter Instance")

	return widget
}
