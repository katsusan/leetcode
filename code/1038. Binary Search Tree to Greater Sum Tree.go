package code

//Definition for a binary tree node.
/*type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}*/

func bstToGst(root *TreeNode) *TreeNode {
	values := ValInOrder(root)

	for i := len(values) - 2; i >= 0; i-- {
		values[i] = values[i] + values[i+1]
	}

	nodes := NodeInOrder(root)

	for j := 0; j < len(nodes); j++ {
		nodes[j].Val = values[j]
	}
	return root
}

//Input: binary search tree
//Oupt: ordered values
func ValInOrder(root *TreeNode) []int {
	if root == nil {
		return []int{}
	}

	return append(ValInOrder(root.Left), append([]int{root.Val}, ValInOrder(root.Right)...)...)

}

//Input: binary search tree
//Oupt: Inorder node list
func NodeInOrder(root *TreeNode) []*TreeNode {
	if root == nil {
		return []*TreeNode{}
	}

	return append(NodeInOrder(root.Left), append([]*TreeNode{root}, NodeInOrder(root.Right)...)...)
}
