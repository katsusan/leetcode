package main

type MedianFinder struct {
	minh   minheap
	minlen int

	maxh   maxheap
	maxlen int
}

type minheap []int

type maxheap []int

func (h *minheap) swim(i int) {
	cur := i
	for cur > 0 {
		par := (cur - 1) / 2
		if (*h)[cur] < (*h)[par] {
			(*h)[cur], (*h)[par] = (*h)[par], (*h)[cur]
		}
		cur = par
	}
}

func (h *minheap) sink(i int) {
	cur := i
	for 2*cur+1 < len(*h) {
		par := 2*cur + 1
		if par+1 < len(*h) && (*h)[par] > (*h)[par+1] {
			par++
		}
		if (*h)[cur] < (*h)[par] {
			break
		}

		(*h)[cur], (*h)[par] = (*h)[par], (*h)[cur]
		cur = par
	}
}

func (h *maxheap) swim(i int) {
	cur := i
	for cur > 0 {
		par := (cur - 1) / 2
		if (*h)[cur] > (*h)[par] {
			(*h)[cur], (*h)[par] = (*h)[par], (*h)[cur]
		}
		cur = par
	}
}

func (h *maxheap) sink(i int) {
	cur := i
	for 2*cur+1 < len(*h) {
		par := 2*cur + 1
		if par+1 < len(*h) && (*h)[par] < (*h)[par+1] {
			par++
		}
		if (*h)[cur] > (*h)[par] {
			break
		}
		(*h)[cur], (*h)[par] = (*h)[par], (*h)[cur]
		cur = par
	}
}

/** initialize your data structure here. */
func Constructor() MedianFinder {
	return MedianFinder{
		minh:   make([]int, 0),
		minlen: 0,
		maxh:   make([]int, 0),
		maxlen: 0,
	}
}

func (this *MedianFinder) AddNum(num int) {
	if this.maxlen == this.minlen {
		if len(this.minh) > 0 && num > this.minh[0] {
			//swap minh[0] and num, then sink minh[0]
			this.minh[0], num = num, this.minh[0]
			this.minh.sink(0)
		}
		//add to left(maxheap)
		this.maxh = append(this.maxh, num)
		this.maxh.swim(len(this.maxh) - 1)
		this.maxlen++
	} else {
		if num < this.maxh[0] {
			//swap maxh[0] and num, then sink maxh[0]
			this.maxh[0], num = num, this.maxh[0]
			this.maxh.sink(0)
		}
		//add to right(minheap)
		this.minh = append(this.minh, num)
		this.minh.swim(len(this.minh) - 1)
		this.minlen++
	}
}

func (this *MedianFinder) FindMedian() float64 {
	if this.maxlen == 0 {
		//both empty
		return 0
	}

	if this.minlen == 0 {
		return float64(this.maxh[0])
	}

	if this.minlen == this.maxlen {
		return (float64(this.minh[0]) + float64(this.maxh[0])) / 2
	} else {
		return float64(this.maxh[0])
	}
}

/**
 * Your MedianFinder object will be instantiated and called as such:
 * obj := Constructor();
 * obj.AddNum(num);
 * param_2 := obj.FindMedian();
 */
