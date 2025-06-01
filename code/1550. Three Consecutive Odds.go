package code

func threeConsecutiveOdds(arr []int) bool {
	odds := 0
	for i := 0; i < len(arr); i++ {
		if arr[i]&1 == 1 {
			odds++
		} else {
			odds = 0
		}

		if odds == 3 {
			return true
		}
	}

	return false
}
