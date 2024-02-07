package code

import "strings"

func gcdOfStrings(str1 string, str2 string) string {
	var gcdStr string
	for i := 0; i < len(str1); i++ {
		// test if str[0:i+1] can divide str1
		t := str1[0 : i+1]
		if canDivide(str1, t) && canDivide(str2, t) {
			gcd := getGCD(len(str1)/len(t), len(str2)/len(t))
			gcdStr = strings.Repeat(t, gcd)
			break
		}
	}
	return gcdStr
}

// Euclidean Algorithm
func getGCD(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}

// canDivide reports if son is able to divide parent.
func canDivide(parent, son string) bool {
	if len(parent)%len(son) != 0 {
		return false
	}

	times := len(parent) / len(son)
	unitLen := len(son)
	for i := 0; i < times; i++ {
		if parent[i*unitLen:(i+1)*unitLen] != son {
			return false
		}
	}

	return true
}
