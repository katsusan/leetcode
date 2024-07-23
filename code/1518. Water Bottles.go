package code

func numWaterBottles(numBottles int, numExchange int) int {
	drunk := numBottles
	full, empty := 0, numBottles
	for empty >= numExchange {
		// exchange
		full = empty / numExchange
		empty -= full * numExchange

		// drink
		drunk += full
		empty += full
	}

	return drunk
}
