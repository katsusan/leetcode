package code

func rangeBitwiseAnd(left int, right int) int {
	bits := 0
	for left != right {
		left >>= 1
		right >>= 1
		bits++
	}
	return left << bits
}
