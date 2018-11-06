package base26

import (
	"math"
)

// ToLetter converts an int to a single letter character
func ToLetter(input int) rune {
	return rune(input + 'A' - 1)
}

// ToNumber converts a single letter character into a number
func ToNumber(input rune) int {
	return int(input - 'A')
}

// ConvertToBase26
func ConvertToBase26(input int) string {

	results := make([]rune, 0)

	var remainder int

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

func ConvertToDecimal(input string) int {

	total := 0

	for i := 0; i < len(input); i++ {
		total += ToNumber(rune(input[i])) * int(math.Pow(26.0, float64(len(input) -1 - i)))
	}

	return total
}