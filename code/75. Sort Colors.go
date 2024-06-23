package code

func sortColors(nums []int) {
	// O(N) time, O(1) space

	cnts := make([]int, 3) // stores counts of 0s, 1s, 2s
	for i := range nums {
		cnts[nums[i]]++
	}

	for i := 0; i < cnts[0]; i++ {
		nums[i] = 0
	}

	for j := cnts[0]; j < cnts[0]+cnts[1]; j++ {
		nums[j] = 1
	}

	for k := cnts[0] + cnts[1]; k < cnts[0]+cnts[1]+cnts[2]; k++ {
		nums[k] = 2
	}
}
