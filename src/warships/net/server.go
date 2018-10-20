package net

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"warships/game"
)

// Server is a struct exposed to clients and acts a buffer between clients
// and the actual Game
type Server struct {
	//Test 	int
	game	*game.Game
}

const RPC_PORT = 51832

// StartGameServer handles all incoming client requests
func StartGameServer(newGame *game.Game) {

	fmt.Println("Starting Server")

	//server := Server{ newGame }

	server := new(Server)
	server.game = newGame
	server.game.Teams = []*game.Team{}
	server.game.Teams = append(server.game.Teams, new(game.Team))
	server.game.Teams = append(server.game.Teams, new(game.Team))
	fmt.Printf("Added %v new team(s)\n", len(server.game.Teams))
	// Register the server for Remote Procedure Calls
	rpc.Register(server)
	rpc.HandleHTTP()

	// Listen on the RPC port for incoming commands
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(RPC_PORT))
	if err != nil {
		panic(err)
	}

	go http.Serve(listener, nil)

	for server.game.Live {

	}

}

func (t *Server) DoubleNum(num int, result *int) error {
	*result = num *2
	return nil
}

// JoinGame joins a Player to the running Server
func (t *Server) JoinGame(username string, player *game.Player) error {

	player, err := t.game.GetSmallestTeam().NewPlayer(username)

	fmt.Printf("New player joined: %v [%#v]\n", username, player)
	fmt.Printf("\tPlayer ID: %v\n", &player)
	fmt.Printf("\tAssigned to team %#v\n", player.Team)

	return err
}