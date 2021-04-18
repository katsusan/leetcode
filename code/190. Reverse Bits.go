package code

func reverseBits(num uint32) uint32 {

	var res uint32
	for i := 0; i < 32; i++ {
		res += ((num & (1 << uint(i))) >> uint(i)) << uint(31-i)
	}
	return res
}
