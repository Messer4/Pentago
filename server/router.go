package server

import (
	"PentagoServer/game"
	"PentagoServer/service"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func Start() {
	rand.Seed(time.Now().UTC().UnixNano())
	router := mux.NewRouter()

	router.HandleFunc("/registration", service.CreateUser).Methods("POST")
	router.HandleFunc("/login", service.Login).Methods("POST")
	router.HandleFunc("/results", service.Results).Methods("GET")


	router.HandleFunc("/game", game.GameRoute).Methods("POST")
	router.HandleFunc("/join", game.JoinRoute).Methods("POST")
	router.HandleFunc("/move", game.MoveHandler).Methods("POST")
	router.HandleFunc("/game/{id}", game.RefreshGameRoute).Methods("GET")
	router.HandleFunc("/last_turn/{id}", game.GetLastTurn).Methods("GET")

	fmt.Println("Listening on port : 6000")
	log.Fatal(http.ListenAndServe(":6000", router))
}
