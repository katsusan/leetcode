package code

import "slices"

// O(n^2) time complexity
func tupleSameProductDepre(nums []int) int {
	slices.Sort(nums)

	combs := 0
	m := make(map[int]int, len(nums))
	for i, v := range nums {
		m[v] = i
	}

	traverse := func(i, j int) {
		//fmt.Printf("will traverse %d to %d\n", i, j)
		prod := nums[i] * nums[j]
		for k := i + 1; k < j; k++ {
			if canDiv(prod, nums[k]) {
				if nums[k]*nums[k] > prod {
					break
				}

				d := prod / nums[k]
				if v, ok := m[d]; ok && v != k {
					//fmt.Printf("hit! %d * %d = %d * %d\n", nums[i], nums[j], nums[k], nums[v])
					combs++
				}
			}
		}
	}

	n := len(nums)
	for step := 0; step < n/2; step++ {
		traverse(step, n-1-step)
		nn := n - 1 - step
		for i := 1; i < nn; i++ {
			traverse(i+step, nn)
			traverse(step, nn-i)
		}
	}

	return 8 * combs
}

// return if x can divide y
func canDiv(x, y int) bool {
	return x%y == 0
}
