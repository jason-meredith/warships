package base26

import (
	"math"
	"unicode"
)

/*********************************************************
 *														 *
 *                   CST8333 Exercise 3					 *
 *					   Jason Meredith					 *
 *														 *
 *	DATE:		December 2, 2018						 *
 *	FILE: 		base26.go								 *
 *	PURPOSE:	Converts base10 integer values into 	 *
 *				base26 string representation of numbers, *
 *				and the reverse, base26 strings back to	 *
 *				base10 ints								 *
 *														 *
 *														 *
 *********************************************************/

// ToLetter converts an int to a single letter character
func ToLetter(input int) rune {
	return rune(input + 'A')
}

// ToNumber converts a single letter character into a number
func ToNumber(input rune) int {
	unicode.ToUpper(input)
	return int(unicode.ToUpper(input) - 'A')
}

// ConvertToBase26 takes a base10 decimal integer and converts it to a base26 number as a string
func ConvertToBase26(input int) string {

	results := make([]rune, 0)

	var remainder int

	if input == 0 {
		return "A"
	}

	for input > 0 {
		remainder = input % 26
		input = input / 26
		results = append(results, ToLetter(remainder))
	}

	for i, j := 0, len(results)-1; i < j; i, j = i+1, j-1 {
		results[i], results[j] = results[j], results[i]
	}

	return string(results)
}

// ConvertToDecimal takes a base26 number as a string and converts it to base10 decimal integer
func ConvertToDecimal(input string) int {

	total := 0

	for i := 0; i < len(input); i++ {
		total += ToNumber(rune(input[i])) * int(math.Pow(26.0, float64(len(input) -1 - i)))
	}

	return total
}