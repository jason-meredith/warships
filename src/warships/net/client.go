package net

import (
	"fmt"
	"net/rpc"
	"strconv"
)


func JoinServer(username, address string) string {

	client, err := rpc.DialHTTP("tcp", address + ":" + strconv.Itoa(RPC_PORT))
	if err != nil {
		panic(err)
	}

	var playerId string

	err = client.Call("Server.JoinGame", &username, &playerId)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Joined Game with ID %v\n", playerId)

	return playerId

}