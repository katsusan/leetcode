package main

func hammingWeight(num uint32) int {
	if num <= 1 {
		return int(num)
	}
	return int(num%2) + hammingWeight(num/2)
}
