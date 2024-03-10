package code

/*
Given an array of integers arr, return true if the number of occurrences
of each value in the array is unique or false otherwise.
*/

func uniqueOccurrences(arr []int) bool {
	occurMap := make(map[int]int)
	for _, v := range arr {
		occurMap[v]++
	}

	dupMap := make(map[int]bool, len(occurMap)/2)
	for _, cnt := range occurMap {
		if dupMap[cnt] {
			return false
		}
		dupMap[cnt] = true
	}

	return true
}
