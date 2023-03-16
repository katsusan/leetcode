package code

import "math"

/*
Given the root of a binary tree, find the maximum value v for which there
exist different nodes a and b where v = |a.val - b.val| and a is an ancestor of b.
*/

func maxAncestorDiff(root *TreeNode) int {
	maxDiff := 0
	pathMin := math.MaxInt
	pathMax := math.MinInt

	var traverse func(node *TreeNode)
	traverse = func(node *TreeNode) {
		if node == nil {
			pathMaxDiff := pathMax - pathMin
			if pathMaxDiff > maxDiff {
				maxDiff = pathMaxDiff
			}
			return
		}

		if node.Val > pathMax {
			pathMax = node.Val
		}

		if node.Val < pathMin {
			pathMin = node.Val
		}

		curMax, curMin := pathMax, pathMin
		traverse(node.Left)
		pathMax, pathMin = curMax, curMin

		traverse(node.Right)
		pathMax, pathMin = curMax, curMin
	}

	traverse(root)
	return maxDiff
}
