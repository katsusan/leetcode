package code

/*
Given an array of integers nums, calculate the pivot index of this array.

The pivot index is the index where the sum of all the numbers strictly to the left of the index is equal to the sum of all the numbers strictly to the index's right.

If the index is on the left edge of the array, then the left sum is 0 because there are no elements to the left. This also applies to the right edge of the array.

Return the leftmost pivot index. If no such index exists, return -1.
*/

func pivotIndex(nums []int) int {
	preSum := make([]int, len(nums))
	sufSum := make([]int, len(nums))

	preSum[0] = 0
	for i := 1; i < len(nums); i++ {
		preSum[i] = preSum[i-1] + nums[i-1]
	}

	sufSum[len(nums)-1] = 0
	for j := len(nums) - 2; j >= 0; j-- {
		sufSum[j] = sufSum[j+1] + nums[j+1]
	}

	for i := 0; i < len(nums); i++ {
		if preSum[i] == sufSum[i] {
			return i
		}
	}

	// fmt.Println("preSum:", preSum)
	// fmt.Println("sufSum:", sufSum)
	return -1
}
