package game

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"warships/base26"
)

/*********************************************************
 *														 *
 *                   	  Warships						 *
 *					   Jason Meredith					 *
 *														 *
 *	DATE:		October 22, 2018						 *
 *	FILE: 		game.go								 	 *
 *	PURPOSE:	Contains all the actual game logic, all	 *
 *				networking and player-management issues	 *
 *				aside. Calculating hits/misses, firing	 *
 *				shots, determining Ship health is all	 *
 *				found here.								 *
 *														 *
 *				 										 *
 *														 *
 *********************************************************/

// Point values
const (
	// Points awarded for scoring hit
	HIT_POINT 			= 8
	HIT_STREAK_POINT 	= 10
	SINK_POINT			= 18
	DISCOVERY_POINT		= 1

	DEPLOY_COST			= 50
	DEPLOY_PER_SQ_COST	= 10
	MOVE_COST			= 1
	NEW_TEAM_COST		= 200
)

// Map icons
const (
	// ICON_ALIVE is the icon represent a portion of a ship that has not been hit
	ICON_ALIVE = '+'

	// ICON_DEAD is the icon represent a portion of a ship that has been hit
	ICON_DEAD = '@'

	// ICON_MISS is the icon representing a missed attempted attack on a team map
	ICON_MISS = "-"

	ICON_HIT = "x"
)

// Orientation is integer used to represent the Orientation enum options. Represents
// the direction a Ship is pointing
type Orientation uint8

// ShotResult is integer used to represent the ShotResult enum options. Represents the
// result of fired shot
type ShotResult uint8

// Orientation is the direction a Ship is pointing towards
const (
	VERTICAL Orientation = iota
	HORIZONTAL
)

// ShotResult is the result of a shot, a Hit if it hits a Ship, a miss if it doesn't and
// a Sink if it was a killing Hit
const (
	SINK ShotResult = iota
	HIT
	MISS
	REPEAT_HIT
)

// Target is the human-way of representing a square most similar to the board game (A1 -> Z26)
type Target struct {
	X string
	Y uint8
}

// Coordinate is the actually row/column square
type Coordinate struct {
	X uint8
	Y uint8
}

// Game represents a running with all Game settings and teams
type Game struct {

	// As long as this is true server will keep running
	Live			bool

	// Game server info
	Port			uint16
	Password		string
	StartTime		time.Time
	AdminPassword	string

	// The maximum ratio as 1:X that teams can be unbalanced
	// before Players from a short-handed Team can no longer
	// switch to a loaded Team

	MaxPlayers		uint8
	ShipLimit		uint8
	BoardSize		uint8

	Teams			[]*Team

	StartDeployPts	int

}

// Ship represents a single ship
type Ship struct {
	Team 		*Team

	Size        uint8
	Orientation Orientation
	Location    Coordinate


	Health		uint8 // Bit-field representing spots hit on this Ship

}

func StringToTarget(targetString string) (Target, error) {

	var re = regexp.MustCompile(`(?m)^([A-Z]+)([0-9]+)$`)

	search := re.FindAllStringSubmatch(strings.ToUpper(targetString), -1)

	if len(search) == 0  {
		return Target{}, errors.New("invalid Target Format (^[A-Z]+[0-9]+$)")
	}

	tokens := search[0]

	x := tokens[1]
	y, _ := strconv.Atoi(tokens[2])

	return Target{
		X: x,
		Y: uint8(y),
	}, nil
}

// ToCoordinate converts a Target (base64/number pair ex: B12) to a Coordinate (X,Y pair)
func (target Target) ToCoordinate() Coordinate {
	return Coordinate{ uint8(base26.ConvertToDecimal(target.X)), target.Y}
}

func (coordinate Coordinate) ToTarget() Target {
	return Target{base26.ConvertToBase26(int(coordinate.X)), coordinate.Y}
}


// FireShot handles a Player firing a shot on another Team. Returns either
// MISS, HIT or SINK
func FireShot(player *Player, targetTeam *Team, target Target) ShotResult {

	// Translate Target to an integer-pair Coordinate
	coordinate := target.ToCoordinate()

	// If there is a Ship at the given Coordinates, CheckLocation will
	// return the Ship that is there; otherwise it will return nil.
	enemyShip := CheckLocation(targetTeam, coordinate)

	// Add the shot the target teams firedupon list
	targetTeam.ShotsUpon = append(targetTeam.ShotsUpon, coordinate)

	// If there is no Ship there, return MISS, otherwise mark a HIT on the Ship
	// and return what hit() returns (HIT or SINK)
	if enemyShip == nil {
		player.HitStreak = 0
		player.Team.Misses[targetTeam] = append(player.Team.Misses[targetTeam], coordinate)
		return MISS
	} else {
		player.Team.Hits[targetTeam] = append(player.Team.Hits[targetTeam], coordinate)
		return enemyShip.Hit(player, coordinate)
	}

}

