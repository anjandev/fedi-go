package main

import (
       "github.com/McKael/madon"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"os"
	"strings"
)

type ClientStruct struct {
    Name   string `json:"Name"`
    ID     string `json:"ID"`
    Secret string `json:"Secret"`
    APIBase string `json:"APIBase"`
    InstanceURL string `json:"InstanceURL"`
    UserToken struct {
        Accesstoken string `json:"access_token"`
	CreatedAt   int64  `json:"created_at"`
	Scope       string `json:"scope"`
        Tokentype   string `json:"token_type"`
    } `json:"UserToken"` // In order to parse properly, the variable names' first char must be capitalized
}

const INSTANCE_FILE = "/home/anjan/output.json"

func setInstance(gClient *madon.Client) (err error){
// TODO: make folder fedi-go
    instance, _ := json.Marshal(gClient)
    err = ioutil.WriteFile(INSTANCE_FILE, instance, 0644)

    return err
}

func readInstance() (client ClientStruct){
    jsonFile, err := os.Open(INSTANCE_FILE)

    if err != nil {
	    panic(err)
    }

    // defer the closing of our jsonFile so that we can parse it later on
    defer jsonFile.Close()

    content, _ := ioutil.ReadAll(jsonFile)


    fmt.Println(string(content))
    dec := json.NewDecoder(strings.NewReader(string(content)))

    var clientStruct ClientStruct

    err = dec.Decode(&clientStruct);

    if err != nil {
	    panic(err)
    }


    fmt.Println(clientStruct)
    return clientStruct
}
