package code

/*
A binary tree is named Even-Odd if it meets the following conditions:
	- The root of the binary tree is at level index 0, its children are at level index 1, their children are at level index 2, etc.
	- For every even-indexed level, all nodes at the level have odd integer values in strictly increasing order (from left to right).
	- For every odd-indexed level, all nodes at the level have even integer values in strictly decreasing order (from left to right).
Given the root of a binary tree, return true if the binary tree is Even-Odd, otherwise return false.
*/

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func isEvenOddTree(root *TreeNode) bool {
	q := []*TreeNode{root}
	nextq := []*TreeNode{}
	depth := 0
	for {
		//fmt.Printf("verifying q=%#v\n", q)
		if !IsStrictlyOrdered(q, depth) {
			return false
		}

		for len(q) != 0 {
			cur := q[0]
			q = q[1:]
			for _, child := range []*TreeNode{cur.Left, cur.Right} {
				if child != nil {
					nextq = append(nextq, child)
				}
			}
		}

		if len(nextq) == 0 {
			break
		}
		q = nextq
		nextq = []*TreeNode{}
		depth++
	}
	return true
}

// reports where x is strictly order.
func IsStrictlyOrdered(x []*TreeNode, depth int) bool {
	if len(x) == 1 {
		if depth%2 == 0 {
			return x[0].Val%2 != 0
		}
		return x[0].Val%2 == 0
	}

	i, j := 0, 1
	// even-indexed level, ascending and odd value
	if depth%2 == 0 {
		for j < len(x) {
			if x[j].Val%2 != 1 || x[i].Val%2 != 1 || x[j].Val <= x[i].Val {
				//fmt.Printf("got x[j].Val=%d, x[i].Val=%d, order=%d\n", x[j].Val, x[i].Val, order)
				return false
			}
			i++
			j++
		}
		return true
	}

	i, j = 0, 1
	// odd-indexed level, descending and even value
	for j < len(x) {
		if x[j].Val%2 != 0 || x[i].Val%2 != 0 || x[j].Val >= x[i].Val {
			//fmt.Printf("got x[j].Val=%d, x[i].Val=%d, order=%d\n", x[j].Val, x[i].Val, order)
			return false
		}
		i++
		j++
	}
	return true
}
