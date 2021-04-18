package code

import (
	"container/heap"
	"fmt"
	"sort"
	"strings"
)

func topKFrequentWords(words []string, k int) []string {
	dic := make(map[string]int, len(words))

	for _, v := range words {
		dic[v]++
	}
	fmt.Println("dic", dic)

	var counter int
	wh := &wordheap{}
	heap.Init(wh)

	for w, f := range dic {
		if counter < k {
			fmt.Println("will push", w, f)
			heap.Push(wh, wordentry{w, f})
			counter++
			continue
		}
		fmt.Println("wh:", wh)
		wh.CompareAndSink(wordentry{w, f})
	}
	fmt.Println("after loop wh:", wh)

	sort.Sort(wh)
	var res []string

	for wh.Len() > 0 {
		res = append(res, wh.Pop().(wordentry).word)
	}
	return res
}

type wordentry struct {
	word string
	freq int
}

type wordheap []wordentry

//w1 < w2 => true
func compare(w1, w2 wordentry) bool {
	if w1.freq < w2.freq {
		return true
	} else if w1.freq > w2.freq {
		return false
	} else {
		comp := strings.Compare(w1.word, w2.word)
		if comp > 0 {
			return true
		}
		return false
	}
}

func (h wordheap) Len() int           { return len(h) }
func (h wordheap) Less(i, j int) bool { return compare(h[i], h[j]) }
func (h wordheap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *wordheap) Push(w interface{}) {
	*h = append(*h, w.(wordentry))
}
func (h *wordheap) Pop() interface{} {
	size := len(*h)
	if size == 0 {
		return nil
	}

	r := (*h)[size-1]
	*h = (*h)[:size-1]
	return r
}
func (h *wordheap) CompareAndSink(w interface{}) {
	if len(*h) == 0 {
		h.Push(w)
	}

	v, ok := w.(wordentry)
	if !ok {
		return
	}

	if compare((*h)[0], v) {
		//need to swap h[0] and v, and then sink the top
		(*h)[0] = v

		var i int
		size := len(*h)
		for 2*i+1 < size {
			p := 2*i + 1
			if p < size-1 && !compare((*h)[p], (*h)[p+1]) {
				//p will point to the smaller one of parents
				p++
			}
			if compare((*h)[i], (*h)[p]) {
				break
			}
			h.Swap(i, p)
			i = p
		}
	}
	return
}
