package main

import (
       "github.com/McKael/madon"
	"encoding/json"
	"io/ioutil"
)


func setInstance(gClient *madon.Client) (err error){
// TODO: make folder fedi-go
    instance, _ := json.Marshal(gClient)
    err = ioutil.WriteFile("/home/anjan/output.json", instance, 0644)

    return err
}
