package code

import (
	"container/heap"
	"sort"
	"strings"
)

func suggestedProducts(products []string, searchWord string) [][]string {
	t := &trie{}
	for _, prod := range products {
		t.insert(prod)
	}

	sug := [][]string{}
	for i := range searchWord {
		allWords := t.search(searchWord[i])
		h := &minheapStrN{}
		for _, w := range allWords {
			heap.Push(h, w)
		}
		sort.Sort(h)
		sug = append(sug, h.c)
	}
	return sug
}

const heapN = 3 // at most three words suggested
const sumLen = 20005

type minheapStrN struct {
	c []string
	n int
}

func (h minheapStrN) Len() int {
	return len(h.c)
}

func (h minheapStrN) Less(i, j int) bool {
	return strings.Compare(h.c[i], h.c[j]) > 0
}

func (h minheapStrN) Swap(i, j int) {
	h.c[i], h.c[j] = h.c[j], h.c[i]
}

func (h *minheapStrN) Pop() any {
	old := h.c
	n := len(h.c)
	x := old[n-1]
	h.c = h.c[:n-1]
	return x
}

func (h *minheapStrN) Push(x any) {

}

type trie struct {
	nex   [10][26]int
	exist [20]bool
	cnt   int
}

func (t *trie) insert(s string) {
	p := 0
	for i := 0; i < len(s); i++ {
		c := int(s[i] - 'a')
		if t.nex[p][c] == 0 {
			t.cnt++
			//fmt.Printf("insert %s: set t.nex[%d][%d]=%d\n", s, p, c, t.cnt)
			t.nex[p][c] = t.cnt
		}
		p = t.nex[p][c]
	}
	//fmt.Printf("insert: set t.exist[%d]=true\n", p)
	t.exist[p] = true
}

func (t *trie) find(s string) bool {
	p := 0
	for i := 0; i < len(s); i++ {
		c := int(s[i] - 'a')
		//fmt.Printf("find %s: traversing t.nex[%d][%d]=%d\n", s, p, c, t.nex[p][c])
		if t.nex[p][c] == 0 {
			return false
		}
		p = t.nex[p][c]
	}
	return t.exist[p]
}

// return all words start with w
func (t *trie) search(b byte) []string {
	return []string{}
}
