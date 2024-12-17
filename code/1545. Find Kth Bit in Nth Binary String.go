package code

func findKthBit1545(n int, k int) byte {
	if n == 1 {
		return '0'
	}

	// S(n)=S(n-1) + "1" + f(S(n-1)) => len(S(n)) = 2^n-1
	mid := 1 << (n - 1)
	switch {
	case k < mid:
		return findKthBit1545(n-1, k)
	case k > mid:
		nk := 1<<n - k
		if findKthBit1545(n-1, nk) == '0' {
			return '1'
		}
		return '0'
	default:
		return '1'
	}
}
