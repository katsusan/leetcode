package code

import "slices"

func maximumSwap670(num int) int {
	bits := make([]int, 0, 9)
	r := num
	for i := 9; i >= 0; i-- {
		b := r / pow10(i)
		if b != 0 || (b == 0 && len(bits) != 0) {
			bits = append(bits, b)
			r -= b * pow10(i)
		}
	}

	swapIdx := -1
	destBit := -1
	for i := 0; i < len(bits); i++ {
		maxBit := slices.Max(bits[i:])
		if bits[i] != maxBit {
			swapIdx = i // swapIdx is the bit we need to do swapping
			destBit = maxBit
			break
		}
	}

	// already maximum value
	if swapIdx == -1 {
		return num
	}

	// swap with the lowest maximum bit
	destIdx := -1
	for j := len(bits) - 1; j > swapIdx; j-- {
		if bits[j] == destBit {
			destIdx = j
			break
		}
	}

	bits[swapIdx], bits[destIdx] = bits[destIdx], bits[swapIdx]

	t := len(bits) - 1
	offset := (bits[destIdx] - bits[swapIdx]) * (pow10(t-destIdx) - pow10(t-swapIdx))
	return num + offset
}

func pow10(n int) int {
	e := 1
	for i := 0; i < n; i++ {
		e = e * 10
	}
	return e
}
