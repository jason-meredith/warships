package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

/*********************************************************
 *														 *
 *                   	  Warships						 *
 *					   Jason Meredith					 *
 *														 *
 *	DATE:		October 22, 2018						 *
 *	FILE: 		input.go								 *
 *	PURPOSE:	Provides helper functions for the menu	 *
 *				such as clearing the screen, displaying  *
 *				multi-option menus and parsing the 		 *
 *				response, as well as collecting info	 *
 *				in a form style input with multiple		 *
 *				fields									 *
 *				 										 *
 *														 *
 *														 *
 *********************************************************/

// clearScreen clears the screen. Currently only works with linux (linux -> clear, windows -> cls) but
// I intend to implement that
func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// MenuOption is a function representing the action of an inputMenu menu option
type MenuOption func()

// inputMenu displays a selection of options with corresponding numbers, collects user input, and
// runs the function associated with the selected Menu option
func inputMenu(prompt string, options ...MenuOption) {
	fmt.Println(prompt)

	reader := bufio.NewReader(os.Stdin)

	var optionSelection MenuOption

	for valid := false; !valid; {

		fmt.Print("> ")
		userInput, err := reader.ReadString('\n')

		option, err := strconv.Atoi(strings.TrimRight(userInput, "\n"))

		if err != nil || option > len(options) || option < 1 {
			fmt.Println("Invalid input, please enter number between 1 and ", len(options))
			continue
		}

		valid = true

		optionSelection = options[option - 1]

	}



	clearScreen()
	optionSelection()
}

// inputOptions collects user text input using an array of string prompts and stores
// responses in key-value map
func inputOptions(prompt string, options ...string) map[string]string {

	reader := bufio.NewReader(os.Stdin)

	var results map[string]string

	results = make(map[string]string)

	fmt.Println(prompt)
	for i := 0; i < len(prompt); i++ {
		fmt.Print("-")
	}
	fmt.Println("\n")

	for _, option := range options {
		fmt.Printf("%-20s\t", option)
		rawInput, _ := reader.ReadString('\n')
		results[option] = strings.TrimRight(rawInput, "\n")
	}

	return results
}

// setupScreen sets up the screen showing credits and the logo
func setupScreen() {
	clearScreen()
	fmt.Println(logo)
	fmt.Println("\t\t ~ By Jason Meredith ~ \n")
}
