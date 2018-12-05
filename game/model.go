package game

import (
	"encoding/json"
	"errors"
	"math/rand"
	"strconv"
	"PentagoServer/pentago"
	"fmt"
	"PentagoServer/db"
)

type Game struct {
	Id       int
	Mode string
	Turn pentago.Piece
	Players  [2]int
	Board pentago.Board
	LastTurn pentago.Move
}

// Move executes the given Move on a game state.
func (g *Game) Move(m pentago.Move) bool {
	if !g.Board.ApplyMove(m, g.Turn) {
		return false
	}
	g.LastTurn = m

	switch g.Turn {
	case pentago.White: g.Turn = pentago.Black
	case pentago.Black: g.Turn = pentago.White
	}
	return true
}


func GetGame(id int) (Game, error) {
	game := Game{}
	client := db.GetClientRedis()
	val, err := client.Get(strconv.Itoa(id)).Result()
	if err != nil {
		return game, err
	}

	errJson := json.Unmarshal([]byte(val), &game)
	if errJson != nil {
		return game, errJson
	}

	return game, nil
}

func GetColor(color int) {}

func SetGame(value Game) error {
	client := db.GetClientRedis()
	str, err := json.Marshal(value)
	if err != nil {
		return err
	}
	errSet := client.Set(strconv.Itoa(value.Id), str, 0).Err()
	if errSet != nil {
		return errSet
	}
	return nil
}

func CreateGame(p GameParams) (Game, error) {
	game := Game{
		Id:      rand.Int(),
		Mode:p.Mode,
		Players: [2]int{p.UserId, 0},
		Board: pentago.NewBoard(),
		LastTurn: pentago.Move{},
		Turn: pentago.White,
	}
	err := SetGame(game)
	if err != nil {
		return game, err
	}
	return game, nil
}

func JoinGame(gameId int, userId int) (Game, error) {
	game, err := GetGame(gameId)
	if err != nil {
		return game, err
	}
	game.Players[1] = userId
	game.Turn = pentago.White
	errSet := SetGame(game)
	if errSet != nil {
		return game, errSet
	}
	return game, nil
}

func Move(m MoveParams) (Game, error) {
	game, err := GetGame(m.Id)
	if err != nil {
		return game, err
	}

	if m.Position[0] < 0 || m.Position[0] > 5 || m.Position[1] < 0 || m.Position[1] > 5 {
		return game, errors.New("The marble is outside the board.")
	}

	userMove := pentago.NewMove(m.Position[0],m.Position[1],m.Quarter,m.Direction)

	if !userMove.IsValid(game.Board) {
		return game, errors.New("There is already a marble at this position.")
	}

	ok := game.Move(userMove)
	if !ok{
		return game, fmt.Errorf("Wrong move")
	}

	if game.Mode == "ai" {
		aiMove := game.Board.BestMove(game.Turn)
		ok := game.Move(aiMove)
		if !ok{
			return game, fmt.Errorf("Wrong move")
		}
	}

	errSet := SetGame(game)
	if errSet != nil {
		return game, errSet
	}

	return game, nil
}