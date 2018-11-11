package main

import (
	"os"
	 // "github.com/therecipe/qt/core"
	// "github.com/therecipe/qt/uitools"
	"github.com/therecipe/qt/widgets"
)

const APPNAME string = "fedi-go"
const APPWEBSITE string = "momi.ca"
const CONFIG_PATH string = "~/.config/fedi-go"

func main() {
    widgets.NewQApplication(len(os.Args), os.Args)


    // add error checking if instance does not exist
    // NewInstance().Show()
    _ = readInstance()


    widgets.QApplication_Exec()
}
