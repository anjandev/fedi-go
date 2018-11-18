package main

import (
	"fmt"
       "github.com/McKael/madon"
)

var timelineOpts struct {
	local, onlyMedia bool
	limit, keep      uint
	sinceID, maxID   int64
}

func timelineGetter(gClient *madon.Client, maxId int64, minId int64, timeline string) (statuses []madon.Status){
    opt := timelineOpts

    var limOpts *madon.LimitParams
    limOpts = new(madon.LimitParams)
    limOpts.Limit = int(opt.limit)
    limOpts.MaxID = maxId
    limOpts.SinceID = minId


    statuses, err := gClient.GetTimelines(timeline, false, false, limOpts)
    if err != nil {
	fmt.Println(err)
    }
    return statuses
}

func timelineGetterStream(gClient *madon.Client) (chan madon.StreamEvent) {

    evChan := make(chan madon.StreamEvent, 10)
    stop := make(chan bool)
    done := make(chan bool)
    var err error
    streamName := "public" // can be user, local, public, or hashtag

    err = gClient.StreamListener(streamName, "", evChan, stop, done)


    if err != nil {
	fmt.Println("Stream Error")
	fmt.Println(err)
    }

    return evChan

}
