package code

/*
You have a long flowerbed in which some of the plots are planted, and some are not. However, flowers cannot be planted in adjacent plots.

Given an integer array flowerbed containing 0's and 1's, where 0 means empty and 1 means not empty, and an integer n, return true if n new
flowers can be planted in the flowerbed without violating the no-adjacent-flowers rule and false otherwise.

Input: flowerbed = [1,0,0,0,1], n = 1
Output: true

Input: flowerbed = [1,0,0,0,1], n = 2
Output: false

*/

func canPlaceFlowers(flowerbed []int, n int) bool {
	i := 0
	for i < len(flowerbed) {
		if flowerbed[i] == 1 {
			i = i + 2
			continue
		} else {
			leftOK := (i > 1 && flowerbed[i-1] == 0) || (i == 0)
			rightOK := (i < len(flowerbed)-1 && flowerbed[i+1] == 0) || (i == len(flowerbed)-1)
			//fmt.Printf("for index:%d leftOK=%v rightOK=%v\n", i, leftOK, rightOK)
			if leftOK && rightOK {
				n--
				if n <= 0 {
					return true
				}
				i = i + 2
				continue
			}
			i++

		}
	}
	return n == 0
}
