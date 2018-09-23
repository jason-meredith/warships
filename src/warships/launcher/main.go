package main

import (
	"fmt"
	"os"
	"warships/net"
)

func main() {

	fmt.Println("Warships Launcher 1.0.0\nBy Jason Meredith")


	if(os.Args[1] == "s") {
		net.StartServer()
	} else if(os.Args[1] == "c") {
		net.StartClient()
	}

}