// getOccupyingSpaces returns an array of Coordinates that are occupied by this Ship
func (ship Ship) GetOccupyingSpaces() []Coordinate {

	var coordinateArray []Coordinate

	if ship.Orientation == VERTICAL {
		for y := uint8(0) ; y < ship.Size; y++ {
			coordinateArray = append(coordinateArray, Coordinate{ ship.Location.X, ship.Location.Y + y})
		}
	} else {
		for x := uint8(0) ; x < ship.Size; x++ {
			coordinateArray = append(coordinateArray, Coordinate{ ship.Location.X + x, ship.Location.Y})
		}
	}

	return coordinateArray
}

// CheckLocation will loop through each Ship on the target Team and see if any Ships occupy the target
// space
func CheckLocation(targetTeam *Team, target Coordinate) *Ship {

	// Loop through enemy Ships
	for _, ship := range targetTeam.Ships {
		// Loop through each Ship's occupying spaces
		for _, coordinate := range ship.GetOccupyingSpaces() {
			if target == coordinate {
				return ship
			}
		}

	}

	return nil
}

// ProduceHitBitmask creates a bitfield of all ones, except for a single zero at the offset position.
// This bitmask & ship.health will zero the bit at the offset.
func ProduceHitBitmask(offset uint8) uint8 {
	return uint8(255 - int(math.Pow(2, float64(8-offset-1))))
}

// GetOffset returns the number of squares a coordinate is away from the "location" of this
// ship
func (ship *Ship) GetOffset(coordinate Coordinate) uint8 {
	var offset uint8

	if ship.Orientation == VERTICAL {
		// How many square from Location was the hit
		offset = coordinate.Y - ship.Location.Y
	} else {
		offset = coordinate.X - ship.Location.X
	}

	return offset
}

// Hit marks a hit on the Ship and returns either HIT if the Ship is still alive or SINK if dead. If the hit occurred
// in a spot on the Ship take has already taken damage, the ship health will remain the same and we return REPEAT_HIT
func (ship *Ship) Hit(player *Player, coordinate Coordinate) ShotResult {

	originalShipHealth := ship.Health

	offset := ship.GetOffset(coordinate)

	bitmask :=  ProduceHitBitmask(offset)

	ship.Health = ship.Health & bitmask

	if ship.Health == originalShipHealth {
		if player != nil {
			player.HitStreak = 0
		}
		return REPEAT_HIT
	}

	if player != nil {
		if player.HitStreak == 0 {
			player.Points += HIT_POINT
		} else {
			player.Points += HIT_STREAK_POINT
		}
		player.HitStreak++
	}

	if ship.Health == 0  {
		if player != nil {
			player.Points += SINK_POINT
		}
		return SINK
	}

	return HIT
}

// NewShip generates a new Ship and adds it to a Team, then returns a pointer to the new Ship
func (team *Team) NewShip(size uint8, orientation Orientation, coordinate Coordinate) (*Ship, error) {

	// Make sure ship limit hasn't been reached
	if len(team.Ships) >= int(team.Game.ShipLimit) &&  int(team.Game.ShipLimit) != 0 {
		return nil, errors.New("ship Limit For This Team Has Been Reached")
	}


	// Make sure not out of bounds of Game area
	if orientation == HORIZONTAL && coordinate.X > (team.Game.BoardSize - size) {
		return nil, errors.New("ship being placed outside horizontal bound")
	}
	if orientation == VERTICAL && coordinate.Y > (team.Game.BoardSize - size) {
		return nil, errors.New("ship being placed outside vertical bound")
	}

	// Make sure no overlaps occur
	if orientation == HORIZONTAL {
		for x := uint8(0); x < size; x++ {
			if CheckLocation(team, Coordinate{x + coordinate.X, coordinate.Y}) != nil {
				return nil, errors.New("ship overlap, cannot place Ship here")
			}
		}
	}

	if orientation == VERTICAL {
		for y := uint8(0); y < size; y++ {
			if CheckLocation(team, Coordinate{coordinate.X, y + coordinate.Y}) != nil {
				return nil, errors.New("ship overlap, cannot place Ship here")
			}
		}
	}

	// Create the ship
	ship := Ship{
		team,
		size,
		orientation,
		coordinate,

		// Bit field of 1s the length of the Ship Size ie Ship Size 4 -> 11110000, 2 -> 11000000
		GetHealthBitfield(size),

	}

	team.Ships = append(team.Ships, &ship)


	return &ship, nil
}

