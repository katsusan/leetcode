package code

import (
	"cmp"
	"slices"
	"strings"
)

/*
You are given two strings order and s. All the characters of order are unique and were sorted in some custom order previously.

Permute the characters of s so that they match the order that order was sorted. More specifically, if a character x occurs before a character y in order, then x should occur before y in the permuted string.

Return any permutation of s that satisfies this property.
*/

func customSortString(order string, s string) string {
	b := []byte(s)

	slices.SortFunc(b, func(a, b byte) int {
		aIdx := strings.IndexByte(order, a)
		bIdx := strings.IndexByte(order, b)

		return cmp.Compare(aIdx, bIdx)
	})

	return string(b)
}
