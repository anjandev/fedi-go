package main

import (
	"regexp"
       "github.com/McKael/madon"
)

func imgHTML(status madon.Status) (html string){
    if !status.Sensitive {
	html = "<!DOCTYPE html> <html> <body><center> <img src=\"" + status.MediaAttachments[0].URL + "\" style=\"max-height:300px; max-width:500px; object-fit: contain\"></center></body> </html> "

    }else if status.Sensitive {
    html = `<div id="spoiler" style="display:none">` + "<center><img src=\"" + status.MediaAttachments[0].URL + "\" style=\"max-height:300px; max-width:500px; object-fit: contain\"></center>" +
	    `</div>
	<a display="initial" id="button" title="Click to show/hide content" type="button" onclick="document.getElementById('spoiler').style.display='block';
	document.getElementById('button').style.display='none'">
        Click to view NSFW content
	</a>`
    }
    return html
}

func textHeight(contentInPost string) (height int) {
    // find number of lines to properly size box
    // I really dont like this because if someone posts a really long 1 liner, this algorithm
    // fails. Try something else. Check line size?
    paras:= regexp.MustCompile("</p>")
    matches := paras.FindAllStringIndex(contentInPost, -1)
    const HEIGHT_OF_LINE = 25
    return (len(matches)+2) * HEIGHT_OF_LINE
}
