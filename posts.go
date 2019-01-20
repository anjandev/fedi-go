package main

import (
	"github.com/therecipe/qt/widgets"
	"github.com/McKael/madon"
	"fmt"
)

var accountsOpts struct {
	accountID             int64
	accountUID            string
	unset                 bool     // TODO remove eventually?
	limit, keep           uint     // Limit the results
	sinceID, maxID        int64    // Query boundaries
	all                   bool     // Try to fetch all results
	onlyMedia, onlyPinned bool     // For acccount statuses
	excludeReplies        bool     // For acccount statuses
	remoteUID             string   // For account follow
	reblogs               bool     // For account follow
	acceptFR, rejectFR    bool     // For account follow_requests
	list                  bool     // For account follow_requests/reports
	accountIDs            string   // For account relationships
	statusIDs             string   // For account reports
	comment               string   // For account reports
	displayName, note     string   // For account update
	profileFields         []string // For account update
	avatar, header        string   // For account update
	defaultLanguage       string   // For account update
	defaultPrivacy        string   // For account update
	defaultSensitive      bool     // For account update
	locked, bot           bool     // For account update
	muteNotifications     bool     // For account mute
	following             bool     // For account search
}

func makePost(status madon.Status, ui_posts *widgets.QVBoxLayout, ui_replyStatus *widgets.QLabel, replyingTo *madon.Status, ui_scrollArea *widgets.QScrollArea, gClient *madon.Client, lastIDchan chan int64) (){

    together := widgets.NewQVBoxLayout()
    together.SetDirection(2)

    interactions := widgets.NewQHBoxLayout()

    reply := widgets.NewQPushButton(nil)
    reply.SetText("reply")
    reply.ConnectClicked(func(bool) {
	*replyingTo = status
	ui_replyStatus.SetText(makeContent(status).Text())
    })
    interactions.InsertWidget(0, reply, 0,0)

    star := widgets.NewQPushButton(nil)
    star.SetText("star")
    star.ConnectClicked(func(bool) {
	if star.Text() == "star" {
	    err := gClient.FavouriteStatus(status.ID)
	    if err == nil {
		star.SetText("starred!")
	    }
	} else if star.Text() == "starred!" {
	    err := gClient.UnfavouriteStatus(status.ID)
	    if err == nil {
		star.SetText("star")
	}
	}
    })
    interactions.InsertWidget(0, star, 0,0)

    repost := widgets.NewQPushButton(nil)
    repost.SetText("repost")
    repost.ConnectClicked(func(bool) {
	if repost.Text() == "repost" {
	    err := gClient.ReblogStatus(status.ID)
	    if err == nil {
		repost.SetText("reposted!")
	    }
	} else if repost.Text() == "reposted!" {
	    err := gClient.UnreblogStatus(status.ID)
	    if err == nil {
		repost.SetText("repost")
	    }
	}
    })
    interactions.InsertWidget(0, repost, 0,0)

    together.InsertLayout(0, interactions, 0)

    if len(status.MediaAttachments) > 0 {
	image := makeImage(status)
	// remove this. This wont print all statuses
	together.InsertWidget(0, image, 0,0)
    }
    post := makeContent(status)
    together.InsertWidget(0, post, 0,0)


    moreStatusDetails := widgets.NewQHBoxLayout()

    date := widgets.NewQPushButton(nil)
    date.SetText(status.CreatedAt.String())
    date.ConnectClicked(func(bool) {
	ui_scrollArea,ui_posts = deletePosts(ui_scrollArea)
    	add2FeedContext(gClient, lastIDchan, replyingTo, ui_replyStatus, ui_posts, ui_scrollArea, status.ID)
    })
    moreStatusDetails.InsertWidget(0, date, 0,0)

    displayName := widgets.NewQPushButton(nil)
    displayName.SetText(status.Account.DisplayName)
    displayName.ConnectClicked(func(bool) {
	ui_scrollArea,ui_posts = deletePosts(ui_scrollArea)
    	add2FeedAccount(gClient, lastIDchan, replyingTo, ui_replyStatus, ui_posts, ui_scrollArea, status.Account.ID)
    })
    moreStatusDetails.InsertWidget(0, displayName, 0,0)



    fullName := widgets.NewQLabel(nil,0)
    fullName.SetText(status.Account.Acct)

    moreStatusDetails.InsertWidget(0, fullName, 0,0)

    moreStatusDetails.InsertWidget(0, postAvatar(status.Account.Avatar), 0,0)
    together.InsertLayout(0, moreStatusDetails, 0)
    ui_posts.AddLayout(together, 0)
}

