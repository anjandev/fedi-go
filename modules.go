package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/webengine"
	"github.com/therecipe/qt/widgets"
       "github.com/McKael/madon"
	"strconv"
)

func makeImage(status madon.Status) (*webengine.QWebEngineView) {
    // TODO: add webm support
    image := webengine.NewQWebEngineView(nil)
    image.SetHtml(imgHTML(status), core.NewQUrl())
    image.SetMinimumHeight(325)
    return image
}

func makeContent(status madon.Status) (*widgets.QTextBrowser) {
    post := widgets.NewQTextBrowser(nil)
	contentInPost := "<p>By: " + status.Account.Username + "</p>\n" + status.Content + "\nID: " + strconv.FormatInt(status.ID, 10)

    post.SetHtml(contentInPost)

    post.SetMinimumHeight(textHeight(contentInPost))
    return post
}
