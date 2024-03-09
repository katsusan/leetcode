package code

func findMaxConsecutiveOnes(nums []int) int {
	i := 0
	maxOnes := 0
	for i < len(nums) {
		if nums[i] == 0 {
			i++
			continue
		}

		nextZero := findNextZero(nums, i)
		maxOnes = max(maxOnes, nextZero-i)
		i = nextZero
	}

	return maxOnes
}

func findNextZero(nums []int, cur int) int {
	i := cur
	for i < len(nums) {
		if nums[i] == 0 {
			return i
		}
		i++
	}
	return i
}
