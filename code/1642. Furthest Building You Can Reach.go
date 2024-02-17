package code

import "container/heap"

/*
ou are given an integer array heights representing the heights of buildings, some bricks, and some ladders.

You start your journey from building 0 and move to the next building by possibly using bricks or ladders.

While moving from building i to building i+1 (0-indexed),

If the current building's height is greater than or equal to the next building's height, you do not need a ladder or bricks.
If the current building's height is less than the next building's height, you can either use one ladder or (h[i+1] - h[i]) bricks.
Return the furthest building index (0-indexed) you can reach if you use the given ladders and bricks optimally.
*/

func furthestBuilding(heights []int, bricks int, ladders int) int {
	h := &minHeap2{} // hold steps for ladders, length <= ladders
	heap.Init(h)
	steps := make([]int, len(heights))
	remains := bricks // hold steps for bricks
	for i := 0; i < len(heights)-1; i++ {
		steps[i] = heights[i+1] - heights[i]
		if steps[i] <= 0 {
			// no need for bricks or ladders
			steps[i] = 0
			continue
		}

		heap.Push(h, steps[i])
		if len(*h) > ladders {
			x := heap.Pop(h).(int)
			remains = remains - x
		}

		if remains < 0 {
			return i
		}
	}
	return len(heights) - 1

}

// avoid confict with leetcode-347
type minHeap2 []int

func (h minHeap2) Len() int {
	return len(h)
}

func (h minHeap2) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h minHeap2) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *minHeap2) Push(elem any) {
	*h = append(*h, elem.(int))
}

func (h *minHeap2) Pop() any {
	e := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return e
}
