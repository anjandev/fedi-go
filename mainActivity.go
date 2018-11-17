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
	opt := timelineOpts
	prevLastId := <- lastIDchan
	statuses := timelineGetter(gClient, opt.maxID, prevLastId)
	if len(statuses) == 0{
	    go func() {
		lastIDchan <- prevLastId
	    }()
	    return
	}
	fmt.Println("PREVIOUS ID")
	fmt.Println(prevLastId)

	go func() {
	lastIDchan <- statuses[0].ID
	}()

	oldStatusStillIn := false
	var oldStatusIdx int

	for i := 0; i < len(statuses); i++{
	    if (statuses[i].ID == prevLastId){
		oldStatusStillIn = true
		oldStatusIdx = i
		break
	    }
	}

	for i := 0; i < len(statuses); i++{
		fmt.Println(statuses[i].ID)
	}

	// TODO: I repeat this code 3 times only changing the initial i. put into function

	if oldStatusStillIn {
		for i := 0; i < oldStatusIdx; i++ {
		    makePost(statuses[i], *ui_posts, gClient)
		}
	} else if !oldStatusStillIn {
		fmt.Println("OLD STATUS NOT IN")
		for i := len(statuses)-1; i >= 0; i-- {
		    makePost(statuses[i], *ui_posts, gClient)
		}
	}
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

    add2Feed(gClient, lastIDchan, ui_posts)

    widget.SetWindowTitle("Fedi-go")

    return widget
}
