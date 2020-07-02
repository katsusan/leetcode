package ds

import (
	"hash"
	"hash/fnv"
	"math"
)

/*
布隆过滤器使用m位的bitmap以及k个哈希函数来判断一个元素肯定不在其中or很有可能在其中，注意其中的关键字"肯定", "很有可能"。
这是由布隆过滤器的一个特性保证，即常规的布隆过滤器里你是无法移除一个元素的，也就是不存在把某位从1置位为0的操作。
以6bits的布隆过滤器为例，
 hash(A) = 0 1 1 1 0 0
 hash(B) = 1 0 0 0 0 0
加入元素A,B后,bits数组变为 [1 1 1 1 0 0]
如果来一个元素C, hash(C) = 0 0 1 1 1 1，则 hash(C) | array != array => C not exists in array.
当然如果另一个元素D满足hash(D) = 0 0 1 1 0 0(也就是A的子集)，则虽然它没被加到布隆过滤器中，但却会被误判为存在，
因此哈希函数的均匀性好坏也会影响误判几率(probability of false positive)。
概念:
	https://en.wikipedia.org/wiki/Bloom_filter
fpp,m,n之间的实验数据：
	http://pages.cs.wisc.edu/~cao/papers/summary-cache/node8.html
*/

// BloomFilter describes Counting Bloom Filters with following properties:
// 	k : hash function amount, in fact, we use hash(k, function) as k different hash functions
//	m : bit length, for counting bloom filters, it is the length of array
//	h : hash function used actually
//	bfList: stores the counting information
type BloomFilter struct {
	k      int //k hash functions
	m      int //m bits
	h      hash.Hash64
	bfList []uint16
}

// NewBloomFilter will initialize a BloomFilter by given element number and false positive probability
func NewBloomFilter(totalNum uint, fpp float64) *BloomFilter {
	b := &BloomFilter{h: fnv.New64()}
	b.CalculatingMK(totalNum, fpp)
	b.bfList = make([]uint16, b.m)
	return b
}

// CalculatingMK will calculating optimal m andk by given n and fpp(false positive probability)
// m = -(n * ln(fpp)) / (ln2)^2
// k = (ln2) * (m / n)
// refer: https://en.wikipedia.org/wiki/Bloom_filter#Optimal_number_of_hash_functions
func (b *BloomFilter) CalculatingMK(num uint, fpp float64) {
	nFloat := float64(num)
	mFloat := -nFloat * math.Log(fpp) / math.Pow(math.Ln2, 2)
	b.m = int(mFloat)
	b.k = int(math.Ln2 * (math.Ceil(mFloat) / nFloat))
}

func (b *BloomFilter) hashWithIndex(idx int, data []byte) int {
	b.h.Reset()
	b.h.Write(data)
	hashres := b.h.Sum64()
	return int((uint64(idx) + hashres) % uint64(b.m))
}

func (b *BloomFilter) Add(element []byte) {
	for i := 0; i < b.k; i++ {
		hashIdx := b.hashWithIndex(i, element)
		b.bfList[hashIdx]++
	}
}

func (b *BloomFilter) Exist(element []byte) bool {
	for i := 0; i < b.k; i++ {
		hashIdx := b.hashWithIndex(i, element)
		if b.bfList[hashIdx] <= 0 {
			return false
		}
	}
	return true
}

func (b *BloomFilter) Remove(element []byte) {
	if !b.Exist(element) {
		return
	}

	for i := 0; i < b.k; i++ {
		hashIdx := b.hashWithIndex(i, element)
		if b.bfList[hashIdx] > 0 {
			b.bfList[hashIdx]--
		}
	}
}
