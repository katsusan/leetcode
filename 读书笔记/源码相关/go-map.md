
# 1. map的结构

```Go
// 编译时生成
type maptype struct {
	typ    _type
	key    *_type
	elem   *_type
	bucket *_type // internal type representing a hash bucket
	// function for hashing keys (ptr to key, seed) -> hash
	hasher     func(unsafe.Pointer, uintptr) uintptr
	keysize    uint8  // size of key slot
	elemsize   uint8  // size of elem slot
	bucketsize uint16 // size of bucket
	flags      uint32
}

// A header for a Go map.
type hmap struct {
	// Note: the format of the hmap is also encoded in cmd/compile/internal/gc/reflect.go.
	// Make sure this stays in sync with the compiler's definition.
	count     int // # live cells == size of map.  Must be first (used by len() builtin)
	flags     uint8
	B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
	noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
	hash0     uint32 // hash seed

	buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
	oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
	nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)

	extra *mapextra // optional fields
}

// A bucket for a Go map.
type bmap struct {
	// tophash generally contains the top byte of the hash value
	// for each key in this bucket. If tophash[0] < minTopHash,
	// tophash[0] is a bucket evacuation state instead.
	tophash [bucketCnt]uint8
	// Followed by bucketCnt keys and then bucketCnt elems.
	// NOTE: packing all the keys together and then all the elems together makes the
	// code a bit more complicated than alternating key/elem/key/elem/... but it allows
	// us to eliminate padding which would be needed for, e.g., map[int64]int8.
	// Followed by an overflow pointer.
}

// A hash iteration structure.
// If you modify hiter, also change cmd/compile/internal/gc/reflect.go to indicate
// the layout of this structure.
type hiter struct {
	key         unsafe.Pointer // Must be in first position.  Write nil to indicate iteration end (see cmd/internal/gc/range.go).
	elem        unsafe.Pointer // Must be in second position (see cmd/internal/gc/range.go).
	t           *maptype
	h           *hmap
	buckets     unsafe.Pointer // bucket ptr at hash_iter initialization time
	bptr        *bmap          // current bucket
	overflow    *[]*bmap       // keeps overflow buckets of hmap.buckets alive
	oldoverflow *[]*bmap       // keeps overflow buckets of hmap.oldbuckets alive
	startBucket uintptr        // bucket iteration started at
	offset      uint8          // intra-bucket offset to start from during iteration (should be big enough to hold bucketCnt-1)
	wrapped     bool           // already wrapped around from end of bucket array to beginning
	B           uint8
	i           uint8
	bucket      uintptr
	checkBucket uintptr
}

```

# 2. 所有map相关函数签名

```Go
// map.go
// 对于fast/slow与fat的选择： 
//  - 先判断出是否要fast即mapaccessN_fastxx还是mapaccessN
//  - 然后再视v的宽度与1024的关系判断是mapaccessN{shoudfast}还是mapccessN_fat
func makemap(t *maptype, hint int, h *hmap) *hmap   // hint > 8
func makemap_small() *hmap      // hint省略或<=8
func makemap64(t *maptype, hint int64, h *hmap) *hmap   //32位环境用64位的hint时
func mapaccess1(t *maptype, h *hmap, key unsafe.Pointer) unsafe.Pointer     // v := map[key]
func mapaccess2(t *maptype, h *hmap, key unsafe.Pointer) (unsafe.Pointer, bool) // v, ok := map[key]
func mapaccessK(t *maptype, h *hmap, key unsafe.Pointer) (unsafe.Pointer, unsafe.Pointer)
func mapaccess1_fat(t *maptype, h *hmap, key, zero unsafe.Pointer) unsafe.Pointer   // v := map[key] && sizeof(v) > 1024
func mapaccess2_fat(t *maptype, h *hmap, key, zero unsafe.Pointer) (unsafe.Pointer, bool) //v, ok := map[key] && sizeof(v) > 1024
func mapassign(t *maptype, h *hmap, key unsafe.Pointer) unsafe.Pointer  // m[key] = v
func mapdelete(t *maptype, h *hmap, key unsafe.Pointer)     //del(m, key)
func mapclear(t *maptype, h *hmap)

// map_faststr.go
func mapaccess1_faststr(t *maptype, h *hmap, ky string) unsafe.Pointer
func mapaccess2_faststr(t *maptype, h *hmap, ky string) (unsafe.Pointer, bool)
func mapassign_faststr(t *maptype, h *hmap, s string) unsafe.Pointer
func mapdelete_faststr(t *maptype, h *hmap, ky string)

// map_fast64.go
func mapaccess1_fast64(t *maptype, h *hmap, key uint64) unsafe.Pointer
func mapaccess2_fast64(t *maptype, h *hmap, key uint64) (unsafe.Pointer, bool)
func mapassign_fast64(t *maptype, h *hmap, key uint64) unsafe.Pointer
func mapassign_fast64ptr(t *maptype, h *hmap, key unsafe.Pointer) unsafe.Pointer
func mapdelete_fast64(t *maptype, h *hmap, key uint64)

// map_fast32.go
func mapaccess1_fast32(t *maptype, h *hmap, key uint32) unsafe.Pointer
func mapaccess2_fast32(t *maptype, h *hmap, key uint32) (unsafe.Pointer, bool)
func mapassign_fast32(t *maptype, h *hmap, key uint32) unsafe.Pointer
func mapassign_fast32ptr(t *maptype, h *hmap, key unsafe.Pointer) unsafe.Pointer
func mapdelete_fast32(t *maptype, h *hmap, key uint32)
```

