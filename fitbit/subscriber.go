package fitbit

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Subscriber struct {
	client *http.Client
	port   int16
	router *mux.Router
}

func newSubscriber(authedClient *http.Client, port int16, router *mux.Router) (*Subscriber, error) {
	var subscriber = new(Subscriber)
	router.HandleFunc("/api/subscriber", sayHello)
	subscriber.client = authedClient
	subscriber.port = port
	subscriber.router = router
	return subscriber, nil

}

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HOLA"))
}
