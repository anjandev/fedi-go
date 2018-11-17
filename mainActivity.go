package main

import (
    "github.com/therecipe/qt/core"
    "github.com/therecipe/qt/uitools"
    "github.com/therecipe/qt/widgets"
    "github.com/McKael/madon"
    "fmt"
)

// Direction to put new items in posts
// fix this. magic numbers
const BOTTOM_TO_TOP int = 3
const TOP_TO_BOTTOM int = 2

func mainActivity(gClient *madon.Client, lastIDchan chan int64) (*widgets.QWidget) {

    var widget = widgets.NewQWidget(nil, 0)

    var loader = uitools.NewQUiLoader(nil)
    var file = core.NewQFile2(":/qml/main.ui")

    file.Open(core.QIODevice__ReadOnly)
    var formWidget = loader.Load(file, widget)
    file.Close()

	var (
	    ui_posts = widgets.NewQVBoxLayoutFromPointer(widget.FindChild("posts", core.Qt__FindChildrenRecursively).Pointer())
	    ui_sendMsg = widgets.NewQPushButtonFromPointer(widget.FindChild("pushButtonPostMsg", core.Qt__FindChildrenRecursively).Pointer())
	    ui_updateFeed = widgets.NewQPushButtonFromPointer(widget.FindChild("updateFeed", core.Qt__FindChildrenRecursively).Pointer())
	    ui_msg = widgets.NewQTextEditFromPointer(widget.FindChild("postMsg", core.Qt__FindChildrenRecursively).Pointer())
	    ui_timelineSelector = widgets.NewQComboBoxFromPointer(widget.FindChild("timelineSelector", core.Qt__FindChildrenRecursively).Pointer())
	)

	
	ui_timelineSelector.ConnectActivated(func(index int) {
		fmt.Println("lol")
	})

    ui_updateFeed.ConnectClicked(func(bool) {
	add2Feed(gClient,lastIDchan, ui_posts, false)
    })


//	for i:= 0; i < len(updates); i++{
//		updates.Data
//	}

    // Sending post handler
    ui_sendMsg.ConnectClicked(func(bool) {
	var newPost madon.PostStatusParams
	newPost.Text = ui_msg.ToPlainText()
	_, error := gClient.PostStatus(newPost)
	if error != nil {
	    fmt.Println(error)
	} else {
		ui_msg.Clear()
	}
    })


    // Fill first open with content :3
    var layout = widgets.NewQVBoxLayout()
    layout.AddWidget(formWidget, 0, 0)
    widget.SetLayout(layout)

    add2Feed(gClient, lastIDchan, ui_posts, true)

    widget.SetWindowTitle("Fedi-go")

    return widget
}
