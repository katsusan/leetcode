package code

func numSubarraysWithSum(nums []int, goal int) int {
	sum := make([]int, len(nums))
	sum[0] = nums[0]
	for i := 1; i < len(sum); i++ {
		sum[i] = sum[i-1] + nums[i]
	}

	cnt := 0
	m := make(map[int][]int)
	m[0] = []int{-1}
	for i := range sum {
		if v, ok := m[sum[i]-goal]; ok {
			cnt += len(v)
		}
		m[sum[i]] = append(m[sum[i]], i)
	}

	return cnt
}
