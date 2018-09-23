package net

import (
	"fmt"
	"os"
)

type Command struct {
	//player Player
	CmdString string
}


func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
