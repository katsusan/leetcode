package main

func swapPairs(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	newhead := head.Next
	newtail := swapPairs(head.Next.Next)
	head.Next, head.Next.Next = newtail, head

	return newhead
}
