package code

/**
 * Forward declaration of guess API.
 * @param  num   your guess
 * @return          -1 if num is higher than the picked number
 *                  1 if num is lower than the picked number
 *               otherwise return 0
 * func guess(num int) int;
 */

func guessNumber(n int) int {
	return guessRange(1, n)
}

func guessRange(a, b int) int {
	mid := (a + b) / 2
	if guess(mid) == 1 {
		return guessRange(mid+1, b)
	} else if guess(mid) == -1 {
		return guessRange(a, mid-1)
	}

	return mid
}

func guess(num int) int
