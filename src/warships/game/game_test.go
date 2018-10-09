package game

import (
	"testing"
)

func TestTeam_NewShip(t *testing.T) {
	game := Game{[]*Team{}, 32, 128}

	team := game.NewTeam()

	// New Ship details
	size := uint8(5);
	orientation := VERTICAL;
	location := Coordinate{ 10, 10}

	testShip, err := team.NewShip(size, orientation, location);

	// Test new Ship
	if err != nil {
		t.Error("Error Thrown: ", err);
	}

	if testShip.team != &team {
		t.Error("Ship Team not set");
	}

	if testShip.location != location {
		t.Error("Ship Location not set");
	}

	if testShip.orientation != VERTICAL {
		t.Error("Ship Orientation not set");
	}

	if testShip.size != size {
		t.Error("Ship Size not set");
	}

	t.Run("Error Check - Boundaries", func(t *testing.T) {
		_, err := team.NewShip(size, HORIZONTAL, Coordinate{125, 0})
		if err == nil {
			t.Error("Creating Ship out of x bounds should have returned error");
		}

		_, err = team.NewShip(size, VERTICAL, Coordinate{0, 125})
		if err == nil {
			t.Error("Creating Ship out of y bounds should have returned error");
		}

	})

	t.Run("Error Check - Overlap", func(t *testing.T) {
		_, err := team.NewShip(size, HORIZONTAL, Coordinate{50, 50})
		_, err = team.NewShip(size, VERTICAL, Coordinate{51, 49})
		if err == nil {
			t.Error("Creating Ship that intersects with existing ship should result in error");
		}


	})
}

func TestShip_GetOccupyingSpaces(t *testing.T) {
	game := Game{[]*Team{}, 32, 128}

	team := game.NewTeam()

	// New Ship details
	size := uint8(5);
	orientation := VERTICAL;
	location := Coordinate{ 10, 10}

	testShip, _ := team.NewShip(size, orientation, location);

	expectedResults := []Coordinate{
		{10, 10},
		{10, 11},
		{10, 12},
		{10, 13},
		{10, 14},
	}

	matches := true

	for i, actual := range testShip.GetOccupyingSpaces() {
		if expectedResults[i] != actual {
			matches = false;
		}
	}

	if matches == false {
		t.Error("Does not match expected results")
	}
}

func TestCheckLocation(t *testing.T) {
	game := Game{[]*Team{}, 32, 128}

	team := game.NewTeam()

	// New Ship details
	size := uint8(5);
	orientation := VERTICAL;
	location := Coordinate{ 10, 10}

	testShip, _ := team.NewShip(size, orientation, location);

	ship := CheckLocation(&team, Coordinate{9,9})

	if ship != nil {
		t.Error("Found ship where no Ship should have existed");
	}

	ship = CheckLocation(&team, Coordinate{10,10})

	if ship != testShip {
		t.Error("Didn't find Ship where one should have existed");
	}
}

func TestGetHealthBitfield(t *testing.T) {

	result := GetHealthBitfield(2)

	// Size 2 should result in bit-field 11000000 = 192
	if result != 192 {
		t.Error("Bitfield returning proper value for size 2")
	}

	result = GetHealthBitfield(4)

	// Size 4 should result in bit-field 11110000 = 240
	if result != 240 {
		t.Error("Bitfield returning proper value")
	}

	result = GetHealthBitfield(6)

	// Size 6 should result in bit-field 11111100 = 252
	if result != 252 {
		t.Error("Bitfield returning proper value")
	}
}


func TestShip_Hit(t *testing.T) {
	game := Game{[]*Team{}, 32, 128}

	team := game.NewTeam()

	// New Ship details
	size := uint8(5);
	orientation := VERTICAL;
	location := Coordinate{ 0, 0}


	testShip, _ := team.NewShip(size, orientation, location);

	// Bitmask: 0111 1111 & 1111 1000 -> 0111 1000 (120)
	result := testShip.Hit(Coordinate{0, 0})
	if testShip.health != 120 {
		t.Error("Health not depleted as expected")
	}

	if result != HIT {
		t.Error("Hit does not return HIT")
	}

	// Bitmask: 1011 1111 & 0111 1000 -> 0011 1000 (56)
	testShip.Hit(Coordinate{0, 1})
	if testShip.health != 56 {
		t.Error("Health not depleted as expected")
	}

	// Bitmask: 1101 1111 & 0011 1000 -> 0001 1000 (24)
	testShip.Hit(Coordinate{0, 2})
	if testShip.health != 24 {
		t.Error("Health not depleted as expected")
	}

	// Bitmask: 1110 1111 & 0001 1000 -> 0000 1000 (8)
	testShip.Hit(Coordinate{0, 3})
	if testShip.health != 8 {
		t.Error("Health not depleted as expected")
	}

	// Bitmask: 1111 0111 & 0000 0000 -> 0000 0000 (0)
	result = testShip.Hit(Coordinate{0, 4})
	if testShip.health != 0  {
		t.Error("Health not depleted as expected")
	}

	if result != SINK {
		t.Error("Killing HIT does not return SINK")
	}

}

func TestProduceHitBitmask(t *testing.T) {
	properBitmask := true

	// Offset 0 should return 0111 1111 (127)
	if ProduceHitBitmask(0 ) != 127 {
		properBitmask = false;
	}

	// Offset 1 should return 1011 1111 (191)
	if ProduceHitBitmask(1 ) != 191 {
		properBitmask = false;
	}

	// Offset 2 should return 1101 1111 (223)
	if ProduceHitBitmask(2 ) != 223 {
		properBitmask = false;
	}


	// Offset 7 should return 1111 1101 (253)
	if ProduceHitBitmask(6 ) != 253 {
		properBitmask = false;
	}

	// Offset 7 should return 1111 1110 (254)
	if ProduceHitBitmask(7 ) != 254 {
		properBitmask = false;
	}

	if !properBitmask {
		t.Error("Not producing proper bitmask")
	}
}