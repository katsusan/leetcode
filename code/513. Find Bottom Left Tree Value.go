package code

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func findBottomLeftValue(root *TreeNode) int {
	leftmost := root
	q := []*TreeNode{root}
	nextq := []*TreeNode{}
	for {
		// drain current layer
		for len(q) != 0 {
			cur := q[0]
			q = q[1:]

			if cur.Left != nil {
				nextq = append(nextq, cur.Left)
			}

			if cur.Right != nil {
				nextq = append(nextq, cur.Right)
			}
		}

		// end layer
		if len(nextq) == 0 {
			break
		}

		q = nextq
		nextq = []*TreeNode{}
		leftmost = q[0]
	}
	return leftmost.Val
}
