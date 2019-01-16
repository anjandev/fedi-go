package main

import (
    "github.com/therecipe/qt/core"
    "github.com/therecipe/qt/uitools"
    "github.com/therecipe/qt/widgets"
    "github.com/McKael/madon"
    "fmt"
    "time"
)

// Direction to put new items in posts
// fix this. magic numbers
const BOTTOM_TO_TOP int = 3
const TOP_TO_BOTTOM int = 2


type QVBoxLayoutCustomSlot struct {
	widgets.QVBoxLayout

	_ func(notification madon.Notification) `slot:"triggerSlot"`
}


func mainActivity(gClient *madon.Client, lastIDchan chan int64) (*widgets.QWidget) {

    var widget = widgets.NewQWidget(nil, 0)

    var loader = uitools.NewQUiLoader(nil)
    var file = core.NewQFile2(":/qml/main.ui")

    file.Open(core.QIODevice__ReadOnly)
    var formWidget = loader.Load(file, widget)
    file.Close()

    var (
	ui_posts = widgets.NewQVBoxLayoutFromPointer(widget.FindChild("posts", core.Qt__FindChildrenRecursively).Pointer())
	ui_leftPanel = widgets.NewQWidgetFromPointer(widget.FindChild("scrollAreaWidgetContents_2", core.Qt__FindChildrenRecursively).Pointer())
	ui_sendMsg = widgets.NewQPushButtonFromPointer(widget.FindChild("pushButtonPostMsg", core.Qt__FindChildrenRecursively).Pointer())
	ui_updateFeed = widgets.NewQPushButtonFromPointer(widget.FindChild("updateFeed", core.Qt__FindChildrenRecursively).Pointer())
	ui_msg = widgets.NewQTextEditFromPointer(widget.FindChild("postMsg", core.Qt__FindChildrenRecursively).Pointer())
	ui_timelineSelector = widgets.NewQComboBoxFromPointer(widget.FindChild("timelineSelector", core.Qt__FindChildrenRecursively).Pointer())
	// ui_scrollAreaContent = widgets.NewQWidgetFromPointer(widget.FindChild("scrollAreaWidgetContents", core.Qt__FindChildrenRecursively).Pointer())
	ui_scrollArea = widgets.NewQScrollAreaFromPointer(widget.FindChild("scrollArea", core.Qt__FindChildrenRecursively).Pointer())
	ui_replyStatus = widgets.NewQLabelFromPointer(widget.FindChild("replyingStatus", core.Qt__FindChildrenRecursively).Pointer())
    )

    var replyingTo madon.Status

    ui_timelineSelector.ConnectActivated(func(index int) {
	// clear channel of old Ids if user has changed the timeline
	if ui_posts.Count() > 0 {
	    fmt.Println(<- lastIDchan)
	}
	// delete all children
//	for i:= 0; i < ui_posts.Count(); i++ {
//	    child := ui_posts.ItemAt(i)
//	    ui_posts.RemoveItem(child)
//	}
	// ui_scrollArea.TakeWidget()
	// add2Feed(gClient,lastIDchan, ui_posts, true, ui_timelineSelector.CurrentText())

	// put this in a function. Will need to for replying
	ui_scrollArea, ui_posts = deletePosts(ui_scrollArea)
	add2Feed(gClient, lastIDchan, &replyingTo, ui_replyStatus, ui_posts, ui_scrollArea, "initialize", ui_timelineSelector.CurrentText(), 0)
    })

    ui_updateFeed.ConnectClicked(func(bool) {
	add2Feed(gClient, lastIDchan, &replyingTo, ui_replyStatus, ui_posts, ui_scrollArea, "updatefeed", ui_timelineSelector.CurrentText(), 0)
    })

    // Sending post handler
    ui_sendMsg.ConnectClicked(func(bool) {
	var newPost madon.PostStatusParams
	newPost.Text = ui_msg.ToPlainText()
	newPost.InReplyTo = replyingTo.ID
	_, error := gClient.PostStatus(newPost)
	if error != nil {
	    fmt.Println(error)
	} else {
	    ui_msg.Clear()
	    ui_replyStatus.SetText("")
	    var newReply madon.Status
	    replyingTo = newReply
	}
    })

    var layout = widgets.NewQVBoxLayout()
    layout.AddWidget(formWidget, 0, 0)
    widget.SetLayout(layout)


    // Fill first open with content :3
    // this should only happen once (in the beginning)
    add2Feed(gClient, lastIDchan, &replyingTo, ui_replyStatus, ui_posts, ui_scrollArea, "initialize", ui_timelineSelector.CurrentText(), 0)

    widget.SetWindowTitle("Fedi-go")

    var ui_notif = NewQVBoxLayoutCustomSlot()

    ui_leftPanel.SetLayout(ui_notif)

    ui_notif.ConnectTriggerSlot(func(notification madon.Notification) {
	if notification.Status != nil {
	    ui_status := widgets.NewQLabel(nil,0)
	    ui_status.SetText(notification.Status.Content)
	    ui_status.SetWordWrap(true)
	    ui_notif.InsertWidget(0, ui_status, 0,0)
	}

	ui_account := widgets.NewQLabel(nil,0)
	ui_account.SetText(notification.Account.Acct)
	ui_account.SetWordWrap(true)
	ui_notif.InsertWidget(0, ui_account, 0,0)

	ui_type := widgets.NewQLabel(nil,0)
	ui_type.SetText(notification.Type)
	ui_type.SetWordWrap(true)
	ui_notif.InsertWidget(0, ui_type, 0,0)
    })

    go func() {
	var lastNotifID int64 = 0

	for {
	    var lopt *madon.LimitParams
	    lopt = new(madon.LimitParams)
	    lopt.All = true
	    notifications, err := gClient.GetNotifications(nil, lopt)

	    if err != nil {
		fmt.Println("Error getting notifications")
		fmt.Println(err)
		continue
	    }


	    for i:= 0; i < 6; i++ {
		fmt.Println(i)
		fmt.Println(notifications[i].ID)
	     	if notifications[i].ID == lastNotifID {
	     		break
	     	}
		ui_notif.TriggerSlot(notifications[i])
	    }

	     lastNotifID = notifications[0].ID

	    time.Sleep(2 * time.Second)
	}
    }()

    return widget
}
