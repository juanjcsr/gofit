package fitbit

import (
	"net/http"

	"fmt"

	"github.com/gorilla/mux"
)

type Subscriber struct {
	client           *http.Client
	port             int16
	router           *mux.Router
	verificationCode string
}

func newSubscriber(authedClient *http.Client, router *mux.Router, verificationCode string) (*Subscriber, error) {
	var subscriber = new(Subscriber)

	subscriber.client = authedClient
	subscriber.router = router
	subscriber.verificationCode = verificationCode

	router.HandleFunc("/api/fitbit_notification", subscriber.getFitBitNotification)
	return subscriber, nil

}

func (s *Subscriber) getFitBitNotification(w http.ResponseWriter, r *http.Request) {
	switch m := r.Method; m {
	case "GET":
		if v := r.URL.Query().Get("verify"); v != "" {
			fmt.Printf("\n SENT VERIFIC CODE: %s", v)
			fmt.Printf("\n VERIFICATION CODE: %s", s.verificationCode)
			if v == s.verificationCode {
				w.WriteHeader(http.StatusNoContent)
				fmt.Fprintf(w, "%s", "VERIFIED")
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "%s", "NOT VERIFIED")
			}

		}
	case "POST":

	}

}
