package main

import "math"

func titleToNumber(s string) int {
	if s == "" {
		return 0
	}

	var sum int
	for i := 0; i < len(s); i++ {
		sum += int(s[i]-64) * int(math.Pow(26, float64(len(s)-1-i)))
	}
	return sum
}
