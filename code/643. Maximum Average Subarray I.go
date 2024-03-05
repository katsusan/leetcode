package code

/*
You are given an integer array nums consisting of n elements, and an integer k.

Find a contiguous subarray whose length is equal to k that has the maximum average value
and return this value. Any answer with a calculation error less than 10-5 will be accepted.
*/

func findMaxAverage(nums []int, k int) float64 {
	if len(nums) <= k {
		return float64(sum(nums)) / float64(k)
	}
	s := sum(nums[0:k])
	maxsum := s
	for i := 1; i <= len(nums)-k; i++ {
		s = s - nums[i-1] + nums[i+k-1]
		maxsum = max(maxsum, s)
	}
	return float64(maxsum) / float64(k)
}

func sum(nums []int) int {
	s := 0
	for _, v := range nums {
		s += v
	}
	return s
}
