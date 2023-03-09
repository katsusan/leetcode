package code

/*
Given the root of a binary tree, return the length of the diameter of the tree.

The diameter of a binary tree is the length of the longest path between any two nodes in a tree. This path may or may not pass through the root.

The length of a path between two nodes is represented by the number of edges between them.

Input: root = [1,2,3,4,5]
Output: 3
Explanation: 3 is the length of the path [4,2,1,3] or [5,2,1,3].
*/

func diameterOfBinaryTree(root *TreeNode) int {
	var maxDiameter int
	maxDepth(&maxDiameter, root)
	return maxDiameter
}

func maxDepth(max *int, node *TreeNode) int {
	if node == nil {
		return 0
	}

	leftMax := maxDepth(max, node.Left)
	rightMax := maxDepth(max, node.Right)

	curDiameter := leftMax + rightMax
	if curDiameter > *max {
		*max = curDiameter
	}

	return 1 + maxInt(leftMax, rightMax)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
