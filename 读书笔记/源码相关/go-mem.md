ref: 
	cnblogs.com/saryli/p/10116579.html
	https://medium.com/a-journey-with-go/go-memory-management-and-allocation-a7396d430f44

	https://groups.google.com/g/golang-nuts/c/umeli6cYXT4/m/y6A92G5oBQAJ
	https://go-review.googlesource.com/c/go/+/205239/		//关于m和p各自拥有一个mcache的问题

1. go在启动时会分配一段连续的虚拟地址空间，分为三个区域：
        spans   ->    bitmap    ->    arena
(X64)   512MB         16GB            512GB

    arena: 通常所说的堆区

    bitmap:  用于描述arena，1个字节描述arean区4个指针大小的内存(X64下为8字节),
            每2位描述1个指针大小的内存，包括should scan和is pointer两个信息。
            shouldscan4,~3,~2,~1,ispointer1,~2,~3,~4这样从中间开始每对bit描述
            一个对象。

            bitmap              |       arena
            byte3, byte2, byte1 | byte1-slot1,byte1-slot2,byte1-slot3,byte1-slot4,...
                            从中间向两边对应
    
    spans:  表示arena区的某一页属于哪个span。该区一个指针8byte对应arena区的一页(8KB)
            *mspan1, *mspan2, *mspan3,... | bitmap | page1, page2, page3,...


