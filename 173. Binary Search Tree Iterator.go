package main

type BSTIterator struct {
	vals []int
	cur  int
}

func Traverse(root *TreeNode) []int {
	if root == nil {
		return []int{}
	}

	return append(Traverse(root.Left), append([]int{root.Val}, Traverse(root.Right)...)...)
}

func BSTIterConstructor(root *TreeNode) BSTIterator {
	return BSTIterator{
		vals: Traverse(root),
		cur:  -1,
	}
}

/** @return the next smallest number */
func (this *BSTIterator) Next() int {
	this.cur++
	if this.cur < len(this.vals) {
		return this.vals[this.cur]
	}
	return 0
}

/** @return whether we have a next smallest number */
func (this *BSTIterator) HasNext() bool {
	return this.cur < len(this.vals)-1
}
