package code

func twoSum_167(numbers []int, target int) []int {
	i := 0
	j := len(numbers) - 1
	for i != j {
		sum := numbers[i] + numbers[j]
		if sum == target {
			return []int{i + 1, j + 1}
		} else if sum > target {
			j--
			continue
		} else if sum < target {
			i++
			continue
		}
	}
	return []int{}
}
