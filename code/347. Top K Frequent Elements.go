package code

import "fmt"

//use k-size heap to store the current top k frequence of nums.
func topKFrequent(nums []int, k int) []int {
	dic := make(map[int]int, 0)

	for _, v := range nums {
		dic[v]++
	}

	h := &minHeap{}
	var counter int
	fmt.Println(dic)
	for key, val := range dic {
		if counter < k {
			//fmt.Println("will push", key, val)
			h.Push(entry{key, val})
			counter++
			continue
		}
		//fmt.Println("will cmpsink", key, val)
		h.CompareAndSink(entry{key, val})
	}

	var res []int

	for _, v := range *h {
		res = append(res, v.num)
	}
	return res

}

type minHeap []entry

type entry struct {
	num  int
	freq int
}

func (h *minHeap) Push(e entry) {
	*h = append(*h, e)

	//swim the last one
	i := len(*h) - 1
	for i > 0 && (*h)[i].freq < (*h)[(i-1)/2].freq {
		(*h)[i], (*h)[(i-1)/2] = (*h)[(i-1)/2], (*h)[i]
		i = (i - 1) / 2
	}
}

func (h *minHeap) CompareAndSink(e entry) {
	if len(*h) == 0 {
		return
	}

	if e.freq > (*h)[0].freq {
		var i int
		(*h)[0] = e
		size := len(*h)
		for 2*i+1 < size {
			p := 2*i + 1
			//look for the smaller one of parents
			if p < size-1 && (*h)[p].freq > (*h)[p+1].freq {
				p++
			}
			if (*h)[i].freq <= (*h)[p].freq {
				break
			}
			(*h)[i], (*h)[p] = (*h)[p], (*h)[i]
			i = p
		}
	}
}
