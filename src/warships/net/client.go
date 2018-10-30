package net

import (
	"bufio"
	"errors"
	"fmt"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

/*********************************************************
 *														 *
 *                   	  Warships						 *
 *					   Jason Meredith					 *
 *														 *
 *	DATE:		October 22, 2018						 *
 *	FILE: 		client.go								 *
 *	PURPOSE:	Handles a client connecting to a server  *
 *				and all outgoing commands				 *
 *				 										 *
 *														 *
 *********************************************************/

// ClientCommand wraps the PlayerID and command input into a single struct to send to server
type ClientCommand struct {
	PlayerId string
	Fields   []string
}

// CreateServerConnection takes a username, password and network address and attempts to connect
// to a Game server running at that location. If a user using that username has never connected
// to that server before a Player is created on the server with the given username and password.
//
// If a Player already exists on that server the password entered must be the password they entered
// when they first logged in or else they must re-login with a new username/password combo.
//
// Upon successful login a string (their UserID) and the RPC Client object are returned, or an error.
func CreateServerConnection(username, password, address string) (string, *rpc.Client, error) {

	// Create connection to server
	client, err := rpc.DialHTTP("tcp", address + ":" + strconv.Itoa(RPC_PORT))
	if err != nil {
		return "", nil, errors.New("unable to connect to that address")
	}

	login := LoginCredentials{
		Username:username,
		Password:password,
	}

	var details JoinDetails

	err = client.Call("Server.JoinGame", &login, &details)
	if err != nil {
		return "", nil, err
	}


	fmt.Printf("Joined Game with ID %v\n", details.PlayerId)
	fmt.Printf("Assigned to team: %v\n", details.TeamName)

	return details.PlayerId, client, nil

}

// AcceptCommands takes a the UserID and RPC Client object and prompts the user for
// keyboard input. If the first token of the input matches a key in hashmap of commands
// their UserID and full input are sent to the server wrapped in ClientCommand struct.
// Response string from server is printed to screen.
func AcceptCommands(playerId string, connection *rpc.Client) {


	reader := bufio.NewReader(os.Stdin)
	var input string

	// Map the client commands to remote function calls
	var commands map[string]string
	commands = make(map[string]string)
	commands["echo"] 		= "Server.EchoTest" // Test command
	commands["target"] 		= "Server.Target"	// Fire a shot at given location
	commands["sweep"]		= "Server.Sweep"	// Check a location for enemies without firing
	commands["map"]			= "Server.Map"		// Show team map
	commands["radar"]		= "Server.Radar"	// Show shots fired on enemy map
	commands["players"]		= "Server.Players"	// Show the player list

	for {

		if input == "quit" {
			os.Exit(0)
		}

		fmt.Printf("> ")

		// Get the input, divide into fields
		input, _ 	= reader.ReadString('\n')
		input 		= strings.TrimRight(input, "\n")
		fields 		:= strings.Fields(input)

		// If the first field in the input matches a key in the commands hashmap
		if _, ok := commands[fields[0]]; ok {
			var response string

			// Wrap command in ClientCommand struct
			command := ClientCommand{ PlayerId: playerId, Fields: fields }

			// Send ClientCommand to Server and print response
			err := connection.Call(commands[fields[0]], &command, &response)
			if err != nil {
				fmt.Println("Error: " + err.Error())
			}

			fmt.Println(response)
		} else {
			fmt.Printf("Command '%v' not found.\n", input)
		}

	}
}


