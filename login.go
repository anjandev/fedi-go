package main

import (
    "github.com/therecipe/qt/core"
    "github.com/therecipe/qt/uitools"
    "github.com/therecipe/qt/widgets"
    "github.com/McKael/madon"
    "fmt"
)

// you do not need to turn an object from a function to change its state
func login(gClient *madon.Client) (*widgets.QWidget) {
	var widget = widgets.NewQWidget(nil, 0)

	var loader = uitools.NewQUiLoader(nil)
	var file = core.NewQFile2(":/qml/login.ui")

	file.Open(core.QIODevice__ReadOnly)
	var formWidget = loader.Load(file, widget)
	file.Close()

	var (
		ui_username = widgets.NewQLineEditFromPointer(widget.FindChild("username", core.Qt__FindChildrenRecursively).Pointer())
		ui_password = widgets.NewQLineEditFromPointer(widget.FindChild("password", core.Qt__FindChildrenRecursively).Pointer())
		ui_submit = widgets.NewQPushButtonFromPointer(widget.FindChild("submit", core.Qt__FindChildrenRecursively).Pointer())
	)

	ui_submit.ConnectClicked(func(_ bool) {
	    var scopes = []string{"read", "write", "follow"}

	    err := gClient.LoginBasic(ui_username.Text(), ui_password.Text(), scopes)
	    if err != nil {
		fmt.Println(err)
		// TODO: make this error message more detailed
		widgets.QMessageBox_Information(nil, "An error ocurred. See terminal output", "An error ocurred. See terminal output", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	    } else {
		// TODO: Save to file
		widgets.QMessageBox_Information(nil, "Authentication successful", "Authentication successful", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		widget.Close()
		lastIDchan := make(chan int64)
		mainActivity(gClient, lastIDchan).Show()
	    }
	})

	var layout = widgets.NewQVBoxLayout()
	layout.AddWidget(formWidget, 0, 0)
	widget.SetLayout(layout)

	widget.SetWindowTitle("Login to the fediverse")

	return widget
}