2. Golang的内存分配
	runtime的内存管理粒度是8KB or bigger，运行时按照这个一定的size将之分为多个mspan结构，
	mspan在mheap.go里定义：
	```Go
	type mspan struct {
		next *mspan     // next span in list, or nil if none
		prev *mspan     // previous span in list, or nil if none
		list *mSpanList // For debugging. TODO: Remove.

		startAddr uintptr // address of first byte of span aka s.base()
		npages    uintptr // number of pages in span
		spanclass   spanClass     // size class and noscan (uint8)
		... //剩余一些gc时的mark/sweep字段
	}
	```
	#mspan具体的size可以在sizeclasses.go里，一共67种，从8到32768#
	// class  bytes/obj  bytes/span  objects  tail waste  max waste
	//     1          8        8192     1024           0     87.50%
	//     2         16        8192      512           0     43.75%
	//     3         32        8192      256           0     46.88%
	//     4         48        8192      170          32     31.52%
	...
	//    63      24576       24576        1           0     11.45%
	//    64      27264       81920        3         128     10.00%
	//    65      28672       57344        2           0      4.91%
	//    66      32768       32768        1           0     12.50%

	由上面的next和prev指针可以看出，mspan结构应该是双向链表中的节点类型->即mSpanList。
	
	mcache：
		类似于tcmalloc，Golang为每个P提供了Thread Local Cache,也就是mcache。
		如果goroutine需要分配内存，会先从mcache上拿取，由于单P上的执行可以看做是串行的，
		所以不需要加锁，一定程度上加快了内存分配的效率。

		mcache包含了所有不同size的mspan作为cache。

	```Go
	//mcache.go
	type struct mcache {
		// spans to allocate from, indexed by spanClass
		// numSpanClasses=67x2
		alloc [numSpanClasses]*mspan 
		...
	}
	```
		每个size的mspan包含2种，scan和noscan类型。
		- scan：包含指针的对象
		- noscan：不包含指针的对象
		这样分类可以在GC扫描时省略一半的扫描量，因为不需要确认noscan里是否有指向其它对象的指针。

		总结：小于32K的对象直接按照size分配给mcache里具体的mspan，
			当mcache里mspan空间不足时，则从mcentral里获取新的mspan。

	mcentral：指定spanClass/size的内存链表，分为noempty和empty两部分。
		- nonempty：至少有一个节点尚未分配的mspan链表、或者已经缓存在mcache里的mspan链表
		- empty：已分配完的mspan链表
		当mcentral里请求新的span时，会先从nonempty链表里选取，然后把这个节点放入empty链表。
		反过来讲，当empty链表里的mspan对象被free掉后会被放入nonempty链表。
		每个spanclass/size的mcentral由mheap维护。

	```Go
	type mcentral struct {
		lock      mutex
		spanclass spanClass
		nonempty  mSpanList // list of spans with a free object, ie a nonempty free list
		empty     mSpanList // list of spans with no free objects (or cached in an mcache)

		// nmalloc is the cumulative count of objects allocated from
		// this mcentral, assuming all spans in mcaches are
		// fully-allocated. Written atomically, read under STW.
		nmalloc uint64
	}
	type mheap struct {
		//...
		central [numSpanClasses]struct {	//numSpanClasses为134
			mcentral mcentral
			//CacheLinePadSize依cpu架构不同而异，x64上为64字节，设定pad可以使得一次lock的时间内获得cache line
			//避免false sharing问题
			pad      [cpu.CacheLinePadSize - unsafe.Sizeof(mcentral{})%cpu.CacheLinePadSize]byte
		}
	}
	```

	mheap：描述了golang管理运行时堆空间的结构，runtime.mheap_是运行时的全局唯一实例。
	mcentral空间不够时会向mheap申请指定size的span。

	- 大于32K的对象(large object)直接由mheap分配
	- 小于16B的对象(tiny object)由mcache的tiny allocator分配
	- 16B到32K的对象按照指定的size class由mcache分配
	- 如若mcache中没有足够的memory block则向mcentral申请
	- 如若mcentral中没有可用的memory block则向mheap申请并使用BestFit算法找到最合适的span
	- mheap中没有合适的span则向OS申请页(至少1MB)

	Go以arena的尺度分配内存，一系列arena的集合则表现为运行时的堆。

	```Go
	type mheap struct {
		// lock must only be acquired on the system stack, otherwise a g
		// could self-deadlock if its stack grows with the lock held.
		lock      mutex
		pages     pageAlloc // page allocation data structure
		sweepgen  uint32    // sweep generation, see comment in mspan; written during STW
		sweepdone uint32    // all spans are swept
		sweepers  uint32    // number of active sweepone calls
		allspans []*mspan 	//allspans is a slice of all mspans ever created. Each mspan appears exactly once.
		sweepSpans [2]gcSweepBuf 	//sweepSpans contains two mspan stacks: one of swept in-use spans, and one of unswept in-use spans.
		pagesInUse         uint64  // pages of spans in stats mSpanInUse; updated atomically
		pagesSwept         uint64  // pages swept this cycle; updated atomically
		pagesSweptBasis    uint64  // pagesSwept to use as the origin of the sweep ratio; updated atomically
		sweepHeapLiveBasis uint64  // value of heap_live to use as the origin of sweep ratio; written with lock, read without
		sweepPagesPerByte  float64 // proportional sweep ratio; written with lock, read without
		scavengeGoal uint64		//the amount of total retained heap memory(measured by heapRetained) that the runtime will try to maintain by returning memory to the OS.
		reclaimIndex uint64		//the page index in allArenas of next page to reclaim, If this is >= 1<<63, the page reclaimer is done scanning the page marks.
		reclaimCredit uintptr	//reclaimCredit is spare credit for extra pages swept

		// Malloc stats.
		largealloc  uint64                  // bytes allocated for large objects
		nlargealloc uint64                  // number of large object allocations
		largefree   uint64                  // bytes freed for large objects (>maxsmallsize)
		nlargefree  uint64                  // number of frees for large objects (>maxsmallsize)
		nsmallfree  [_NumSizeClasses]uint64 // number of frees for small objects (<=maxsmallsize)

		//2-level mapping of L1 map and many L2 maps
		//arenaL1Bits = 6, arenaL1Bits需要在BSS段占用 PtrSize*(1<<arenaL1Bits)的大小所以不能太大
		//arenaL2Bits = 20，即总共(2^6)*(2^20)的二维heapArena数组
		arenas [1 << arenaL1Bits]*[1 << arenaL2Bits]*heapArena	

		heapArenaAlloc linearAlloc	//pre-reserved space for allocating heapArena objects
		arenaHints *arenaHint		//a list of addresses at which to attempt to add more heap arenas
		arena linearAlloc		//a pre-reserved space for allocating heap arenas(the actual arenas)
		allArenas []arenaIdx 	//the arenaIndex of every mapped arena, can be used to iterate through the address space
		sweepArenas []arenaIdx	// a snapshot of allArenas taken at the beginning of the sweep cycle

		// curArena is the arena that the heap is currently growing
		// into. This should always be physPageSize-aligned.
		curArena struct {
			base, end uintptr
		}

		_ uint32 // ensure 64-bit alignment of central

		// central free lists for small size classes.
		// the padding makes sure that the mcentrals are
		// spaced CacheLinePadSize bytes apart, so that each mcentral.lock
		// gets its own cache line.
		// central is indexed by spanClass.
		central [numSpanClasses]struct {
			mcentral mcentral
			pad      [cpu.CacheLinePadSize - unsafe.Sizeof(mcentral{})%cpu.CacheLinePadSize]byte
		}

		spanalloc             fixalloc // allocator for span*
		cachealloc            fixalloc // allocator for mcache*
		specialfinalizeralloc fixalloc // allocator for specialfinalizer*
		specialprofilealloc   fixalloc // allocator for specialprofile*
		speciallock           mutex    // lock for special record allocators.
		arenaHintAlloc        fixalloc // allocator for arenaHints

		unused *specialfinalizer // never set, just here to force the specialfinalizer type into DWARF
	}
	```


