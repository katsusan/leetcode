package code

func divideArray(nums []int) bool {
	cnts := make([]int, 501)
	for _, v := range nums {
		cnts[v]++
	}

	for _, c := range cnts {
		if c%2 != 0 {
			return false
		}
	}

	return true
}
