package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"warships/game"
	"warships/net"
)

var logo = `
██╗    ██╗ █████╗ ██████╗ ███████╗██╗  ██╗██╗██████╗ ███████╗
██║    ██║██╔══██╗██╔══██╗██╔════╝██║  ██║██║██╔══██╗██╔════╝
██║ █╗ ██║███████║██████╔╝███████╗███████║██║██████╔╝███████╗
██║███╗██║██╔══██║██╔══██╗╚════██║██╔══██║██║██╔═══╝ ╚════██║
╚███╔███╔╝██║  ██║██║  ██║███████║██║  ██║██║██║     ███████║
 ╚══╝╚══╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═╝╚═╝     ╚══════╝
`

/*********************************************************
 *														 *
 *                   	  Warships						 *
 *					   Jason Meredith					 *
 *														 *
 *	DATE:		October 22, 2018						 *
 *	FILE: 		launcher.go								 *
 *	PURPOSE:	Entry point to the main menu			 *
 *				 										 *
 *				The main function loads you into a 		 *
 *				menu giving you the option of Start 	 *
 *				Game, Join Game or Quit.				 *
 *														 *
 *				Selecting Start Game will launch an		 *
 *				options menu allowing you to set up		 *
 *				preferences for the new Game Server.	 *
 *														 *
 *				Join Game will prompt you to enter		 *
 *				server IP Address, Username and Password *
 *														 *
 *				Quit will exit the program				 *
 *														 *
 *														 *
 *********************************************************/

// main Launches the game into the main menu, showing logo, credits and initial options
func main() {

	clearScreen()
	fmt.Println(logo)
	fmt.Println("\t\t ~ By Jason Meredith ~ \n")
	time.Sleep(700 * time.Millisecond)



	inputMenu("Select an option\n1. Start Game\n2. Join Game\n3. Quit",
		func() { startServer() },
		func() { joinGame() },
		func() { os.Exit(0)})

}

// startServer shows the menu screen for starting a new server
func startServer() {

	const PASSWRD		= "Password"
	const ADMIN_PASSWRD	= "Admin Password"
	const MAX_PLAYERS	= "Max Players"
	const SHIP_LIMIT	= "Ship Limit"
	const BOARD_SIZE	= "Board Size"


	setupScreen()

	options := inputOptions(
		"Start New Game",
					MAX_PLAYERS,
					BOARD_SIZE,
					SHIP_LIMIT,
					PASSWRD,
					ADMIN_PASSWRD,
		)

	maxPlayers, err := strconv.Atoi(options[MAX_PLAYERS])
	shipLimit, err := strconv.Atoi(options[SHIP_LIMIT])
	boardSize, err := strconv.Atoi(options[BOARD_SIZE])

	if err != nil {
		// TODO: Handle this error pls
	}

	newGame := game.Game{}
	newGame.Live =          true
	newGame.Port =          net.RPC_PORT
	newGame.Password =      options[PASSWRD]
	newGame.StartTime =     time.Now()
	newGame.AdminPassword = options[ADMIN_PASSWRD]
	newGame.MaxPlayers =    uint8(maxPlayers)
	newGame.ShipLimit =     uint8(shipLimit)
	newGame.BoardSize =     uint8(boardSize)
	newGame.Teams =         []*game.Team{}



	clearScreen()
	net.StartGameServer( &newGame )

	// Proceed to game
}

// joinGame shows the menu screen for joining a running game server
func joinGame() {

	const SERV_ADDR = "Server Address"
	const PASSWRD	= "Password"
	const USERNAME	= "Username"

	setupScreen()


	success := false

	for !success {
		options := inputOptions(
			"Joining Game",
			SERV_ADDR,
			PASSWRD,
			USERNAME,
		)

		playerId, connection, err := net.CreateServerConnection(
			strings.TrimRight(options[USERNAME], "\n"),
			strings.TrimRight(options[PASSWRD], "\n"),
			options[SERV_ADDR])
		if err == nil {
			success = true
			net.AcceptCommands(playerId, connection)
		} else {
			fmt.Println("Error: " + err.Error())
			time.Sleep(700 * time.Millisecond)
			setupScreen()
		}
	}


}