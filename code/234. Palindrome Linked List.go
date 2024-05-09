package code

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func isPalindrome(head *ListNode) bool {
	fast, slow := head, head
	for fast.Next != nil && fast.Next.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}

	var mid = slow
	var half0, half1 *ListNode
	if fast.Next == nil {
		half1 = mid.Next
		mid.Next = nil
		half0 = reverseLinkedList(head).Next
	} else {
		half1 = mid.Next
		mid.Next = nil
		half0 = reverseLinkedList(head)
	}

	return LinkedListEq(half0, half1)
}

func reverseLinkedList(head *ListNode) *ListNode {
	cur := head
	prev := (*ListNode)(nil)
	next := cur.Next
	for cur != nil {
		if next == nil {
			cur.Next = prev
			return cur
		} else if next.Next == nil {
			next.Next = cur
			cur.Next = prev
			return next
		}

		temp := next.Next
		cur.Next = prev
		next.Next = cur
		prev = next
		cur = temp
		next = cur.Next
	}

	return next
}

func LinkedListEq(m, n *ListNode) bool {
	for m != nil && n != nil {
		if m.Val != n.Val {
			return false
		}
		m = m.Next
		n = n.Next
	}

	return m == nil && n == nil
}
