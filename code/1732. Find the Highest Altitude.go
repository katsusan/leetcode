package code

func largestAltitude(gain []int) int {
	highest := 0
	a := make([]int, len(gain)+1)
	a[0] = 0
	for i := 1; i < len(a); i++ {
		a[i] = a[i-1] + gain[i-1]
		highest = max(highest, a[i])
	}
	return highest
}
