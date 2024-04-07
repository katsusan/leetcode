package code

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func deleteMiddle(head *ListNode) *ListNode {
	dummy := &ListNode{Next: head}
	oldslow, slow, fast := dummy, dummy, dummy
	for fast != nil && fast.Next != nil {
		oldslow = slow
		slow = slow.Next
		fast = fast.Next.Next
	}

	if fast == nil {
		// delete slow node
		oldslow.Next = oldslow.Next.Next
	} else if fast.Next == nil {
		// delete slow.Next node
		slow.Next = slow.Next.Next
	}
	return dummy.Next
}
