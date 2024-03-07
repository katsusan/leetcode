package code

/*
Given a string s consisting only of characters 'a', 'b', and 'c'. You are asked to apply the following algorithm on the string any number of times:

	Pick a non-empty prefix from the string s where all the characters in the prefix are equal.
	Pick a non-empty suffix from the string s where all the characters in this suffix are equal.
	The prefix and the suffix should not intersect at any index.
	The characters from the prefix and suffix must be the same.
	Delete both the prefix and the suffix.
Return the minimum length of s after performing the above operation any number of times (possibly zero times).

Eg. "cabaabac" => 0
	"ca" => 2
*/

func minimumLength(s string) int {
	i, j := 0, len(s)-1
	for i < j {
		if s[i] != s[j] {
			return j - i + 1
		}

		cur := s[i]
		for i <= j && s[i] == cur {
			i++
		}

		for j > i && s[j] == cur {
			j--
		}
	}

	return j - i + 1
}

// func minimumLength(s string) int {
//     i, j := 0, len(s) - 1
//     for i < j {
//         if s[i] != s[j] {
//             return j - i + 1
//         }

//         preSteps := preLastByte(s[i+1:], s[i])
//         //fmt.Printf("i advance %d steps\n", preSteps + 1)
//         i += (preSteps + 1)

//         sufSteps := sufFirstByte(s[:j], s[j])
//         //fmt.Printf("j back %d steps\n", sufSteps + 1)
//         j -= (sufSteps + 1)
//     }

//     //fmt.Printf("finally i=%d j=%d\n", i, j)
//     if i > j {
//         return 0
//     }
//     return i - j + 1
// }

// // return continous length from starting byte
// func preLastByte(s string, b byte) int {
//     idx := 0
//     for idx < len(s) && s[idx] == b {
//         idx++
//     }

//     return idx
// }

// // return continous length from ending byte
// func sufFirstByte(s string, b byte) int {
//     idx := len(s) - 1
//     for idx >= 0 && s[idx] == b {
//         idx--
//     }

//     return len(s) - 1 - idx
// }
