package code

// O(1) space and O(n) time
func fib(n int) int {
	var a, b = 0, 1
	for n > 0 {
		a, b = b, a+b
		n--
	}
	return a
}
