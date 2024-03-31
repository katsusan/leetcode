package code

func tribonacci(n int) int {
	a, b, c := 0, 1, 1

	for i := 0; i < n; i++ {
		tmpb, tmpc := b, c
		c = a + b + c
		b = tmpc
		a = tmpb
	}

	return a
}
