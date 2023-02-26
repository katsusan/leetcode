package code

func deleteDuplicates(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}
	var pslow, pfast = head, head

	for pfast != nil {
		if pfast.Val > pslow.Val {
			pslow.Next = pfast
			pslow = pfast
			continue
		}

		pfast = pfast.Next
	}

	pslow.Next = nil

	return head
}
