package code

import "slices"

func kidsWithCandies(candies []int, extraCandies int) []bool {
	maxCandy := slices.Max(candies)
	var cheat = make([]bool, len(candies))
	for i := range cheat {
		if candies[i]+extraCandies >= maxCandy {
			cheat[i] = true
		}
	}
	return cheat
}
