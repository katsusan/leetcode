package code

func isUgly(n int) bool {
	if n == 0 {
		return false
	}

	for {
		if n == 1 {
			return true
		}

		switch {
		case n%2 == 0:
			n = n / 2
		case n%3 == 0:
			n = n / 3
		case n%5 == 0:
			n = n / 5
		default:
			return false
		}
	}

	return false
}
