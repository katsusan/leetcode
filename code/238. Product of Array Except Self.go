package code

/*
Given an integer array nums, return an array answer such that answer[i] is equal to the product of all the elements of nums except nums[i].

The product of any prefix or suffix of nums is guaranteed to fit in a 32-bit integer.

You must write an algorithm that runs in O(n) time and without using the division operation.

Input: nums = [1,2,3,4]
Output: [24,12,8,6]
*/

func productExceptSelf(nums []int) []int {
	pre := make([]int, len(nums))
	suf := make([]int, len(nums))

	pre[0] = 1
	for i := 1; i < len(pre); i++ {
		pre[i] = pre[i-1] * nums[i-1]
	}

	suf[len(suf)-1] = 1
	for j := len(suf) - 2; j >= 0; j-- {
		suf[j] = suf[j+1] * nums[j+1]
	}

	p := make([]int, len(nums))
	p[0] = suf[0]
	p[len(nums)-1] = pre[len(nums)-1]
	for i := 1; i < len(nums)-1; i++ {
		p[i] = pre[i] * suf[i]
	}
	return p
}
