package game

import (
	"errors"
	"math"
)

type Orientation uint8
type ShotResult uint8

const (
	VERTICAL Orientation = iota
	HORIZONTAL
)

const (
	SINK ShotResult = iota
	HIT
	MISS
)

type Target struct {
	x				rune
	y 				uint8
}

type Coordinate struct {
	x 				uint8
	y 				uint8
}

type Game struct {

	teams			[]*Team

	shipLimit		uint8

	boardSize		uint8


}

type Ship struct {
	team 		*Team

	size        uint8
	orientation Orientation
	location    Coordinate


	health		uint8			// Bit-field representing spots hit on this Ship

}

// ToCoordinate converts a Target (base64/number pair ex: B12) to a Coordinate (x,y pair)
// TODO: Targets are currently Rune/Int, limited x axis to 26. Make this a base 26 number (A-Z, AA, AB, AC, etc)
func (target Target) ToCoordinate() Coordinate {
	x := uint8(target.x - 'A')

	return Coordinate{ x, target.y }
}


// FireShot handles a Player firing a shot on another Team. Returns either
// MISS, HIT or SINK
func FireShot(player *Player, targetTeam *Team, target Target) ShotResult {

	// Translate Target to an integer-pair Coordinate
	coordinate := target.ToCoordinate();

	// If there is a Ship at the given Coordinates, CheckLocation will
	// return the Ship that is there; otherwise it will return nil.
	enemyShip := CheckLocation(targetTeam, coordinate);

	// If there is no Ship there, return MISS, otherwise mark a HIT on the Ship
	// and return what hit() returns (HIT or SINK)
	if enemyShip == nil {
		return MISS;
	} else {
		return enemyShip.Hit(coordinate);
	}

}

// getOccupyingSpaces returns an array of Coordinates that are occupied by this Ship
func (ship Ship) GetOccupyingSpaces() []Coordinate {

	coordinateArray := []Coordinate{};

	if ship.orientation == VERTICAL {
		for y := uint8(0) ; y < ship.size; y++ {
			coordinateArray = append(coordinateArray, Coordinate{ ship.location.x, ship.location.y + y})
		}
	} else {
		for x := uint8(0) ; x < ship.size; x++ {
			coordinateArray = append(coordinateArray, Coordinate{ ship.location.x + x, ship.location.y})
		}
	}

	return coordinateArray
}

// CheckLocation will loop through each Ship on the target Team and see if any Ships occupy the target
// space
func CheckLocation(targetTeam *Team, target Coordinate) *Ship {

	// Loop through enemy Ships
	for _, ship := range targetTeam.ships {
		// Loop through each Ship's occupying spaces
		for _, coordinate := range ship.GetOccupyingSpaces() {
			if target == coordinate {
				return ship
			}
		}

	}

	return nil
}

func ProduceHitBitmask(offset uint8) uint8 {
	return uint8(255 - int(math.Pow(2, float64(8-offset-1))))
}

// hit marks a hit on the Ship and returns either HIT if the Ship is still alive or SINK if dead
func (ship *Ship) Hit(coordinate Coordinate) ShotResult {

	var offset uint8

	if ship.orientation == VERTICAL {
		// How many square from location was the hit
		offset = coordinate.y - ship.location.y;
	} else {
		offset = coordinate.x - ship.location.x;
	}

	bitmask :=  ProduceHitBitmask(offset)

	ship.health = ship.health & bitmask

	if ship.health == 0  {
		return SINK
	}

	return HIT
}

// NewShip generates a new Ship and adds it to a Team, then returns a pointer to the new Ship
func (team *Team) NewShip(size uint8, orientation Orientation, coordinate Coordinate) (*Ship, error) {

	// Make sure ship limit hasn't been reached
	if len(team.ships) >= int(team.game.shipLimit) {
		return nil, errors.New("Ship Limit For This Team Has Been Reached");
	}


	// Make sure not out of bounds of Game area
	if orientation == HORIZONTAL && coordinate.x > (team.game.boardSize - size) {
		return nil, errors.New("ship being placed outside horizontal bound");
	}
	if orientation == VERTICAL && coordinate.y > (team.game.boardSize - size) {
		return nil, errors.New("ship being placed outside vertical bound");
	}

	// Make sure no overlaps occur
	if orientation == HORIZONTAL {
		for x := uint8(0); x < size; x++ {
			if CheckLocation(team, Coordinate{x + coordinate.x, coordinate.y}) != nil {
				return nil, errors.New("ship overlap, cannot place Ship here");
			}
		}
	}

	if orientation == VERTICAL {
		for y := uint8(0); y < size; y++ {
			if CheckLocation(team, Coordinate{coordinate.x, y + coordinate.y}) != nil {
				return nil, errors.New("ship overlap, cannot place Ship here");
			}
		}
	}

	// Create the ship
	ship := Ship{
		team,
		size,
		orientation,
		coordinate,

		// Bit field of 1s the length of the Ship size ie Ship size 4 -> 11110000, 2 -> 11000000
		GetHealthBitfield(size),

	}

	team.ships = append(team.ships, &ship);


	return &ship, nil;
}

// GetHealthBitfield returns a bitfield representing the Ship and thats parts of it that are hit and unscathed.
func GetHealthBitfield(size uint8) uint8 {
	return (uint8(math.Pow(2, float64(size))) - 1) << (8 - size);
}