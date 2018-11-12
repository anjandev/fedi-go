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

func timelineGetter(gClient *madon.Client) (statuses []madon.Status){
    opt := timelineOpts

    var limOpts *madon.LimitParams
    limOpts = new(madon.LimitParams)
    limOpts.Limit = int(opt.limit)
    limOpts.MaxID = opt.maxID
    limOpts.SinceID = opt.sinceID

    statuses, err := gClient.GetTimelines("public", false, false, limOpts)
    if err != nil {
	fmt.Println(err)
    }
    return statuses
}
