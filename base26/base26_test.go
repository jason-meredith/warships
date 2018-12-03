package base26

import (
	"testing"
)

func TestConvertToBase26(t *testing.T) {
	if ConvertToBase26(0) != "A" {
		t.Error("Convert Int (1) To Base26 (A) fail")
	}

	if ConvertToBase26(25) != "Z" {
		t.Error("Convert Int (25) To Base26 (Z) fail")
	}

	// Rollover from ... X, Y, Z to -> BA, BB, BC, BD, ...
	if ConvertToBase26(26) != "BA" {
		t.Error("Rollover Int (26) To Base26 (BA) fail")
	}

	if ConvertToBase26(27) != "BB" {
		t.Error("Convert Int (27) To Base26 (BB) fail")
	}

	// Rolling over from B* to C*
	if ConvertToBase26(51) != "BZ" {
		t.Error("Convert Int (51) To Base26 (BZ) fail")
	}

	if ConvertToBase26(52) != "CA" {
		t.Error("Convert Int (52) To Base26 (CA) fail")
	}

	if ConvertToBase26(53) != "CB" {
		t.Error("Convert Int (53) To Base26 (CB) fail")
	}

}

func TestConvertToDecimal(t *testing.T) {
	if ConvertToDecimal("A") != 0 {
		t.Error("Convert String (A) To Decimal (0) fail")
	}

	if ConvertToDecimal("Z") != 25 {
		t.Error("Convert String (Z) To Decimal (25) fail")
	}

	// Rollover from ... X, Y, Z to -> BA, BB, BC, BD, ...
	if ConvertToDecimal("BA") != 26 {
		t.Error("Rollover String (BA) To Decimal (26) fail")
	}

	if ConvertToDecimal("BB") != 27 {
		t.Error("Convert String (BB) To Decimal (27) fail")
	}

	// Rolling over from B* to C*
	if ConvertToDecimal("BZ") != 51 {
		t.Error("Convert String (BZ) To Decimal (51) fail")
	}

	if ConvertToDecimal("CA") != 52 {
		t.Error("Convert String (CA) To Decimal (52) fail")
	}

	if ConvertToDecimal("CB") != 53 {
		t.Error("Convert String (CB) To Decimal (53) fail")
	}

	if ConvertToDecimal("cb") != 53 {
		t.Error("Lowercase Convert String (CB) To Decimal (53) fail")
	}
}

func TestToLetter(t *testing.T) {
	if ToLetter(0) != 'A' {
		t.Error("Convert int (0) to letter 'A' fail ")
	}

	if ToLetter(25) != 'Z' {
		t.Error("Convert int (25) to letter 'Z' fail ")
	}
}

func TestToNumber(t *testing.T) {
	if ToNumber('A') != 0 {
		t.Error("Convert letter 'A' to int (0) fail")
	}

	if ToNumber('Z') != 25 {
		t.Error("Convert letter 'Z' to int (25) fail")
	}

	if ToNumber('a') != 0 {
		t.Error("Convert lowercase letter 'a' to int (0) fail")
	}

	if ToNumber('z') != 25 {
		t.Error("Convert letter lower case 'z' to int (25) fail")
	}
}