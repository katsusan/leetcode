package code

import "sort"

func maxDistance(position []int, m int) int {
	sort.Ints(position)
	n := len(position)
	low, high := 1, (position[n-1]-position[0])/(m-1)
	ans := 1

	for low <= high {
		mid := low + (high-low)/2
		if canPlace(position, mid, m) {
			ans = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return ans
}

func canPlace(position []int, d, balls int) bool {
	cnt := 1
	lastP := position[0]
	for i := 1; i < len(position); i++ {
		if position[i]-lastP >= d {
			cnt++
			lastP = position[i]
		}

		if cnt >= balls {
			return true
		}
	}
	return false
}
