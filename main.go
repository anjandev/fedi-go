package main

import (
	"os"
	 // "github.com/therecipe/qt/core"
	// "github.com/therecipe/qt/uitools"
	"github.com/therecipe/qt/widgets"
       "github.com/McKael/madon"
	"fmt"
)

const APPNAME string = "fedi-go"
const APPWEBSITE string = "momi.ca"
const CONFIG_PATH string = "/home/anjan/.config/fedi-go"



func main() {
    widgets.NewQApplication(len(os.Args), os.Args)
    var mainWindow *widgets.QMainWindow
    mainWindow = widgets.NewQMainWindow(nil, 0)
    mainWindow.SetWindowTitle("fedi")

    // Rewrite auth.go From
    if _, err := os.Stat(CONFIG_PATH + "/oAuth.json"); os.IsNotExist(err) {
	NewInstance().Show()
	widgets.QApplication_Exec()
    }

    client := readInstance()
    gClient, err := madon.RestoreApp(APPNAME, client.InstanceURL, client.ID, client.Secret, nil)
    if err != nil {
	fmt.Println(err)
    }

    login(gClient).Show()

    lastIDchan := make(chan int64)
    mainActivity(gClient, lastIDchan).Show()

    widgets.QApplication_Exec()
}