// GetHealthBitfield returns a bitfield representing the Ship and thats parts of it that are hit and unscathed.
func GetHealthBitfield(size uint8) uint8 {
	return (uint8(math.Pow(2, float64(size))) - 1) << (8 - size)
}


// BoardCoordinates creates an iterator that iterates over ever Coordinate of the Board
func (game *Game) BoardCoordinates() <-chan Coordinate {

	channel := make(chan Coordinate)

	go func() {

		var x, y uint8

		// Iterate over each grid pair the
		for x = 0; x < game.BoardSize; x ++ {
			for y = 0; y < game.BoardSize; y ++ {
				channel <- Coordinate{x, y}
			}
		}

		close(channel)
	}()

	return channel
}

// ShipIcon returns, given a Ship and a coordinate, what icon should be
// displayed at this  coordinate?
func (ship *Ship) ShipIcon(coordinate Coordinate) rune {
	health := ship.Health
	offset := ship.GetOffset(coordinate)

	var icon rune

	// ProduceHitBitmask will produced a binary of all 1s except at the offset
	// Ex: 3 -> 11111011
	// I will use the complement (0000 0100) to see if there is a 1 or a 0 at
	// that position in the Health value (1111 1111)&(0000 0100) = that spot
	// is alive, (1111 1011)&(0000 0100) = that spot is dead
	if health & ^ProduceHitBitmask(offset) > 0 {
		icon = ICON_ALIVE
	} else {
		icon = ICON_DEAD
	}

	return icon
}

type ShipCoord struct {
	ship 	*Ship
	coord 	Coordinate
	icon 	rune
}

// ShipCoordinates creates an iterator that iterates through the Coordinates
// of all the ship on a Team. It uses the channel iterator pattern, meaning
// it creates a Goroutine that loops through all the values and pushes them
// into a channel to be consumed in the main thread
func (team *Team) ShipCoordinates() chan ShipCoord {
	channel := make(chan ShipCoord)

	go func() {
		defer close(channel)

		// Loop over each ship the belongs to this team
		for _, ship := range team.Ships {
			// Loop over each coordinate of each ship
			for _, coord := range ship.GetOccupyingSpaces() {

				shipCoord := ShipCoord{
					ship: ship,
					coord: coord,
					icon: ship.ShipIcon(coord),
				}

				channel <- shipCoord
			}
		}

	}()

	return channel
}

func (game *Game) UniqueTeamName(name string) bool {
	for _, team := range game.Teams {
		if team.Name == name {
			return false
		}
	}

	return true
}



func (game *Game) GetRadar(team *Team, targetTeam *Team) [][]string {

	boardSize := game.BoardSize


	board := make([][] string, boardSize)
	for x := 0; x < int(boardSize); x++ {
		board[x] = make([] string, boardSize)
	}

	// Initialize grid with spaces
	for coord := range game.BoardCoordinates() {
		board[coord.X][coord.Y] = "_|"
	}

	// Add in hits
	for _, hit := range team.Hits[targetTeam] {
		board[hit.X][hit.Y] = ICON_HIT + "|"
	}

	// Add in misses
	for _, miss := range team.Misses[targetTeam] {
		board[miss.X][miss.Y] = ICON_MISS + "|"
	}


	return board
}

func (game *Game) GetMap(team *Team) [][]string {

	boardSize := game.BoardSize


	board := make([][] string, boardSize)
	for x := 0; x < int(boardSize); x++ {
		board[x] = make([] string, boardSize)
	}

	// Initialize grid with spaces
	for coord := range game.BoardCoordinates() {
		board[coord.X][coord.Y] = "_|"
	}

	// Add in shots upon our team
	for _, miss := range team.ShotsUpon {
		board[miss.X][miss.Y] = ICON_MISS + "|"
	}

	// Add in our ships
	for shipCoord := range team.ShipCoordinates() {
		board[shipCoord.coord.X][shipCoord.coord.Y] = string(shipCoord.icon) + "|"
	}


	return board

}