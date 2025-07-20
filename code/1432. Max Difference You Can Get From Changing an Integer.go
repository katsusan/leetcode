package code

import (
	"slices"
)

func maxDiff(num int) int {
	// convert 1234 into []int{4,3,2,1}
	digits := make([]int, 0, 9)
	for n := num; n > 0; {
		digits = append(digits, n%10)
		n = n / 10
	}

	// replace highest non-9 digit with 9
	h := slices.Clone(digits)
	digitPicked := -1
	for i := len(h) - 1; i >= 0; i-- {
		if h[i] != 9 && digitPicked == -1 {
			digitPicked = h[i]
		}

		if h[i] == digitPicked {
			h[i] = 9
		}
	}

	// replace highest non-1 digit with 1
	//   - highest digit != 1 -> dPicked = l(len(h[l]-1))
	//   - highest digit == 1 -> rest highest digit && it != 1
	l := slices.Clone(digits)
	leadingDigit := l[len(l)-1]
	digitPicked = -1
	if leadingDigit != 1 {
		digitPicked = leadingDigit
	}
	for i := len(l) - 1; i >= 0; i-- {
		if l[i] != 1 && l[i] != 0 && digitPicked == -1 {
			digitPicked = l[i]
		}

		if l[i] == digitPicked {
			if leadingDigit != 1 {
				l[i] = 1
			} else {
				l[i] = 0
			}
		}
	}

	diff := 0
	accu := 1
	for i := 0; i < len(h); i++ {
		diff += (h[i] - l[i]) * accu
		accu *= 10
	}

	return diff
}
