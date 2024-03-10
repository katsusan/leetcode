package code

/*
Given two 0-indexed integer arrays nums1 and nums2, return a list answer of size 2 where:
	-answer[0] is a list of all distinct integers in nums1 which are not present in nums2.
	-answer[1] is a list of all distinct integers in nums2 which are not present in nums1.
*/

func findDifference(nums1 []int, nums2 []int) [][]int {
	var s1 SetGeneric[int] = make(map[int]struct{}, len(nums1)/2)
	var s2 SetGeneric[int] = make(map[int]struct{}, len(nums2)/2)

	for _, v := range nums1 {
		s1.Add(v)
	}

	for _, v := range nums2 {
		s2.Add(v)
	}

	ans1 := make([]int, 0, len(s1)/2)
	ans2 := make([]int, 0, len(s2)/2)

	for e := range s1 {
		if !s2.Exist(e) {
			ans1 = append(ans1, e)
		}
	}

	for e := range s2 {
		if !s1.Exist(e) {
			ans2 = append(ans2, e)
		}
	}

	return [][]int{ans1, ans2}
}

type SetGeneric[T int] map[T]struct{}

func (s SetGeneric[T]) Add(e T) {
	s[e] = struct{}{}
}

func (s SetGeneric[T]) Exist(e T) bool {
	if _, ok := s[e]; ok {
		return true
	}
	return false
}
