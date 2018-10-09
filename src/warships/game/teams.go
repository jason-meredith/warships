package game



type Player struct {

	username 	string
	score 		int

	// Pointer to team this player is on
	team 		*Team


}

type Team struct {

	game 		*Game

	// Array of pointers to Player
	players 	[]*Player

	// Number of players currently on team
	numPlayers 	int


	// A log of shots fired by this team and shots upon this team
	shotsFired		[]Target
	shotsUpon		[]Target

	// This team's ships
	ships			[]*Ship

}

// Adds a new Player to a Team, returns pointer to new Player
// This should be used to instantiate a new Player
func (team *Team) NewPlayer (username string) (*Player, error) {

	// Create new player
	newPlayer := Player{username, 0, team}

	// Add reference to player to Team.players array
	team.players = append(team.players, &newPlayer)

	// Increment number of Players on Team
	team.numPlayers += 1

	return &newPlayer, nil

}

// Instantiates a new Team
func (game *Game) NewTeam() Team {
	return Team{ game,make([]*Player, 20), 0, []Target{}, []Target{}, []*Ship{} }
}

// Switches a Players team
func SwitchTeam(player *Player, destTeam *Team) {

	// Get the Players original Team
	originalTeam := *player.team

	// Find the index in the Team.players array of the Player
	playerIndex := originalTeam.findPlayerIndex(player)

	// Remove that element from array
	originalTeam.players = append(originalTeam.players[:playerIndex],
		originalTeam.players[playerIndex+1:]...)

	// Change the Player.team value
	*player.team = *destTeam

	// Increment destTeam player count
	destTeam.numPlayers++

	// Add Player to destTeam.players array
	destTeam.players = append(destTeam.players, player)

}

// Find the index of Player in a Teams Player array
func (team Team) findPlayerIndex(player *Player) int {
	playerList := team.players

	for i, v := range playerList {
		if(v == player) {
			return i
		}
	}

	return -1
}
/*

func main() {

	team := NewTeam()
	otherTeam := NewTeam()

	fmt.Printf("Original\nTeam 1: %x\nTeam 2: %x\n\n", team, otherTeam)

	player, _:= (&team).NewPlayer("Jason")


	fmt.Printf("After add player\nTeam 1: %x\nTeam 2: %x\n\n", team, otherTeam)

	SwitchTeam(player, &otherTeam)


	fmt.Printf("After switch player\nTeam 1: %x\nTeam 2: %x\n\n", team, otherTeam)

	//fmt.Println(team)
	//fmt.Println(player)
}*/