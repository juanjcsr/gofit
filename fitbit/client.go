package fitbit

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/fitbit"
)

const (
	accessTokenString  = "access_token"
	refreshTokenString = "refresh_token"
	expiryTokenString  = "expiry_token_time"
	tokenTypeString    = "token_type"
)

// var fitbitConf *oauth2.Config
var oauthCode chan oauth2.Token
var chanError chan error

var oauthToken *oauth2.Token

var fitbitConf = &oauth2.Config{
	ClientID:     os.Getenv("FITBIT_CLIENT"),
	ClientSecret: os.Getenv("FITBIT_SECRET"),
	//RedirectURL:  "http://localhost:3000/oauth",
	RedirectURL: os.Getenv("FITBIT_REDIRECT_URL"),
	Endpoint:    fitbit.Endpoint,
	Scopes: []string{
		"activity",
		"heartrate",
		"location",
		"nutrition",
		"profile",
		"settings",
		"sleep",
		"social",
		"weight",
	},
}

type Client struct {
	Client     *http.Client
	User       *UserService
	Activities *ActivityService
	Heart      *HeartService
}

func NewFitbitClient() (*Client, error) {
	prefs := new(Preferences)
	if _, err := prefs.Open(); err != nil {
		fmt.Printf("error opening prefs %v", err)
	}
	defer prefs.Close()

	token, noToken := getFitBitKeys(prefs)
	if noToken != nil {
		oauthCode = make(chan oauth2.Token)
		chanError = make(chan error)

		code, _ := getServerlessAuthCode()
		t, _ := fitbitConf.Exchange(oauth2.NoContext, code)
		token = t

		// http.HandleFunc("/", handeAuth)
		// http.HandleFunc("/oauth", handleOauth)

		// go func() {
		// 	http.ListenAndServe(":3000", nil)
		// }()
		// select {
		// case code := <-oauthCode:
		// 	fmt.Printf("YA LLEGOO %v", code)
		prefs.Update(accessTokenString, token.AccessToken)
		prefs.Update(expiryTokenString, token.Expiry.String())
		prefs.Update(tokenTypeString, token.TokenType)
		prefs.Update(refreshTokenString, token.RefreshToken)
		// 	token = &code

		// case chanErr := <-chanError:
		// 	fmt.Printf("ERROR EN ALGUN LUGAR %s", chanErr)
		// 	return nil, fmt.Errorf("fitbit client: can't create token")
		// }
	}
	fmt.Printf("TENGO UN TOKEN %v", token.AccessToken)
	// tokenSource := fitbitConf.TokenSource(oauth2.NoContext, token)
	// transport := &oauth2.Transport{Source: ts}

	tokenSource := fitbitConf.TokenSource(oauth2.NoContext, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		log.Fatalln(err)
	}

	if newToken.AccessToken != token.AccessToken {
		prefs.Update(accessTokenString, newToken.AccessToken)
		prefs.Update(expiryTokenString, newToken.Expiry.String())
		prefs.Update(tokenTypeString, newToken.TokenType)
		prefs.Update(refreshTokenString, newToken.RefreshToken)
	}

	// client := fitbitConf.Client(oauth2.NoContext, token)
	client := oauth2.NewClient(oauth2.NoContext, tokenSource)
	fClient := &Client{
		Client:     client,
		User:       newUserService(client),
		Activities: newActivityService(client),
		Heart:      newHeartService(client),
	}
	// fmt.Printf("TOKEN: %v\n", token.AccessToken)
	return fClient, nil
}

func createOrGetPreferences() (*Preferences, error) {
	prefs := new(Preferences)
	if _, err := prefs.Open(); err != nil {
		return nil, fmt.Errorf("error opening prefs %v", err)
	}
	return prefs, nil
}

func handleOauth(w http.ResponseWriter, r *http.Request) {
	// TODO: check for redirect state first
	code := r.FormValue("code")
	token, err := fitbitConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		// fmt.Printf("Failed code exchanging with '%s'\n", err)
		chanError <- fmt.Errorf("failed code exchanging with '%s'", err)
	}
	fmt.Printf("REFRESH TIME: %v\n", token.Expiry)

	oauthCode <- *token
}

func handeAuth(w http.ResponseWriter, r *http.Request) {
	url := fitbitConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func getServerlessAuthCode() (string, error) {
	the_url := fitbitConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Go to: %s\n", the_url)
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		return "", fmt.Errorf("could not read from cli: %v", err)
	}
	// fmt.Println(response)
	return response, nil
}

func getFitBitKeys(prefs *Preferences) (*oauth2.Token, error) {
	token := prefs.Read(accessTokenString)
	tokenType := prefs.Read(tokenTypeString)
	refreshToken := prefs.Read(refreshTokenString)
	expiry := prefs.Read(expiryTokenString)

	if token == "" || tokenType == "" || refreshToken == "" || expiry == "" {
		return nil, fmt.Errorf("fitbit keys: no keys")
	}

	tk := &oauth2.Token{}
	tk.AccessToken = token
	tk.TokenType = tokenType
	tk.RefreshToken = refreshToken
	// expT.Add(time.Duration(expTime) * time.Second)
	tk.Expiry, _ = time.Parse("2006-01-02 15:04:05 -0700 MST", expiry)
	// fmt.Printf("REFRESH EXPIRY ERR: %v\n", err)
	// fmt.Printf("REFRESH EXPIRY DB RAW: %v\n", expiry)
	// fmt.Printf("REFRESH EXPRIY DB XXX: %v\n", tk.Expiry)
	return tk, nil
}
