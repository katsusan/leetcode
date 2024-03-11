package code

import "slices"

/*
You are given an array nums consisting of positive integers.

Return the total frequencies of elements in nums such that those elements all have the maximum frequency.

The frequency of an element is the number of occurrences of that element in the array.
*/

func maxFrequencyElements(nums []int) int {
	m := make(map[int]int, len(nums)/2)

	for _, v := range nums {
		m[v]++
	}

	a := make([]int, 0, 8)
	for _, v := range m {
		a = append(a, v)
	}

	maxFreq := slices.Max(a)
	cnt := 0
	for _, c := range a {
		if c == maxFreq {
			cnt++
		}
	}

	return cnt * maxFreq
}
