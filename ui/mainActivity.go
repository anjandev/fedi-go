package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/webengine"
	"github.com/therecipe/qt/uitools"
	"github.com/therecipe/qt/widgets"
       "github.com/McKael/madon"
	"regexp"
	"fmt"
)

func mainActivity(statuses []madon.Status) *widgets.QWidget {

    var widget = widgets.NewQWidget(nil, 0)

    var loader = uitools.NewQUiLoader(nil)
    var file = core.NewQFile2(":/qml/main.ui")

    file.Open(core.QIODevice__ReadOnly)
    var formWidget = loader.Load(file, widget)
    file.Close()

	var (
		ui_posts = widgets.NewQVBoxLayoutFromPointer(widget.FindChild("posts", core.Qt__FindChildrenRecursively).Pointer())
	)

    var layout = widgets.NewQVBoxLayout()
    layout.AddWidget(formWidget, 0, 0)
//	layout.AddStretch(0)
    widget.SetLayout(layout)

    for i := 0; i < len(statuses); i++{
	post := widgets.NewQTextBrowser(nil)
	contentInPost := "<p>By: " + statuses[i].Account.Username + "</p>\n" + statuses[i].Content

	post.SetHtml(contentInPost)

	// find number of lines to properly size box
	// I really dont like this because if someone posts a really long 1 liner, this algorithm
	// fails. Try something else. Check line size?
	paras:= regexp.MustCompile("</p>")
	matches := paras.FindAllStringIndex(contentInPost, -1)
	const HEIGHT_OF_LINE = 25
	fmt.Println(len(matches))
	    post.SetMinimumHeight((len(matches)+2) * HEIGHT_OF_LINE)
	    // minimum
	// post.SetSizePolicy2(1, 0)
	    // post.AdjustSize()
	    // adjust to contents
	// post.SetSizeAdjustPolicy(2)

	ui_posts.InsertWidget(i, post, 0,0)

	if len(statuses[i].MediaAttachments) > 0 {
	    image := webengine.NewQWebEngineView(nil)
	    fmt.Println(statuses[i].MediaAttachments[0].URL)
	    image.SetHtml("<!DOCTYPE html> <html> <body> <img src=" + statuses[i].MediaAttachments[0].URL + "> </body> <html> ", core.NewQUrl())
		image.SetMinimumHeight(100)
	    // remove this. This wont print all statuses
	    i = i +1
	    ui_posts.InsertWidget(i, image, 0,0)
	}
//	for i2 := 0; i2 < len(statuses[i].MediaAttachments); i++ {
//	    img := widgets.NewQWebEngineView(nil)
//	    url = statuses[i].MediaAttachments[i].URL
//	    webEngineView
//	}
    }

    widget.SetWindowTitle("Fedi-go")

    return widget
}
