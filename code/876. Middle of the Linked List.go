package code

func middleNode(head *ListNode) *ListNode {
	pfast, pslow := head, head

	for pfast.Next != nil && pfast.Next.Next != nil {
		pslow = pslow.Next
		pfast = pfast.Next.Next
	}

	if pfast.Next == nil {
		return pslow
	}

	return pslow.Next
}
