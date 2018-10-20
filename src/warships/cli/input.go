package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func clearScreen() {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}

type MenuOption func()
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
		results[option], _ = reader.ReadString('\n')
	}

	return results
}

func setupScreen() {
	clearScreen()
	fmt.Println(" ~ By Jason Meredith ~ \n")
	fmt.Println(logo)
}
