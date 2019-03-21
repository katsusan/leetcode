package main

//use Newton iteration
func mySqrt(x int) int {
	k := x
	for k*k > x {
		k = (k + x/k) / 2
	}
	return k
}
