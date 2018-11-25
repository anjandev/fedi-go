package main
import (
       "fmt"
       "github.com/McKael/madon"
       "errors"
)

func getClient(instance string) (gClient *madon.Client){

    var scopes = []string{"read", "write", "follow"}
    gClient, err := madon.NewApp(APPNAME, APPWEBSITE, scopes, madon.NoRedirect, instance) 

    // TODO: add error checking for bad instance
    if err != nil {
	panic(err)
    }
    return gClient
}

func getAuthOAuth(instance string, gClient *madon.Client) (url string){
    var scopes = []string{"read", "write", "follow"}
    url, err := gClient.LoginOAuth2("", scopes)

    if err != nil {
	fmt.Println(err)
	return ""
    } else {
	fmt.Println(url)
	return url
    }
}

func getAuthBasic(login string, password string, gClient *madon.Client) (authed bool){
    var scopes = []string{"read", "write", "follow"}
    err := gClient.LoginBasic(login, password, scopes)
    if err != nil {
	fmt.Println(err)
	return false
    }
    return true
}


func oAuth2ExchangeCode(tokenCode string, gClient *madon.Client) (err error) {
    var scopes = []string{"read", "write", "follow"}
	// (gClient != nil thanks to PreRun)

	if tokenCode == "" {
		return errors.New("no tokenCode entered")
	}

	// The code has been set; proceed with token exchange
	consentURL , err := gClient.LoginOAuth2(tokenCode, scopes)
	fmt.Println(consentURL)

	if err != nil {
		return err
	}

	if gClient.UserToken != nil {
		fmt.Println("Login successful.\n")
	}
	return nil
}

