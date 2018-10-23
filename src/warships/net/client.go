package net

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

type CommandArgs struct {
	input 		string
	playerId	string
}


func ConnectToServer(address string) *rpc.Client {
	client, err := rpc.DialHTTP("tcp", address + ":" + strconv.Itoa(RPC_PORT))
	if err != nil {
		panic(err)
	}

	return client
}


// JoinServer takes a username and ip address and attempts to connect to running server
// create a new Player on that server and returns the ID of that new Player
func JoinServer(username, password, address string) (string, error) {

	client := ConnectToServer(address)

	login := LoginCredentials{
		Username:username,
		Password:password,
	}

	var details JoinDetails

	err := client.Call("Server.JoinGame", &login, &details)
	if err != nil {
		return "", err
	}


	fmt.Printf("Joined Game with ID %v\n", details.PlayerId)
	fmt.Printf("Assigned to team: %v\n", details.TeamName)

	return details.PlayerId, nil

}

func AcceptCommands(playerId string, address string) {

	// Get a connection to the server
	client := ConnectToServer(address)

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

		if _, ok := commands[fields[0]]; ok {
			var response string

			command := ClientCommand{PlayerId:playerId, Fields: fields }

			err := client.Call(commands[fields[0]], &command, &response)
			if err != nil {
				panic(err)
			}

			fmt.Println(response)
		} else {
			fmt.Printf("Command '%v' not found.\n", input)
		}

	}
}
