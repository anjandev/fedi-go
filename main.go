package main

import (
	"os"
	 // "github.com/therecipe/qt/core"
	// "github.com/therecipe/qt/uitools"
	"github.com/therecipe/qt/widgets"
        "github.com/McKael/madon"
	"fmt"
	"github.com/therecipe/qt/core"
)

const APPNAME string = "fedi-go"
const APPWEBSITE string = "momi.ca"
var CONFIG_PATH string = (core.QDir_HomePath() + "/.config/fedi-go")
var CONFIG_PATH_OAUTH string = CONFIG_PATH + "/OAuth.json"



func main() {
    widgets.NewQApplication(len(os.Args), os.Args)
    var mainWindow *widgets.QMainWindow
    mainWindow = widgets.NewQMainWindow(nil, 0)
    mainWindow.SetWindowTitle("fedi")

    if _, err := os.Stat(CONFIG_PATH); err != nil {
	if os.IsNotExist(err) {
		os.Mkdir(CONFIG_PATH, os.FileMode(0700))
	}
    }

    // Rewrite auth.go From
    if _, err := os.Stat(CONFIG_PATH_OAUTH); os.IsNotExist(err) {
	NewInstance().Show()
	widgets.QApplication_Exec()
    }


    // fix it so it checks if basic auth
    client := readInstance()
    gClient, err := madon.RestoreApp(APPNAME, client.InstanceURL, client.ID, client.Secret, nil)
    if err != nil {
    	fmt.Println(err)
    }
    

    login(gClient).Show()

    widgets.QApplication_Exec()
}
