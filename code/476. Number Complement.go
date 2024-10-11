package code

func findComplement(num int) int {
	bits := 0
	n := num
	for n > 0 {
		bits++
		n = n >> 1
	}

	return (1<<bits - 1) - num
}
