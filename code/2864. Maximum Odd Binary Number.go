package code

import "strings"

/*
You are given a binary string s that contains at least one '1'.

You have to rearrange the bits in such a way that the resulting binary number is the maximum odd binary number that can be created from this combination.

Return a string representing the maximum odd binary number that can be created from the given combination.

Note that the resulting string can have leading zeros.
*/

func maximumOddBinaryNumber(s string) string {
	cntOne := strings.Count(s, "1")
	cntZero := len(s) - cntOne
	return strings.Repeat("1", cntOne-1) + strings.Repeat("0", cntZero) + "1"
}
