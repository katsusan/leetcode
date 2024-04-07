package code

type RecentCounter struct {
	q []int
}

func Constructor993() RecentCounter {
	return RecentCounter{
		q: make([]int, 0, 8),
	}
}

func (this *RecentCounter) Ping(t int) int {
	curq := this.q
	curq = append(curq, t)
	oldest := t - 3000
	i := 0
	for i < len(curq) {
		if curq[i] >= oldest {
			break
		}
		i++
	}
	this.q = curq[i:]
	//fmt.Println("this.q:", this.q)
	return len(this.q)
}

/**
 * Your RecentCounter object will be instantiated and called as such:
 * obj := Constructor();
 * param_1 := obj.Ping(t);
 */
