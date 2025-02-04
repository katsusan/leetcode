package code

func maxAscendingSum(nums []int) int {
	masm := 0
	for i := 0; i < len(nums); {
		s := nums[i]
		j := i + 1
		for j < len(nums) {
			if nums[j] <= nums[j-1] {
				break
			}
			s += nums[j]
			j++
		}
		masm = max(masm, s)
		i = j
	}
	return masm
}
