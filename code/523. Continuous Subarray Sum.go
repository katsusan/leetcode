package code

func checkSubarraySum(nums []int, k int) bool {
	if len(nums) <= 1 {
		return false
	}

	presum := make([]int, len(nums))
	presum[0] = nums[0]

	for i := 1; i < len(nums); i++ {
		presum[i] = presum[i-1] + nums[i]
	}

	r := make([]int, len(presum))
	lastIdx := make(map[int]int)

	for i := range r {
		r[i] = presum[i] % k
		if r[i] == 0 && i > 0 { // n*k && length>=2
			return true
		}

		if _, ok := lastIdx[r[i]]; !ok {
			lastIdx[r[i]] = i
		} else {
			// found n*k circle
			if i-lastIdx[r[i]] >= 2 {
				// length is at least 2
				return true
			}
		}
	}

	return false
}
