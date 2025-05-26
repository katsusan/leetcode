package code

import "slices"

func findEvenNumbers(digits []int) []int {
	digitMap := make([]int, 10)
	for _, v := range digits {
		digitMap[v]++
	}

	r := make([]int, 0, 8)
	for i := 100; i < 1000; i++ {
		if i%2 != 0 {
			continue
		}
		bm := slices.Clone(digitMap)

		c := i % 10
		bm[c]--
		if bm[c] < 0 {
			continue
		}

		b := (i / 10) % 10
		bm[b]--
		if bm[b] < 0 {
			continue
		}

		a := i / 100
		bm[a]--
		if bm[a] < 0 {
			continue
		}

		r = append(r, i)
	}

	return r
}
