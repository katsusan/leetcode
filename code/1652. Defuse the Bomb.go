package code

import "slices"

func decrypt1652(code []int, k int) []int {
	switch {
	case k > 0:
		return doDecrypt(code, k)
	case k < 0:
		slices.Reverse(code)
		r := doDecrypt(code, -k)
		slices.Reverse(r)
		return r
	default:
		return make([]int, len(code))
	}
}

// only handles k > 0
func doDecrypt(code []int, k int) []int {
	r := make([]int, len(code))
	n := len(code)
	for i := 0; i < n; i++ {
		s := 0
		for j := i + 1; j < k+i+1; j++ {
			s += code[j%n]
		}
		r[i] = s
	}

	return r
}
