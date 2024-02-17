package code

import "slices"

/*
Given an array of integers arr and an integer k. Find the least number of unique integers after removing exactly k elements.

Eg.
Input: arr = [4,3,1,1,3,3,2], k = 3
Output: 2
Explanation: Remove 4, 2 and either one of the two 1s or three 3s. 1 and 3 will be left.

*/

func findLeastNumOfUniqueInts(arr []int, k int) int {
	var CntMap = make(map[int]int)
	for i := range arr {
		if cnt, ok := CntMap[arr[i]]; ok {
			CntMap[arr[i]] = cnt + 1
		} else {
			CntMap[arr[i]] = 1
		}
	}

	var CntArr = make([]int, len(CntMap))
	var idx = 0
	for _, cnt := range CntMap {
		CntArr[idx] = cnt
		idx++
	}
	slices.Sort(CntArr)

	remains := k
	idx = 0
	for idx < len(CntArr) && remains-CntArr[idx] >= 0 {
		remains = remains - CntArr[idx]
		idx++
	}
	//fmt.Println("CntArr=", CntArr, "idx=", idx)
	return len(CntArr) - idx
}
