package game

import (
	"crypto/md5"
	"fmt"
	"io"
	"math"
	"time"
)

type Player struct {

	Username string
	Score    int

	// Pointer to Team this player is on
	Team *Team

	Id 		string

}

type Team struct {

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

// Adds a new Player to a Team, returns pointer to new Player
// This should be used to instantiate a new Player
func (team *Team) NewPlayer (username string) (*Player, error) {

	// Generate New Player ID
	h := md5.New()
	seed := username + time.Now().String()
	io.WriteString(h, seed)
	id := fmt.Sprintf("%x", h.Sum(nil))

	// Create new player
	newPlayer := Player{username, 0, team, id}

	// Add reference to player to Team.Players array
	team.Players = append(team.Players, &newPlayer)

	// Increment number of Players on Team
	team.NumPlayers += 1


	return &newPlayer, nil

}

// Instantiates a new Team
func (game *Game) NewTeam() *Team {
	team := Team{ game,
	make([]*Player, 20),
	0,
	[]Target{},
	[]Target{},
	[]*Ship{},
	}

	game.Teams = append(game.Teams, &team)

	return &team
}

// Find and return a Player using their Player ID
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

// Switches a Players Team
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

// Find the index of Player in a Teams Player array
func (team Team) findPlayerIndex(player *Player) int {
	playerList := team.Players

	for i, v := range playerList {
		if(v == player) {
			return i
		}
	}

	return -1
}
/*

func main() {

	Team := NewTeam()
	otherTeam := NewTeam()

	fmt.Printf("Original\nTeam 1: %X\nTeam 2: %X\n\n", Team, otherTeam)

	player, _:= (&Team).NewPlayer("Jason")


	fmt.Printf("After add player\nTeam 1: %X\nTeam 2: %X\n\n", Team, otherTeam)

	SwitchTeam(player, &otherTeam)


	fmt.Printf("After switch player\nTeam 1: %X\nTeam 2: %X\n\n", Team, otherTeam)

	//fmt.Println(Team)
	//fmt.Println(player)
}*/