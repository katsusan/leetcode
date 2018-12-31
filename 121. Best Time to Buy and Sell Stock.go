package main

import "fmt"

//use divide and conqure, O(nlogn)
func maxProfit(prices []int) int {
	length := len(prices)
	if length < 2 {
		return 0
	}
	if length == 2 {
		if prices[1]-prices[0] > 0 {
			return prices[1] - prices[0]
		} else {
			return 0
		}
	}

	mid := (length + 1) / 2
	maxleft := maxProfit(prices[:mid])
	maxright := maxProfit(prices[mid:])

	minl := prices[0]
	maxr := prices[length-1]
	for i := 0; i < mid; i++ {
		if prices[i] < minl {
			minl = prices[i]
		}
	}
	for j := length - 1; j >= mid; j-- {
		if prices[j] > maxr {
			maxr = prices[j]
		}
	}

	fmt.Println(maxleft, maxright, maxr-minl)

	if maxleft >= maxright && maxleft >= maxr-minl {
		fmt.Println(prices, maxleft)
		return maxleft
	}

	if maxr-minl >= maxright && maxr-minl >= maxleft {
		fmt.Println(prices, maxr-minl)
		return maxr - minl
	}

	fmt.Println(prices, maxright)
	return maxright
}

//O(n) solution.
func maxProfit1(prices []int) int {
	max := 0
	seg := 0

	for i := 1; i <= len(prices)-1; i++ {
		up := prices[i] - prices[i-1]
		seg = seg + up
		if seg < 0 {
			seg = 0
		}

		if seg > max {
			max = seg
		}
	}
	return max

}