func add2FeedContext (gClient *madon.Client, lastIDchan chan int64, replyingTo *madon.Status, ui_replyStatus *widgets.QLabel, ui_posts *widgets.QVBoxLayout, ui_scrollArea *widgets.QScrollArea, ID int64) () {

    contexts := make(chan *madon.Context)

    go func() {
	context, err := gClient.GetStatusContext(ID)
	// TODO: add status posting algorithm
	if err != nil {
	    fmt.Println(err)
	    return
	}
	contexts <- context
    }()

    context := <- contexts

    statuses := context.Descendants
    for i := len(statuses)-1; i >= 0; i--{
	ui_posts.SetDirection(3)
	makePost(statuses[i], ui_posts, ui_replyStatus, replyingTo, ui_scrollArea, gClient, lastIDchan)
    }

    ui_posts.SetDirection(3)
    lineSeperator := widgets.NewQProgressBar(nil)
    lineSeperator.SetTextVisible(false)
    lineSeperator.SetValue(100)
    ui_posts.AddWidget(lineSeperator, 0, 0)

    stat, error := gClient.GetStatus(ID)
    if error != nil {
	fmt.Println(error)
    }

    // dereference stat
    status := *stat
    ui_posts.SetDirection(3)
    makePost(status, ui_posts, ui_replyStatus, replyingTo, ui_scrollArea, gClient, lastIDchan)
    ui_posts.SetDirection(3)

    lineSeperator2 := widgets.NewQProgressBar(nil)
    lineSeperator2.SetTextVisible(false)
    lineSeperator2.SetValue(100)

    ui_posts.AddWidget(lineSeperator2, 0, 0)

    statuses = context.Ancestors
    for i := len(statuses)-1; i >= 0; i--{
	ui_posts.SetDirection(3)
	makePost(statuses[i], ui_posts, ui_replyStatus, replyingTo, ui_scrollArea, gClient, lastIDchan)
    }



}


func add2FeedInit (gClient *madon.Client, firstID *int64, lastIDchan chan int64, replyingTo *madon.Status, ui_replyStatus *widgets.QLabel, ui_posts *widgets.QVBoxLayout, ui_scrollArea *widgets.QScrollArea, timeline string) () {
    opt := timelineOpts
    statuses := timelineGetter(gClient, opt.maxID, opt.sinceID, timeline)
    *firstID = statuses[len(statuses)-1].ID
    go func() {
	lastIDchan <- statuses[0].ID}()

    for i := 0; i < len(statuses); i++{
	ui_posts.SetDirection(2)
	makePost(statuses[i], ui_posts, ui_replyStatus, replyingTo, ui_scrollArea, gClient, lastIDchan)
    }
}

func add2FeedUpdate (gClient *madon.Client, lastIDchan chan int64, replyingTo *madon.Status, ui_replyStatus *widgets.QLabel, ui_posts *widgets.QVBoxLayout, ui_scrollArea *widgets.QScrollArea, timeline string) () {
    prevLastId := <- lastIDchan
    opt := timelineOpts
    statuses := timelineGetter(gClient, opt.maxID, prevLastId, timeline)

    if len(statuses) == 0{
	go func() {
	    lastIDchan <- prevLastId
	}()
	return
    }

    go func() {
	lastIDchan <- statuses[0].ID
    }()

    for i := len(statuses)-1; i >= 0; i--{
	ui_posts.SetDirection(3)
	makePost(statuses[i], ui_posts, ui_replyStatus, replyingTo, ui_scrollArea, gClient, lastIDchan)
    }
}