# 4. fast优化

golang对符合特定条件的map会有额外的fast路径优化，如下：
    - 如果map<k,v>中的v大小超过128字节则不走fast
    - 对于key长
        + 如果为32位的key
            * key中不存在指向堆的指针则走mapfast32
            * key中存在指向堆的指针则走mapfastptr
        + 如果为64位的key
            * key中不存在指向堆的指针则走mapfast64
            * key中存在指向堆的指针则走mapfastptr
        + 如果key为string
            走mapfaststr
    - 其它情况不走fast

```Go
// walk.go
func mapfast(t *types.Type) int {
	// Check runtime/map.go:maxElemSize before changing.
	if t.Elem().Width > 128 {
		return mapslow
	}
	switch algtype(t.Key()) {
	case AMEM32:
		if !t.Key().HasPointers() {
			return mapfast32
		}
		if Widthptr == 4 {
			return mapfast32ptr
		}
		Fatalf("small pointer %v", t.Key())
	case AMEM64:
		if !t.Key().HasPointers() {
			return mapfast64
		}
		if Widthptr == 8 {
			return mapfast64ptr
		}
		// Two-word object, at least one of which is a pointer.
		// Use the slow path.
	case ASTRING:
		return mapfaststr
	}
	return mapslow
}
```

# 5. map中具体的设计

注：编译时maptype中hasher函数由genhash(t.Key())确定，bucketsize由dowidth(bucket)确定，
    其中每个bucket可以看作一个结构体，由以下成员构成，padding和结构体对齐规则类似。(目前bucketCnt为常量8)
        - tophash [8]uint8  
        - key x 8   // key的size x 8，比如int64的key就是8x8
        - elem x 8  // elem的size x 8， 比如string的elem就是16x8
        - overflow uintptr  //溢出时不一定会马上扩容，此时overflow代表溢出时的额外bmap
    比如说，map[int]string的bucketsize为8+8x8+16x8+8=208字节。
    另外tophash[i]代表第i个key的前8位，所以每次先遍历比较tophash与计算所得hash的前8位，相符再取对应的key和elem操作。

key
↓
hash = maptype.hasher(key, maptype.hash0)
↓
bucketindex = hash & (2<<B-1)
↓
b = hmap.bucket + maptype.bucketsize * bucketindex  // 获取桶的地址
↓
按照以下逻辑遍历桶：
    - 比较hash的高8位(特别地，小于0x5时有特殊含义)与桶的tophash[i]，不相等则说明key不匹配
    - 调用maptype.key.equal(key, k)判断给定key和桶里的k是否相等，相等则匹配成功根据i计算出elem的地址并返回
    - 否则根据访问形式返回空值或 空值,false

关于hash算法，以x64为例：
	runtim.memhash64 //通用的64位key的hash算法，由编译时的genhash返回
		- 若runtime.useAeshash为true则使用硬件的aesenc指令计算
		- 否则调用runtime.memhash64Fallback (基于xxhash和cityhash的自定义算法)


# 6. map扩容

当向map里写入新的key/value时，满足以下两个条件之一发生扩容
    - overLoadFactor(h.count+1, h.B)为true  // 即 h.count+1 > 8 && h.count > 6.5 * 2^B
    - tooManyOverflowBuckets(h.noverflow, h.B)为true
        + 若h.B<15，则返回 h.noverflow >= 2^B
        + 若h.B>=15, 则返回 h.noverflow >= 2^15
      简单地说桶个数低于32768时溢出桶数不能超过桶个数，桶个数大于32768时溢出桶数不能大于32768，否则就扩容

扩容有两种形式: SameSize Grow和Double Grow，即原尺寸调整和双倍扩容，
若触发了负载因子上限导致扩容会进行双倍扩容，否则进行原尺寸调整，将稀疏的元素聚集在一起。
(evacuate函数比较复杂有空再梳理)

# 7. map遍历

在栈上申请一个it *hiter变量,调用mapiterinit(t *maptype, h *hmap, it *hiter)初始化遍历结构体
(遍历map的随机因子就在此产生)，里面最后调用了func mapiternext(it *hiter)，当遍历了一圈后将
it.key和it.elem都设为nil，这样外面的range语法糖就可以通过检测it.key是否为nil来判断遍历是否结束。








