package main

//Method 1. use hash table to verify if it has been traversed.
func detectCycleHash(head *ListNode) *ListNode {
	cm := make(map[*ListNode]bool, 0)
	for tmp := head; tmp != nil; tmp = tmp.Next {
		if _, found := cm[tmp]; found {
			return tmp
		} else {
			cm[tmp] = true
		}
	}

	return nil
}

//Method 2. from Freud's theroy of  tortoise and hare race.
func detectCycle(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return nil
	}

	tot := head
	hare := head

	for {
		if hare != nil && hare.Next != nil {
			tot = tot.Next
			hare = hare.Next.Next
			if tot == hare {
				break
			}
		} else {
			return nil
		}
	}

	//put hare to beginning
	hare = head
	for tot != hare {
		tot = tot.Next
		hare = hare.Next
	}
	return hare
}
