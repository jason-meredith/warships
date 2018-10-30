package net

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"time"
	"warships/game"
)


/*********************************************************
 *														 *
 *                   	  Warships						 *
 *					   Jason Meredith					 *
 *														 *
 *	DATE:		October 22, 2018						 *
 *	FILE: 		server.go								 *
 *	PURPOSE:	Contains the Server struct, a wrapper	 *
 *				struct containing the game that has an	 *
 *				RPC listener that listens for incoming	 *
 *				Client commands. Any function attached	 *
 *				to the Server struct is contained here	 *
 *				and defines a command that can be called *
 *				by a Client.							 *
 *														 *
 *				The Game object has a lot of functions	 *
 *				attached to it for internal and 		 *
 *				administrative use. Not every function on*
 *				the Game needs to be attached to the RCP *
 *				Listener. The Server will listen for	 *
 *				commands and delegate actions to the game*
 *				where appropriate.						 *
 *				 										 *
 *														 *
 *********************************************************/

// Server is a struct exposed to clients and acts a buffer between clients
// and the actual Game
type Server struct {
	//Test 	int
	game	*game.Game
}

// LoginCredentials is the Username/Password combo passed by the client when
// attempting to log in
type LoginCredentials struct {
	Username	string
	Password	string
}

// JoinDetails is information sent back to the Client after a successful login
// telling the Client program their PlayerID and the team they've been assigned
type JoinDetails struct {
	PlayerId 	string
	TeamName 	string
}

// RPC_PORT is the TCP port that the server listens to
const RPC_PORT = 51832

// StartGameServer creates the Server using a new Game, sets up the RPC Listener
// and handles all incoming Client requests.
func StartGameServer(newGame *game.Game) {


	timeStamp()
	fmt.Println("Starting Server")

	// Create the Server object using the Game generated and passed to us by the CLI
	server := new(Server)
	server.game = newGame
	server.game.Teams = []*game.Team{}
	server.game.NewTeam()
	server.game.NewTeam()

	// Register the server for Remote Procedure Calls
	rpc.Register(server)
	rpc.HandleHTTP()

	// Listen on the RPC port for incoming commands
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(RPC_PORT))
	if err != nil {
		fmt.Println("Error encountered when starting server... is port " + strconv.Itoa(RPC_PORT) + " open?")
		fmt.Println("Program will now exit. Try changing the Server Listen Port or freeing port " + strconv.Itoa(RPC_PORT))
		os.Exit(1)
	}

	go http.Serve(listener, nil)

	// Loop for as long as Game is 'live'
	for server.game.Live {

	}

}

// timeStamp is the header for all server log entries
func timeStamp() {
	fmt.Printf("\n --- %v ---\n", time.Now())
}

// PrintTeamCounts prints a list of all the Teams and the number of users on each team
func PrintTeamCounts(game *game.Game) {
	for _, team := range game.Teams {
		fmt.Printf("\t%v: %v\n", team.Name, team.NumPlayers)
	}
}

// JoinGame joins a Player to the running Server using LoginCredentials.
func (t *Server) JoinGame(login LoginCredentials, info *JoinDetails) error {

	player, existing, err := t.game.Join(login.Username, login.Password)
	if err != nil {
		return err
	}

	// Details to send back to Client
	*info = JoinDetails{
		player.Id,
		player.Team.Name,
	}

	// Print details about this incoming command to the log
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

// EchoTest is used to confirm we are connected and the Client can send commands,
// the Server can receive them, and the Server can send a response that the Client
// can receive
func (t *Server) EchoTest(args ClientCommand, response *string) error {

	*response = fmt.Sprintf("Echo command successful\n%#v\n", args)

	player := t.game.GetPlayerById(args.PlayerId)

	timeStamp()
	fmt.Printf("Echo command received\n")
	fmt.Printf("\t-Player: %v (%v)\n", player.Username, args.PlayerId)
	fmt.Printf("\t-Fields: %v\n", args.Fields)

	return nil

}