package main

import (
	"github.com/therecipe/qt/widgets"
    "github.com/McKael/madon"
   "fmt"
)

func makePost(status madon.Status, ui_posts widgets.QVBoxLayout, gClient *madon.Client) (widgets.QVBoxLayout){
    ui_posts.SetDirection(2)

    interactions := widgets.NewQHBoxLayout()

    reply := widgets.NewQPushButton(nil)
    reply.SetText("reply")
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

    ui_posts.InsertLayout(0, interactions, 0)

    if len(status.MediaAttachments) > 0 {
	image := makeImage(status)
	// remove this. This wont print all statuses
	ui_posts.InsertWidget(0, image, 0,0)
    }
    post := makeContent(status)
    ui_posts.InsertWidget(0, post, 0,0)
    return ui_posts
}

func add2Feed (gClient *madon.Client, lastIDchan chan int64, ui_posts *widgets.QVBoxLayout) () {
    opt := timelineOpts
    statuses := timelineGetter(gClient, opt.maxID, opt.sinceID)

    fmt.Println(len(statuses))
    for i := len(statuses)-1; i >= 0; i--{
	makePost(statuses[i], *ui_posts, gClient)
    }

	go func() {
		lastIDchan <- statuses[0].ID}()
}
