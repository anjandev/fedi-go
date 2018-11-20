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

func makeContent(status madon.Status) (*widgets.QLabel) {
    contentInPost := "<p>By: " + status.Account.Username + "</p>\n" + status.Content + "\nID: " + strconv.FormatInt(status.ID, 10)
    post := widgets.NewQLabel(nil, 0)
    post.SetText(contentInPost)

    post.SetWordWrap(true)
    post.SetOpenExternalLinks(true)
    // set text selectable by mouse: https://doc.qt.io/qt-5/qt.html#TextInteractionFlag-enum
    post.SetTextInteractionFlags(1)

    return post
}
