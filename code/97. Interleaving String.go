package code

import "fmt"

//time complexity: O(2^(m+n))
func isInterleave(s1 string, s2 string, s3 string) bool {
	if len(s1)+len(s2) != len(s3) {
		return false
	}

	if s1 == "" || s2 == "" {
		return (s1 + s2) == s3
	}

	t1, t2, t3 := s1[len(s1)-1], s2[len(s2)-1], s3[len(s3)-1]

	if t1 == t3 && t2 != t3 {
		return isInterleave(s1[:len(s1)-1], s2, s3[:len(s3)-1])
	} else if t2 == t3 && t1 != t3 {
		return isInterleave(s1, s2[:len(s2)-1], s3[:len(s3)-1])
	} else if t1 == t3 && t2 == t3 {
		return isInterleave(s1[:len(s1)-1], s2, s3[:len(s3)-1]) || isInterleave(s1, s2[:len(s2)-1], s3[:len(s3)-1])
	}

	return false
}

//time complexity: O(m*n)
func isInterleaveDP(s1 string, s2 string, s3 string) bool {
	if len(s1)+len(s2) != len(s3) {
		return false
	}

	arr := make([]bool, len(s2)+1)

	for i := 0; i <= len(s1); i++ {
		for j := 0; j <= len(s2); j++ {
			if i == 0 && j == 0 {
				arr[j] = true
			} else if i == 0 {
				arr[j] = arr[j-1] && (s2[j-1] == s3[j-1])
			} else if j == 0 {
				arr[j] = arr[j] && (s1[i-1] == s3[i-1])
			} else {
				arr[j] = (arr[j] && s1[i-1] == s3[i+j-1]) || (arr[j-1] && s2[j-1] == s3[i+j-1])
			}
		}
		fmt.Println(arr)
	}

	return arr[len(s2)]
}
