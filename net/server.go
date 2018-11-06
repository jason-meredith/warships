package net

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"time"
	"warships/base26"
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
	fmt.Printf("\t-Listening on port %d\n", newGame.Port)
	fmt.Printf("\t-Max Players: %d\n", newGame.MaxPlayers)
	fmt.Printf("\t-Ship Limit: %d\n", newGame.ShipLimit)
	fmt.Printf("\t-Board Size: %d\n", newGame.BoardSize)

	// Create the Server object using the Game generated and passed to us by the CLI
	server := new(Server)
	server.game = newGame
	server.game.Teams = []*game.Team{}
	teamA := server.game.NewTeam()
	teamB := server.game.NewTeam()


	//TODO///////////////////  SHIP TEST DELEEEEETE


	teamA.NewShip(5, game.HORIZONTAL, game.Coordinate{2,2})

	teamA.NewShip(5, game.VERTICAL, game.Coordinate{2,4})

	teamB.NewShip(5, game.VERTICAL, game.Coordinate{2,2})
	teamB.NewShip(5, game.HORIZONTAL, game.Coordinate{4,4})


	//TODO//////////////////////////////////////////

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
		fmt.Printf("New player connected: %v\n", login.Username)
	} else {
		fmt.Printf("Existing player reconnected: %v\n", login.Username)
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

func (t *Server) Map(args ClientCommand, response *string) error {

	// Get the Team Map based on the Player who called the command
	player := t.game.GetPlayerById(args.PlayerId)
	teamMap := t.game.GetMap(player.Team)

	// Parse full command to determine section of map to render

	// Produce a string and put in response
	output := "  "

	// Top row
	for x := 0; x <= int(t.game.BoardSize); x++ {
		output += fmt.Sprintf("%-2v", base26.ConvertToBase26(x))
	}
	output += "\n"
	for y:= 0; y < int(t.game.BoardSize); y++ {
		output += fmt.Sprintf("%3v ", strconv.Itoa(y))
		for x:= 0; x < int(t.game.BoardSize); x++ {
			output += teamMap[x][y]
		}
		output += "\n"
	}

	*response = output

	timeStamp()
	fmt.Printf("Map Request\n")
	fmt.Printf("\t-Player: %v (%v)\n", player.Username, args.PlayerId)


	return nil
}

// Teams serves a list of all the Teams playing on this server, with a * in front
// of the calling Player's Team
func (t *Server) Teams(args ClientCommand, response *string) error {
	output := ""

	player := t.game.GetPlayerById(args.PlayerId)

	for id, team := range t.game.Teams {
		var strId string = ""

		if player.Team == team {
			strId = fmt.Sprintf("*%v", id + 1)
		} else {
			strId = fmt.Sprintf("%v", id + 1)
		}

		output += fmt.Sprintf("%3v:\t%v\n", strId, team.Name)
	}

	*response = output


	timeStamp()
	fmt.Printf("Team List Request\n")
	fmt.Printf("\t-Player: %v (%v)\n", player.Username, args.PlayerId)

	return nil
}

// Players serves a list of Players on a given team# (team# based on Teams command)
func (t *Server) Players(args ClientCommand, response *string) error {
	output := ""

	teamNum, err := strconv.Atoi(args.Fields[1])
	if err != nil {
		return errors.New("Team selection invalid")
	}

	if teamNum < 1 || teamNum > len(t.game.Teams) {
		return errors.New("Team selection out of range")

	}

	team := t.game.Teams[teamNum - 1]

	output += fmt.Sprintf("\n%v [ %v player(s) ]\n", team.Name, team.NumPlayers)
	output += fmt.Sprintf("%8v %-20v\n", "Points", "Username")

	for _, player := range team.Players {
		output += fmt.Sprintf("%8v %-20v\n", player.Points, player.Username )
	}

	*response = output

	return nil
}

// Target fires a shot
func (t *Server) Target(args ClientCommand, response *string) error {

	player := t.game.GetPlayerById(args.PlayerId)

	// command structure: 	target [team#] [Target{}]
	// 						target 2 G7

	if len(args.Fields) < 3 {
		return errors.New("not enough arguments to perform target command: target <team#> <target_coordinate>")
	}

	teamNum, err := strconv.Atoi(args.Fields[1])
	team := t.game.Teams[teamNum - 1]
	if team == player.Team {
		return errors.New("you cannot target you're own team")
	}

	// Parse into Target{} (split letters from numbers)
	target, err := game.StringToTarget(args.Fields[2])
	if err != nil {
		return err
	}

	output := ""

	shotResult := game.FireShot(player, team, target)

	if shotResult == game.HIT {
		output += "Shot confirmed HIT!\n"
		output += fmt.Sprint("%v hit streak\n", player.HitStreak)
	} else if shotResult == game.REPEAT_HIT {
		output += "Shot confirmed HIT but no further damage inflicted!\n"
	} else if shotResult == game.MISS {
		output += "Shot confirmed MISS!\n"
	} else if shotResult == game.SINK {
		output += "Shot confirmed HIT... enemy ship SUNK!\n"
	}

	*response = output

	return nil

}