package code

/*
Given the root of a binary tree, each node in the tree has a distinct value.

After deleting all nodes with a value in to_delete, we are left with a forest
 (a disjoint union of trees).

Return the roots of the trees in the remaining forest. You may return the
result in any order.
*/

func delNodes(root *TreeNode, to_delete []int) []*TreeNode {
	forest := make([]*TreeNode, 0, 4)
	deleteSet := make(map[int]struct{}, len(to_delete))
	for _, d := range to_delete {
		deleteSet[d] = struct{}{}
	}

	if _, exist := deleteSet[root.Val]; !exist {
		forest = append(forest, root)
	}

	traverse1110(root, deleteSet, &forest)
	return forest
}

func traverse1110(node *TreeNode, toDelete map[int]struct{}, fst *[]*TreeNode) *TreeNode {
	if node == nil {
		return node
	}

	node.Left = traverse1110(node.Left, toDelete, fst)
	node.Right = traverse1110(node.Right, toDelete, fst)
	if _, exist := toDelete[node.Val]; exist {
		if node.Left != nil {
			*fst = append(*fst, node.Left)
		}
		if node.Right != nil {
			*fst = append(*fst, node.Right)
		}
		return nil
	}

	return node
}
