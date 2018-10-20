package net

import (
	"fmt"
	"net/rpc"
	"strconv"
	"warships/game"
)

func JoinServer(username, address string) {

	client, err := rpc.DialHTTP("tcp", address + ":" + strconv.Itoa(RPC_PORT))
	if err != nil {
		panic(err)
	}

	var player game.Player

	err = client.Call("Server.JoinGame", &username, &player)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Joined Game: %#v\n", player )

}