package code

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
	res := new(ListNode)
	pr := res

	p1 := l1
	p2 := l2
	for p1 != nil || p2 != nil { //O(len(l1)+len(l2))
		pr.Next = new(ListNode)
		pr = pr.Next

		if p2 == nil && p1 != nil {
			pr.Val = p1.Val
			p1 = p1.Next
			continue
		}

		if p1 == nil && p2 != nil {
			pr.Val = p2.Val
			p2 = p2.Next
			continue
		}

		if p1.Val < p2.Val {
			pr.Val = p1.Val
			p1 = p1.Next
		} else {
			pr.Val = p2.Val
			p2 = p2.Next
		}

	}
	pr = nil
	return res.Next

}
