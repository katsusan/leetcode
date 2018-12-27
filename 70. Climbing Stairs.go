package main

import "fmt"

//use cache

var cachemap map[int]int = make(map[int]int)

func climbStairs(n int) int {

	if r, found := cachemap[n]; found {
		return r
	}

	switch n {
	case 0:
		cachemap[0] = 0
		return 0
	case 1:
		cachemap[1] = 1
		return 1
	case 2:
		cachemap[2] = 2
		return 2
	default:
		cachemap[n] = climbStairs(n-1) + climbStairs(n-2)
		return cachemap[n]
	}
}

//use DP
func climbStairsDP(n int) int {

	if n < 4 {
		return n
	}

	dparr := make([]int, n+1)
	copy(dparr[0:4], []int{0, 1, 2, 3})

	for i := 4; i <= n; i++ {
		dparr[i] = dparr[i-1] + dparr[i-2]
	}
	fmt.Println(dparr)
	return dparr[n]

}
