package ds

import (
	"strconv"
	"testing"
)

func TestBloomFilterExist(t *testing.T) {
	var total int = 1 << 20
	bf := NewBloomFilter(uint(total), 0.0003)

	for i := 0; i < total; i++ {
		bf.Add([]byte(strconv.Itoa(i)))
	}

	var fpcount int

	for i := 0; i < total; i++ {
		if bf.Exist([]byte(strconv.Itoa(i))) {
			fpcount++
		}
	}

	if fpcount != 1<<20 {
		t.Errorf("without Remove, all should exist, expect fpcount=1<<20, got %d\n", fpcount)
	}

}

func TestBloomFilterRemove(t *testing.T) {
	var total int = 1 << 20
	bf := NewBloomFilter(uint(total), 0.0003)

	bf.Add([]byte("aa"))

	if !bf.Exist([]byte("aa")) {
		t.Error("only one element added, expect it exists\n")
	}

	bf.Remove([]byte("aa"))
	if bf.Exist([]byte("aa")) {
		t.Error("the only element removed, expect it doesn't exist\n")
	}
}
