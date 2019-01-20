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

    var firstID int64 = 0
    var replyingTo madon.Status

    ui_scrollArea.VerticalScrollBar().ConnectValueChanged(func(value int) {
	if ui_scrollArea.VerticalScrollBar().Value() == ui_scrollArea.VerticalScrollBar().Maximum() {
	    add2FeedBack (gClient, &firstID, lastIDchan, &replyingTo, ui_replyStatus, ui_posts, ui_scrollArea, ui_timelineSelector.CurrentText())
	}
    })


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
	add2FeedInit(gClient, &firstID, lastIDchan, &replyingTo, ui_replyStatus, ui_posts, ui_scrollArea, ui_timelineSelector.CurrentText())
    })

    ui_updateFeed.ConnectClicked(func(bool) {
	add2FeedUpdate(gClient, lastIDchan, &replyingTo, ui_replyStatus, ui_posts, ui_scrollArea, ui_timelineSelector.CurrentText())
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
    add2FeedInit(gClient, &firstID, lastIDchan, &replyingTo, ui_replyStatus, ui_posts, ui_scrollArea, ui_timelineSelector.CurrentText())

    widget.SetWindowTitle("Fedi-go")

    var ui_notif = NewQVBoxLayoutCustomSlot()

    ui_leftPanel.SetLayout(ui_notif)

    ui_notif.ConnectTriggerSlot(func(notification madon.Notification) {
	var card = widgets.NewQHBoxLayout()

	var details = widgets.NewQVBoxLayout()

	if notification.Status != nil {
	    ui_status := widgets.NewQLabel(nil,0)
	    ui_status.SetText("<b>" + notification.Status.Content + "</b>")
	    ui_status.SetWordWrap(true)
	    details.InsertWidget(0, ui_status, 0,0)
	}

	ui_account := widgets.NewQLabel(nil,0)
	ui_account.SetText(notification.Account.Acct)
	ui_account.SetWordWrap(true)
	details.InsertWidget(0, ui_account, 0,0)

	ui_type := widgets.NewQLabel(nil,0)
	ui_type.SetText(notification.Type)
	ui_type.SetWordWrap(true)
	details.InsertWidget(0, ui_type, 0,0)

        card.InsertLayout(0, details, 0)
        card.InsertWidget(0, postAvatar(notification.Account.Avatar), 0,0)

	ui_notif.InsertLayout(0, card, 0)
    })

    go func() {
	var lastNotifID int64 = 0

	for {
	    var lopt *madon.LimitParams
	    lopt = new(madon.LimitParams)
	    // initially only get 10 notifications
	    // but after initialization, get as many notifications as possible with each poll
	    lopt.SinceID = lastNotifID

	    if lastNotifID == 0 {
		lopt.All = false
		lopt.Limit = 10
	    }

	    notifications, err := gClient.GetNotifications(nil, lopt)

	    if err != nil {
		fmt.Println("Error getting notifications")
		fmt.Println(err)
		continue
	    }

	    if len(notifications) == 0 {
		continue
	    }

	    for i:= len(notifications)-1; i >= 0; i-- {
		ui_notif.TriggerSlot(notifications[i])
	    }

	    lastNotifID = notifications[0].ID

	    time.Sleep(time.Minute * 5)
	}
    }()

    return widget
}
