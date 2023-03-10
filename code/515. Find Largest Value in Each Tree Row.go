package code

/*
Given the root of a binary tree, return an array of the largest value in each row of the tree (0-indexed).

Input: root = [1,3,2,5,3,null,9]
Output: [1,3,9]
*/

// use Breadth-First-Search(BFS)
func largestValues(root *TreeNode) []int {
	if root == nil {
		return nil
	}

	var largest = make([]int, 0, 2)
	var q = []*TreeNode{root}

	for len(q) > 0 {
		qlen := len(q)
		maxValue := q[0].Val
		for i := 0; i < qlen; i++ {
			cur := q[0]
			q = q[1:]
			if cur.Left != nil {
				q = append(q, cur.Left)
			}

			if cur.Right != nil {
				q = append(q, cur.Right)
			}

			if cur.Val > maxValue {
				maxValue = cur.Val
			}
		}
		largest = append(largest, maxValue)
	}

	return largest
}
