package main

func maxProfitMulti(prices []int) int {
	max := 0
	reg := 0
	up := 0
	for i := 1; i < len(prices); i++ {
		up = prices[i] - prices[i-1]
		if up > 0 {
			reg += up
		}
		if up <= 0 {
			max += reg
			reg = 0
		}

	}

	if up > 0 {
		max += reg
	}

	return max

}
