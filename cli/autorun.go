package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

type Result struct{
	output string
	err error
}

func maifn() {
	serverComplete := make(chan Result)


	var cancel context.CancelFunc

	go RunServer(12, 12, 12, "pass", "admin", serverComplete, &cancel)

	result := <- serverComplete

	if result.err != nil {
		fmt.Println("Error")
		print(result.err)
	}

	//println(result.output)

}


func RunServer(maxPlayers, boardSize, shipLimit int, pass, adminPass string, output chan Result, cancelFunc *context.CancelFunc) {
	fmt.Println("Running Server...")
	path := "/home/jason/go/src/warships/cli/cli"

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Minute)
	*cancelFunc = cancel
	defer cancel()

	cmd := exec.CommandContext(ctx, path)



	var in bytes.Buffer
	// Start new game
	in.Write([]byte("1\n"))
	in.Write([]byte(strconv.Itoa(maxPlayers)))
	in.Write([]byte(strconv.Itoa(boardSize)))
	in.Write([]byte(strconv.Itoa(shipLimit)))
	in.Write([]byte(pass))
	in.Write([]byte(adminPass))

	cmd.Stdin = &in

	var out bytes.Buffer
	cmd.Stdout = &out

	var errbuf bytes.Buffer
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil {
		//log.Fatal(err)
		output <- Result{out.String() + errbuf.String(), err}
	}

	output <- Result{out.String(), nil }


}
