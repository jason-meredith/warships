package net

import (
	"fmt"
	"net/rpc"
	"testing"
	"time"
)


func runServer() {

	path := "$GOPATH/src/warships/main/mainfff"
	serverComplete := make(chan Result)

	go RunServer(path,12, 12, 12, "pass", "adminpass", 32, serverComplete)


}

func endServer(connection *rpc.Client) {
	command := ClientCommand{ PlayerId: "0", Fields: []string{"shutdown", "adminpass"} }

	// Send ClientCommand to Server and print response
	connection.Call("Server.Shutdown", &command, nil)

}


func TestCreateServerConnection(t *testing.T) {

	t.Skip("runServer() doesn't work just yet in CI environment")

	runServer()
	fmt.Println("Server running, starting client")
	time.Sleep(1 * time.Second)
	userId, connection, err := CreateServerConnection("j", "j", "127.0.0.1")
	defer endServer(connection)

	if userId == "" {
		t.Error("Error creating connection to server, no userid")
	}

	if connection == nil {
		t.Error("Error creating connection to server, no RPC client")
	}

	if err != nil {
		t.Error("Error creating connection to server")
	}

}

func TestServer_Target(t *testing.T) {


	t.Skip("runServer() doesn't work just yet in CI environment")

	runServer()
	fmt.Println("Server running, starting client")
	time.Sleep(1 * time.Second)
	userId, connection, _ := CreateServerConnection("j", "j", "127.0.0.1")
	CreateServerConnection("k","k", "127.0.0.1")
	defer endServer(connection)

	var response string


	// Try targeting your their own team
	command := ClientCommand{ PlayerId: userId, Fields: []string{"target", "1", "A1"} }
	err := connection.Call("Server.Target", &command, &response)
	if err.Error() != "you cannot target your own team" {
		t.Error("Should have returned error when targeting own team")
	}

	// Try targeting non-existant team
	command = ClientCommand{ PlayerId: userId, Fields: []string{"target", "3", "A1"} }
	err = connection.Call("Server.Target", &command, &response)
	if err.Error() != "not a valid target number. Run 'teams' to see a list of teams and their team#" {
		t.Error("Should have returned error when targeting non-existing team number")
	}


	// Try targeting  non-existant team
	command = ClientCommand{ PlayerId: userId, Fields: []string{"target", "0", "A1"} }
	err = connection.Call("Server.Target", &command, &response)
	if err.Error() != "not a valid target number. Run 'teams' to see a list of teams and their team#" {
		t.Error("Should have returned error when targeting non-existing team number")
	}


	// Try targeting your their own team
	command = ClientCommand{ PlayerId: userId, Fields: []string{"target", "2" } }
	err = connection.Call("Server.Target", &command, &response)
	if err.Error() != "not enough arguments to perform target command: target <team#> <target_coordinate>" {
		t.Error("Should have returned error when not enough args in command")
	}


	// Try targeting your their own team



	command = ClientCommand{ PlayerId: userId, Fields: []string{"target", "2", "c4"} }
	err = connection.Call("Server.Target", &command, &response)
	if response != "Shot confirmed HIT!\n1 hit streak\n" {
		t.Error("Hit should have registered as a hit")
	}



}

