package main

import (
	"log"
	"net/http"

	"fmt"

	"os"

	"github.com/gorilla/mux"
	fitbit "github.com/juanjcsr/gofit/fitbit"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
	"golang.org/x/oauth2"
	ofb "golang.org/x/oauth2/fitbit"
)

type Server struct {
	authedClient *http.Client
	router       *mux.Router
}

var fitbitServerConf = &oauth2.Config{
	ClientID:     os.Getenv("FITBIT_CLIENT"),
	ClientSecret: os.Getenv("FITBIT_SECRET"),
	RedirectURL:  os.Getenv("FITBIT_REDIRECT_URL"),
	Endpoint:     ofb.Endpoint,
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
	client       *fitbit.Client
	notification *fitbit.Subscriber
}

func NewFitbitServer(port int16) {
	// prefs, err := createOrGetPreferences()
	// if err != nil {
	// 	log.Printf("server: could not open prefs file %v", prefs)
	// }
	// defer prefs.Close()

	// fitClient, err := NewFitbitClient()
	// if err != nil {
	// 	log.Fatalf("could not create fitbit client: %v", err)
	// }

	var fs = new(FitServer)
	// fs.client = fitClient

	address := fmt.Sprintf(":%d", port)
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// redirectURL := getLoginURL()
	// verificationCode := os.Getenv("FITBIT_NOTIFY_VERIFICATION")

	// suscriber, err := newSubscriber(fitClient.Client, router, verificationCode)
	// if err != nil {
	// 	log.Fatalf("could not create subscriber endpoint: %v", err)
	// }
	// fs.notification = suscriber

	// router.HandleFunc("/api/activities/summary", fs.getHolaFunc)
	// router.HandleFunc("/api/activities/log", fs.getActivitiesLog)
	// uORM := &UserORM{db: db}

	router.Mount("/api/activities", fs.activitiesRouter())
	router.Mount("/api/users", userRouter())
	log.Fatal(http.ListenAndServe(address, router))
}

func (fs *FitServer) activitiesRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(WithUser)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("activities: index"))
	})
	r.Get("/summary", fs.getHolaFunc)
	r.Get("/log", fs.getActivitiesLog)
	return r
}

func WithUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("HOLA :D")
		next.ServeHTTP(w, r)
	})
}

func handleAuth(w http.ResponseWriter, r *http.Request) {

}

func getLoginURL() string {
	url := fitbitServerConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return url
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

type ErrorResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrorResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request",
		ErrorText:      err.Error(),
	}
}
