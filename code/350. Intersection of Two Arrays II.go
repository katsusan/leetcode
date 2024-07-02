package code

func intersect350(nums1 []int, nums2 []int) []int {
	m1 := freqs(nums1)
	m2 := freqs(nums2)
	r := []int{}

	for n := range m1 {
		for range min(m1[n], m2[n]) {
			r = append(r, n)
		}
	}
	return r
}

func freqs(nums []int) map[int]int {
	m := make(map[int]int, len(nums))
	for _, v := range nums {
		m[v]++
	}
	return m
}
