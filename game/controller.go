package game

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type GameParams struct {
	UserId int
	Mode   string
}

type JoinParams struct {
	UserId int
	Id     int
}

type MoveParams struct {
	UserId   int
	Id       int
	Position [2]int
	Quarter   int
	Direction int
}


func GameRoute(writer http.ResponseWriter, router *http.Request) {
	decoder := json.NewDecoder(router.Body)
	var params GameParams
	err := decoder.Decode(&params)
	if err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		writer.Write([]byte("422 - Cannot read the request's body."))
		return
	}

	game, err := CreateGame(params)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("500 - Cannot initiate the game."))
		return
	}

	payload, err := json.Marshal(game)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("500 - Something bad happened."))
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(payload)
}

func JoinRoute(writer http.ResponseWriter, router *http.Request) {
	decoder := json.NewDecoder(router.Body)
	var params JoinParams
	err := decoder.Decode(&params)
	if err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		writer.Write([]byte("422 - Cannot read the request's body."))
		return
	}

	game, err := JoinGame(params.Id, params.UserId)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("400 - Cannot join the game."))
		return
	}

	payload, err := json.Marshal(game)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("500 - Something bad happened."))
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(payload)
}

func MoveHandler(writer http.ResponseWriter, router *http.Request) {
	decoder := json.NewDecoder(router.Body)
	var params MoveParams

	err := decoder.Decode(&params)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusUnprocessableEntity)
		writer.Write([]byte("422 - Cannot read the request's body."))
		return
	}

	game, err := Move(params)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("400 - Cannot add the marble to the game."))
		return
	}

	payload, err := json.Marshal(game)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("500 - Something bad happened."))
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(payload)
}


func RefreshGameRoute(writer http.ResponseWriter, router *http.Request) {
	vars := mux.Vars(router)
	gameId, _ := strconv.Atoi(vars["id"])

	game, err := GetGame(gameId)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("400 - Cannot load the game."))
		return
	}

	payload, err := json.Marshal(game)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("500 - Something bad happened."))
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(payload)
}

func GetLastTurn(writer http.ResponseWriter, router *http.Request) {
	vars := mux.Vars(router)
	gameId, _ := strconv.Atoi(vars["id"])

	game, err := GetGame(gameId)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("400 - Cannot load the game."))
		return
	}

	payload, err := json.Marshal(game.LastTurn)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("500 - Something bad happened."))
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(payload)
}