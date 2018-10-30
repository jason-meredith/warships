package game

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"math"
)


/*********************************************************
 *														 *
 *                   	  Warships						 *
 *					   Jason Meredith					 *
 *														 *
 *	DATE:		October 22, 2018						 *
 *	FILE: 		players.go								 *
 *	PURPOSE:	All the Player and Team management is	 *
 *				located here. Creating new Players, new	 *
 *				Teams, moving Players between Teams,	 *
 *				authenticating return users etc.		 *
 *														 *
 *				 										 *
 *														 *
 *********************************************************/


// Player represents a single human user connected to game
type Player struct {

	Username string
	Password string

	// Pointer to Team this player is on
	Team *Team

	Id 		string

}

// Team is a collection of Players working together on the same team
type Team struct {

	Name 		string

	Game *Game

	// Array of pointers to Player
	Players []*Player

	// Number of Players currently on Team
	NumPlayers int


	// A log of shots fired by this Team and shots upon this Team
	ShotsFired []Target
	ShotsUpon  []Target

	// This Team's Ships
	Ships []*Ship

}

// GetSmallestTeam when called on a Game returns the Team in the game with
// the least userse
func (game *Game) GetSmallestTeam() *Team {

	var smallestTeam *Team
	smallestAmount := math.MaxInt32

	for _, team := range game.Teams {
		if team.NumPlayers < smallestAmount {
			smallestAmount = team.NumPlayers
			smallestTeam = team
		}
	}

	return smallestTeam

}

// Adds a new Player to a Game, returns pointer to new Player
// This should be used to instantiate a new Player
func (game *Game) Join (username, password string) (*Player, bool, error) {

	// Generate ID based on Username
	id := RandomId(username, 32)

	// Check to see if that user is already registered in this Game
	// If not register the user
	player := game.GetPlayerById(id)
	if player == nil {

		team := game.GetSmallestTeam()

		// Create new player
		newPlayer := Player{username, password, team, id}

		// Add reference to player to Team.Players array
		team.Players = append(team.Players, &newPlayer)

		// Increment number of Players on Team
		team.NumPlayers += 1

		return &newPlayer, false, nil

	} else {
		// If Player already exists check the password
		if password == player.Password {
			return player, true, nil
		} else {
			return nil, true, errors.New("incorrect password")
		}
	}

}

func RandomId(input string, length int) string {
	// Generate random ID
	h := md5.New()
	seed := input
	io.WriteString(h, seed)
	id := fmt.Sprintf("%x", h.Sum(nil))
	return id[:length]
}

// NewTeam Instantiates a new Team
func (game *Game) NewTeam() *Team {



	team := Team{
	"",
	game,
	[]*Player{},
	0,
	[]Target{},
	[]Target{},
	[]*Ship{},
	}

	teamId := fmt.Sprintf("%p", &team)
	team.Name = fmt.Sprintf("Fleet-%v", RandomId(teamId, 5))

	game.Teams = append(game.Teams, &team)

	return &team
}

// GetPlaerById finds and return a Player using their Player ID
func (game *Game) GetPlayerById(id string) *Player {
	var result *Player = nil

	// Iterate through each Team
	for _, team := range game.Teams {
		// Iterate through each Player on each Team
		for _, player := range team.Players {
			if player.Id == id {
				result = player
			}
		}

	}

	return result
}

// SwitchTeam switches a Players Team
func SwitchTeam(player *Player, destTeam *Team) {

	// Get the Players original Team
	originalTeam := *player.Team

	// Find the index in the Team.Players array of the Player
	playerIndex := originalTeam.findPlayerIndex(player)

	// Remove that element from array
	originalTeam.Players = append(originalTeam.Players[:playerIndex],
		originalTeam.Players[playerIndex+1:]...)

	// Change the Player.Team value
	*player.Team = *destTeam

	// Increment destTeam player count
	destTeam.NumPlayers++

	// Add Player to destTeam.Players array
	destTeam.Players = append(destTeam.Players, player)

}

// findPlayerIndex finds the index of Player in a Teams Player array
func (team Team) findPlayerIndex(player *Player) int {
	playerList := team.Players

	for i, v := range playerList {
		if(v == player) {
			return i
		}
	}

	return -1
}