package code

func canBeEqual(target []int, arr []int) bool {
	marks := make([]int16, 1001)
	for i := range target {
		marks[target[i]]++
		marks[arr[i]]--
	}

	for _, v := range marks {
		if v != 0 {
			return false
		}
	}

	return true
}
