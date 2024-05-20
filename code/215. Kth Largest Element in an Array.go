package code

import "container/heap"

func findKthLargest(nums []int, k int) int {
	minh := &minheapN{
		items: make([]int, 0, k),
		size:  k,
	}

	for _, v := range nums {
		heap.Push(minh, v)
	}

	return minh.items[0]
}

type minheapN struct {
	items []int
	size  int
}

func (h minheapN) Len() int {
	return len(h.items)
}

func (h minheapN) Less(i, j int) bool {
	return h.items[i] < h.items[j]
}

func (h minheapN) Swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

func (h *minheapN) Push(x any) {
	v := x.(int)
	if len(h.items) < h.size {
		h.items = append(h.items, v)
		return
	}

	if v > h.items[0] {
		// repace the minimum value with new one
		h.items[0] = v
		heap.Fix(h, 0)
	}

}

func (h *minheapN) Pop() any {
	old := h.items
	n := len(old)
	x := old[n-1]
	h.items = old[:n-1]
	return x
}
