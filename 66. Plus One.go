package main

//URL: https://leetcode.com/problems/plus-one/
//use recursion
func plusOne(digits []int) []int {
	if len(digits) == 0 {
		return []int{}
	}

	ca, tail := plusOneWithCarry(digits[len(digits)-1])
	if len(digits) == 1 {
		if ca == 0 {
			return []int{tail}
		} else {
			return []int{ca, tail}
		}
	}

	if ca == 0 {
		return append(digits[:len(digits)-1], tail)
	}

	return append(plusOne(digits[:len(digits)-1]), tail)
}

func plusOneWithCarry(digit int) (ca int, res int) {
	if digit == 9 {
		return 1, 0
	}
	return 0, digit + 1
}
