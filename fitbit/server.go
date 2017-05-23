package fitbit

import (
	"log"
	"net/http"

	"fmt"

	"os"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/fitbit"
)

type Server struct {
	authedClient *http.Client
	router       *mux.Router
}

var fitbitServerConf = &oauth2.Config{
	ClientID:     os.Getenv("FITBIT_CLIENT"),
	ClientSecret: os.Getenv("FITBIT_SECRET"),
	RedirectURL:  os.Getenv("FITBIT_REDIRECT_URL"),
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

type FitServer struct {
	client *Client
}

func NewFitbitServer(port int16) {
	prefs, err := createOrGetPreferences()
	if err != nil {
		log.Printf("server: could not open prefs file %v", prefs)
	}
	defer prefs.Close()

	fitClient, err := NewFitbitClient()
	if err != nil {
		log.Fatalf("could not create fitbit client: %v", err)
	}

	var fs = new(FitServer)
	fs.client = fitClient

	address := fmt.Sprintf(":%d", port)
	router := mux.NewRouter()
	router.HandleFunc("/api/hola", fs.getHolaFunc)
	router.HandleFunc("/api/activities", fs.getActivitiesLog)
	log.Fatal(http.ListenAndServe(address, router))
}

func handleAuth(w http.ResponseWriter, r *http.Request) {

}

func (fs *FitServer) getActivitiesLog(w http.ResponseWriter, r *http.Request) {
	resp, err := fs.client.Activities.GetActivityLog("-", "2017-05-23", "", 2)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
	}
	tcx, _ := fs.client.Activities.GetTCXData(resp.Activities[0])
	fmt.Fprintf(w, "%s", tcx)
}

func (fs *FitServer) getHolaFunc(w http.ResponseWriter, r *http.Request) {
	resp, err := fs.client.Activities.GetActivitySummary("", "2017-05-21")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, "%s", resp.ToJson())
}
