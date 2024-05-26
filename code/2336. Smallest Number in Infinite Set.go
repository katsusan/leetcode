package code

import "container/heap"

type SmallestInfiniteSet struct {
	h     minxheap
	exist map[int]bool
}

func Constructor2336() SmallestInfiniteSet {
	s := SmallestInfiniteSet{}
	s.exist = make(map[int]bool, 1000)
	for i := 1; i <= 1000; i++ {
		heap.Push(&s.h, i)
		s.exist[i] = true
	}

	return s
}

func (this *SmallestInfiniteSet) PopSmallest() int {
	x := heap.Pop(&this.h).(int)
	delete(this.exist, x)
	return x
}

func (this *SmallestInfiniteSet) AddBack(num int) {
	if !this.exist[num] {
		heap.Push(&this.h, num)
		this.exist[num] = true
	}
}

/**
 * Your SmallestInfiniteSet object will be instantiated and called as such:
 * obj := Constructor();
 * param_1 := obj.PopSmallest();
 * obj.AddBack(num);
 */

type minxheap []int

func (h minxheap) Len() int           { return len(h) }
func (h minxheap) Less(i, j int) bool { return h[i] < h[j] }
func (h minxheap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *minxheap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *minxheap) Pop() any {
	n := len(*h)
	old := *h
	*h = old[:n-1]
	x := old[n-1]
	return x
}
