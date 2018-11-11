package main

import (
       "github.com/McKael/madon"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

type ClientStruct struct {
    Name   string `json:"Name"`
    ID     string `json:"ID"`
    Secret string `json:"Secret"`
    APIBase string `json:"APIBase"`
    InstanceURL string `json:"InstanceURL"`
    UserToken []struct {
        access_token string `json:"access_token"`
        token_type   string `json:"token_type"`
    } `json:"UserToken"`
}

const INSTANCE_FILE = "/home/anjan/output.json"

func setInstance(gClient *madon.Client) (err error){
// TODO: make folder fedi-go
    instance, _ := json.Marshal(gClient)
    err = ioutil.WriteFile(INSTANCE_FILE, instance, 0644)

    return err
}

func readInstance() (client ClientStruct){
    content, err := ioutil.ReadFile(INSTANCE_FILE)
    str := string(content)

    fmt.Println(str)

    if err != nil {
	    panic(err)
    }

    var clientStruct ClientStruct

    err = json.Unmarshal([]byte(str), &clientStruct)

//    if err != nil {
//	    panic(err)
//    }


    return clientStruct
}
