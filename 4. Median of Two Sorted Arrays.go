package main

import "fmt"

func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	tmparr := make([]int, len(nums1)+len(nums2))

	mergesort(tmparr, nums1, nums2)

	fmt.Println(tmparr)

	mid := (len(nums1) + len(nums2)) / 2
	if (len(nums1)+len(nums2))%2 != 0 {
		return float64(tmparr[mid])
	} else {
		return (float64(tmparr[mid]) + float64(tmparr[mid-1])) / 2
	}
}

func mergesort(dst []int, nums1 []int, nums2 []int) {
	if len(dst) < len(nums1)+len(nums2) {
		return
	}

	var k1, k2, x int
	for k1 < len(nums1) && k2 < len(nums2) {
		if nums1[k1] < nums2[k2] {
			dst[x] = nums1[k1]
			x++
			k1++
		} else {
			dst[x] = nums2[k2]
			x++
			k2++
		}
	}

	if k1 == len(nums1) {
		copy(dst[x:], nums2[k2:])
	} else if k2 == len(nums2) {
		copy(dst[x:], nums1[k1:])
	}
}
