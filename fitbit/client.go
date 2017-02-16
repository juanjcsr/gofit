package fitbit

import (
	"fmt"
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
	RedirectURL:  "http://localhost:3000/oauth",
	Endpoint:     fitbit.Endpoint,
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

type FitBitClient struct {
	Client *http.Client
}

func NewFitbitClient() (*FitBitClient, error) {
	prefs := new(Preferences)
	if _, err := prefs.Open(); err != nil {
		fmt.Printf("error opening prefs %v", err)
	}
	defer prefs.Close()

	token, noToken := getFitBitKeys(prefs)
	if noToken != nil {
		oauthCode = make(chan oauth2.Token)
		chanError = make(chan error)

		http.HandleFunc("/", handeAuth)
		http.HandleFunc("/oauth", handleOauth)

		go func() {
			http.ListenAndServe(":3000", nil)
		}()
		select {
		case code := <-oauthCode:
			fmt.Printf("YA LLEGOO %v", code)
			prefs.Update(accessTokenString, code.AccessToken)
			prefs.Update(expiryTokenString, code.Expiry.String())
			prefs.Update(tokenTypeString, code.TokenType)
			prefs.Update(refreshTokenString, code.RefreshToken)
			token = &code

		case chanErr := <-chanError:
			fmt.Printf("ERROR EN ALGUN LUGAR %s", chanErr)
			return nil, fmt.Errorf("fitbit client: can't create token")
		}
	}
	// fmt.Printf("TENGO UN TOKEN %v", token)
	// tokenSource := fitbitConf.TokenSource(oauth2.NoContext, token)
	// transport := &oauth2.Transport{Source: ts}

	client := fitbitConf.Client(oauth2.NoContext, token)
	fClient := &FitBitClient{Client: client}
	return fClient, nil
}

func handleOauth(w http.ResponseWriter, r *http.Request) {
	// TODO: check for redirect state first
	code := r.FormValue("code")
	token, err := fitbitConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		// fmt.Printf("Failed code exchanging with '%s'\n", err)
		chanError <- fmt.Errorf("failed code exchanging with '%s'", err)
	}

	oauthCode <- *token
}

func handeAuth(w http.ResponseWriter, r *http.Request) {
	url := fitbitConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
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
	tk.Expiry, _ = time.Parse("2006-01-02 15:04:05.999999999 -0600 CST", expiry)

	return tk, nil
}
