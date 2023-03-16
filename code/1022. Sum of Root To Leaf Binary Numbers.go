package code

/*
You are given the root of a binary tree where each node has a value 0 or 1.
Each root-to-leaf path represents a binary number starting with the most
significant bit.

For example, if the path is 0 -> 1 -> 1 -> 0 -> 1, then this could represent
01101 in binary, which is 13.
*/

func sumRootToLeaf(root *TreeNode) int {
	var sum = new(int)
	traverse(root, 0, sum)
	return *sum
}

func traverse(node *TreeNode, pathVal int, sum *int) {
	if node == nil {
		return
	}

	if node.Left == nil && node.Right == nil {
		*sum = *sum + pathVal<<1 + node.Val
		return
	}

	pathVal = pathVal<<1 + node.Val
	traverse(node.Left, pathVal, sum)
	traverse(node.Right, pathVal, sum)
	pathVal = pathVal >> 1
}
