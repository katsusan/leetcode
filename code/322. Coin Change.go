package code

func coinChange(coins []int, amount int) int {
	return coinChangeWithCache(coins, amount, make(map[int]int))
}

func coinChangeWithCache(coins []int, amount int, cache map[int]int) int {
	if amount < 0 {
		return -1
	} else if amount == 0 {
		return 0
	}

	var min = -2
	for _, coin := range coins {
		var s int
		if v, ok := cache[amount-coin]; ok {
			s = v
		} else {
			s = coinChangeWithCache(coins, amount-coin, cache)
			cache[amount-coin] = s
		}

		if s != -1 {
			if min < 0 {
				min = s
			} else if s < min {
				min = s
			}
		}
	}
	return min + 1
}
