package code

func missingNumber(nums []int) int {
	// rationale: 0 ^ N = N; a ^ b ^ c = a ^ c ^ b
	d := 0
	i := 0
	for i = 0; i < len(nums); i++ {
		d = d ^ i ^ nums[i]
	}
	return d ^ i
}
