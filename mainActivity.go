package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/uitools"
	"github.com/therecipe/qt/widgets"
       "github.com/McKael/madon"
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
    widget.SetLayout(layout)

    for i := 0; i < 20; i++{
	post := widgets.NewQTextBrowser(nil)
	post.SetText("By: " + statuses[i].Account.Username + statuses[i].Content)

	ui_posts.InsertWidget(i, post, 0,0)
    }

    widget.SetWindowTitle("Fedi-go")

    return widget
}
