package sort

//QuickSort : implement of  recursive quicksort
func QuickSort(nums []int) {
	if len(nums) < 2 {
		return
	}
	qsortex(nums, 0, len(nums)-1)
}

//specify the first number as the "flag"
//return the right postion of first number(after sorted)
//and ensure it in the right place
func qsortex(data []int, start, end int) {
	flag := data[start]
	pos := start
	i, j := start, end

	for i <= j {
		for j >= pos && data[j] >= flag {
			j--
		}

		if j >= pos {
			data[pos], data[j] = data[j], data[pos]
			pos = j
		}

		for i <= pos && data[i] <= flag {
			i++
		}

		if i <= pos {
			data[pos], data[i] = data[i], data[pos]
			pos = i
		}

		if pos-start > 1 {
			qsortex(data, start, pos-1)
		}
		if end-pos > 1 {
			qsortex(data, pos+1, end)
		}
	}
}
