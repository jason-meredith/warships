package net

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"time"
	"warships/game"
)

// Server is a struct exposed to clients and acts a buffer between clients
// and the actual Game
type Server struct {
	//Test 	int
	game	*game.Game
}

type LoginCredentials struct {
	Username	string
	Password	string
}

type JoinDetails struct {
	PlayerId 	string
	TeamName 	string
}

type ClientCommand struct {
	PlayerId	string
	Fields 		[]string
}

const CONNECT_PORT = 51831
const RPC_PORT = 51832

// StartGameServer handles all incoming client requests
func StartGameServer(newGame *game.Game) {



	fmt.Println("Starting Server")

	//server := Server{ newGame }

	server := new(Server)
	server.game = newGame
	server.game.Teams = []*game.Team{}
	server.game.NewTeam()
	server.game.NewTeam()
	//server.game.Teams = append(server.game.Teams, new(game.Team))
	//server.game.Teams = append(server.game.Teams, new(game.Team))
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

// timeStamp is the header for all server log entries
func timeStamp() {
	fmt.Printf("\n --- %v ---\n", time.Now())
}

func PrintTeamCounts(game *game.Game) {
	for _, team := range game.Teams {
		fmt.Printf("\t%v: %v\n", team.Name, team.NumPlayers)
	}
}

// JoinGame joins a Player to the running Server
func (t *Server) JoinGame(login LoginCredentials, info *JoinDetails) error {

	player, existing, err := t.game.Join(login.Username, login.Password)
	if err != nil {
		return err
	}


	*info = JoinDetails{
		player.Id,
		player.Team.Name,
	}

	timeStamp()
	if !existing {
		fmt.Printf("New player joined: %v\n", login.Username)
	} else {
		fmt.Printf("Existing player rejoined: %v\n", login.Username)
	}
	fmt.Printf("\t-Player ID: %v\n", info.PlayerId)
	fmt.Printf("\t-Assigned to team %v (%p)\n", info.TeamName, player.Team)
	fmt.Printf("\n\t[Teams]\n")
	PrintTeamCounts(t.game)

	return err
}

func (t *Server) EchoTest(args ClientCommand, response *string) error {

	*response = fmt.Sprintf("Echo command successful\n%#v\n", args)

	player := t.game.GetPlayerById(args.PlayerId)

	timeStamp()
	fmt.Printf("Echo command received\n")
	fmt.Printf("\t-Player: %v (%v)\n", player.Username, args.PlayerId)
	fmt.Printf("\t-Fields: %v\n", args.Fields)

	return nil

}