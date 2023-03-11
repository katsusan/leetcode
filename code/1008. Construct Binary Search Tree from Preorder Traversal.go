package code

/*
A binary search tree is a binary tree where for every node, any descendant of Node.left has a
value strictly less than Node.val, and any descendant of Node.right has a value strictly greater than Node.val.

A preorder traversal of a binary tree displays the value of the node first, then traverses Node.left, then traverses Node.right.

Input: preorder = [1,3]
Output: [1,null,3]
*/

func bstFromPreorder(preorder []int) *TreeNode {
	if len(preorder) == 0 {
		return nil
	}
	if len(preorder) == 1 {
		return &TreeNode{Val: preorder[0]}
	}
	rNodeIdx := findRightNodeIdx(preorder)
	return &TreeNode{
		Val:   preorder[0],
		Left:  bstFromPreorder(preorder[1:rNodeIdx]),
		Right: bstFromPreorder(preorder[rNodeIdx:]),
	}
}

// returns the index of first elem which is greater than preorder[0].
// len(preorder) is assured to be >=1.
func findRightNodeIdx(preorder []int) int {
	root := preorder[0]
	for i := 0; i < len(preorder); i++ {
		if preorder[i] > root {
			return i
		}
	}

	// no right nodes
	return len(preorder)
}
