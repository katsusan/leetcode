package code

/*
Given the head of a linked list and a value x, partition it such that all nodes less than x come before nodes greater than or equal to x.

Example:
	Input: head = [1,4,3,2,5,2], x = 3
	Output: [1,2,2,4,3,5]
*/
func partition(head *ListNode, x int) *ListNode {
	lt := &ListNode{Val: -1} // nodes whose value < x
	gt := &ListNode{Val: -1} // nodes whose value >= x
	p1, p2 := lt, gt

	cur := head

	for cur != nil {
		if cur.Val < x {
			p1.Next = cur
			p1 = p1.Next
		} else {
			p2.Next = cur
			p2 = p2.Next
		}

		temp := cur.Next
		cur.Next = nil
		cur = temp
	}

	// merge two linked lists
	p1.Next = gt.Next
	return lt.Next
}
