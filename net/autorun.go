package net

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
)
/*********************************************************
 *														 *
 *                   	  Warships						 *
 *					   Jason Meredith					 *
 *														 *
 *	DATE:		December 2, 2018						 *
 *	FILE: 		game.go								 	 *
 *	PURPOSE:	Runs an instance of the Warships server  *
 *				automatically for testing purposes		 *
 *														 *
 *				 										 *
 *														 *
 *********************************************************/


type Result struct{
	output string
	err error
}


func RunServer(path string, maxPlayers, boardSize, shipLimit int, pass, adminPass string, deployPoints int, output chan Result) {

	fmt.Println("Running Server...")


	cmd := exec.Command(path)



	var in bytes.Buffer
	// Start new game
	in.Write([]byte("1\n"))
	in.Write([]byte(strconv.Itoa(maxPlayers)))
	in.Write([]byte(strconv.Itoa(boardSize)))
	in.Write([]byte(strconv.Itoa(shipLimit)))
	in.Write([]byte(pass))
	in.Write([]byte(adminPass))
	in.Write([]byte(strconv.Itoa(deployPoints)))

	cmd.Stdin = &in

	var out bytes.Buffer
	cmd.Stdout = &out

	var errbuf bytes.Buffer
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil {
		//log.Fatal(err)
		println(errbuf.String())
		output <- Result{out.String() + errbuf.String(), err}
	}

	output <- Result{out.String(), nil }


}