func add2FeedBack (gClient *madon.Client, firstID *int64, lastIDchan chan int64, replyingTo *madon.Status, ui_replyStatus *widgets.QLabel, ui_posts *widgets.QVBoxLayout, ui_scrollArea *widgets.QScrollArea, timeline string) () {
	// TODO: Add support for accounts

    statuses := timelineGetter(gClient, *firstID, 0, timeline)

    if len(statuses) == 0{
	return
    }

    *firstID = statuses[len(statuses)-1].ID

    for i := 0; i < len(statuses); i++{
	ui_posts.SetDirection(2)
	makePost(statuses[i], ui_posts, ui_replyStatus, replyingTo, ui_scrollArea, gClient, lastIDchan)
    }
}


func add2FeedAccount (gClient *madon.Client, lastIDchan chan int64, replyingTo *madon.Status, ui_replyStatus *widgets.QLabel, ui_posts *widgets.QVBoxLayout, ui_scrollArea *widgets.QScrollArea, ID int64) () {

    opt := timelineOpts
    var limOpts *madon.LimitParams
    limOpts = new(madon.LimitParams)
    limOpts.Limit = int(opt.limit)
    limOpts.MaxID = opt.maxID
    limOpts.SinceID = opt.sinceID

    statuses, err := gClient.GetAccountStatuses(ID, false, false, false, limOpts)
    if err != nil {
	fmt.Println(err)
    }

    interactions := widgets.NewQHBoxLayout()

    var accIDs []int64
    accIDs = append(accIDs, ID)

    currentRelationship, errRelation := gClient.GetAccountRelationships(accIDs)
    if errRelation != nil{
	fmt.Println(errRelation)
    }

    follow := widgets.NewQPushButton(nil)
    if ! currentRelationship[0].Following {
	follow.SetText("follow")
    } else {
	follow.SetText("followed!")
    }

    accOpt := accountsOpts

    followReposts := &accOpt.reblogs
    follow.ConnectClicked(func(bool) {
	if follow.Text() != "followed!" {
	    _, error := gClient.FollowAccount(ID, followReposts)
	    if error != nil{
		    widgets.QMessageBox_Information(nil, "Authentication successful", error.Error(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	    } else {
		follow.SetText("followed!")
	    }
	} else {
	    _, error := gClient.UnfollowAccount(ID)
	    if error != nil{
		    widgets.QMessageBox_Information(nil, "Authentication successful", error.Error(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	    } else {
		follow.SetText("follow")
	    }
	}
    })
    interactions.InsertWidget(0, follow, 0,0)

    muteNotification := &accOpt.muteNotifications
    mute := widgets.NewQPushButton(nil)

    if ! currentRelationship[0].Muting {
	mute.SetText("mute")
    } else {
	mute.SetText("muted!")
    }

    mute.ConnectClicked(func(bool) {
	if mute.Text() != "muted!" {
	    _, error := gClient.MuteAccount(ID, muteNotification)
	    if error != nil{
		    widgets.QMessageBox_Information(nil, "Authentication successful", error.Error(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	    } else {
		mute.SetText("muted!")
	    }
	} else {
	    _, error := gClient.UnmuteAccount(ID)
	    if error != nil{
		    widgets.QMessageBox_Information(nil, "Authentication successful", error.Error(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	    } else {
		mute.SetText("mute")
	    }
	}
    })
    interactions.InsertWidget(0, mute, 0,0)


    block := widgets.NewQPushButton(nil)

    if ! currentRelationship[0].Blocking {
	block.SetText("block")
    } else {
	block.SetText("blocked!")
    }

    block.ConnectClicked(func(bool) {
	if block.Text() != "blocked!" {
	    _, error := gClient.BlockAccount(ID)
	    if error != nil{
		widgets.QMessageBox_Information(nil, "Authentication successful", error.Error(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	    } else {
		block.SetText("blocked!")
	    }
	} else {
	    _, error := gClient.UnblockAccount(ID)
	    if error != nil{
		widgets.QMessageBox_Information(nil, "Authentication successful", error.Error(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	    } else {
		block.SetText("block")
	    }
	}
    })
    interactions.InsertWidget(0, block, 0,0)


    for i := 0; i < len(statuses); i++{
	ui_posts.SetDirection(2)
	makePost(statuses[i], ui_posts, ui_replyStatus, replyingTo, ui_scrollArea, gClient, lastIDchan)
    }

    ui_posts.InsertLayout(0, interactions, 0)
}
