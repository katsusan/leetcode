package code

// problem 82 may be solved with two pointers, which is more efficient.
// TODO: rewrite with two-pointers.

func deleteDuplicates_82(head *ListNode) *ListNode {
	dummy := new(ListNode)
	dummy.Next = head
	cur := dummy

	for cur != nil {
		var (
			next  *ListNode
			isdup bool
		)
		for {
			if next, isdup = checkdup(cur); isdup {
				cur.Next = next
				continue
			}
			break
		}
		cur = cur.Next
	}

	return dummy.Next
}

// checkdup will iterate head and return if subsequent node have same value with head.
// return the first node whose value is different from head.Next(if head has Next).
// NOTE: this function doesn't modify the linked list.
// exampe:
//  1 -> nil    =>  (nil, false)
//  1 -> 2 -> nil   =>  (2, false)
//  1 -> 2 -> 2 -> ... -> k -> ...   =>  (k, true)
func checkdup(head *ListNode) (*ListNode, bool) {
	if head == nil {
		return nil, false
	}
	var dups int
	start := head.Next
	end := start

	for end != nil && end.Val == start.Val {
		end = end.Next
		dups++
	}

	return end, dups > 1
}
