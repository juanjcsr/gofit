package fitbit

import (
	"log"
	"net/http"

	"fmt"

	"github.com/gorilla/mux"
)

type Server struct {
	authedClient *http.Client
	router       *mux.Router
}

func NewFitbitServer(port int16) {
	// prefs, _ := createOrGetPreferences()
	address := fmt.Sprintf(":%d", port)
	router := mux.NewRouter()
	router.HandleFunc("/api/hola", getHolaFunc)
	log.Fatal(http.ListenAndServe(address, router))

}

func getHolaFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
