package code

/*
Given two integer arrays nums1 and nums2, return an array of their intersection.
Each element in the result must be unique and you may return the result in any order.
*/

func intersection(nums1 []int, nums2 []int) []int {
	ma := make(map[int]bool, len(nums1)/2)
	arr := make([]int, 0, 8)
	for i := range nums1 {
		ma[nums1[i]] = false
	}

	for j := range nums2 {
		if added, ok := ma[nums2[j]]; ok && !added {
			ma[nums2[j]] = true
			arr = append(arr, nums2[j])
		}
	}

	return arr
}
