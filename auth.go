package main
import (
       "fmt"
       "github.com/McKael/madon"
       "errors"
)

func getAuth(instance string) (url string, gClient *madon.Client){

    var scopes = []string{"read", "write", "follow"}
    gClient, err := madon.NewApp(APPNAME, APPWEBSITE, scopes, madon.NoRedirect, instance) 

    if err != nil {
	fmt.Println(err)
	return "", nil
    }

    url, err = gClient.LoginOAuth2("", scopes)

    if err != nil {
	fmt.Println(err)
	return "", nil
    } else {
	fmt.Println(url)
	return url, gClient
    }
}

// I removed gClient *madon.Client from output

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

func getAuthResume(client ClientStruct) (gClient *madon.Client){
	gClient, err := madon.RestoreApp(APPNAME, client.InstanceURL, client.ID, client.Secret, nil)
	if err != nil {
		fmt.Println(err)
	}

	return gClient
}
