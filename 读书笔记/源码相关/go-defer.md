/*
* 参考：https://github.com/golang/proposal/blob/master/design/34481-opencoded-defers.md
*/

# 1. golang的defer实现按照演化进程分为3个阶段：
    在defer语句处插入下述三种函数
    - runtime.deferproc：Go1.13以前的实现，所有_defer结构体都在堆上分配，即newdefer(siz)，它依次从当前p的deferpool，全局deferpool，堆上来尝试申请_defer的空间。
    - runtime.deferprocStack：Go1.13开始在栈上分配_defer的空间，用deferprocStack来将gp._defer指向栈上的该地址，不仅cache友好还免去了可能现分配内存的迟滞。
    - inline code：Go1.14开始如果满足一定条件会在返回语句之前直接插入defer代码，条件如下：
        + 编译选项里没有加入-N(即禁止优化)
        + defer语句不超过8个，且defer语句数与return数乘积不超过15
        + defer不在循环体内

    然后如果不是open-coded的话在每个return之后执行runtime.deferreturn

# 2. 非open-coded情况下的defer执行
    golang用_defer结构体来表示一次defer执行,定义如下
```Go
// A _defer holds an entry on the list of deferred calls.
// If you add a field here, add code to clear it in freedefer and deferProcStack
// This struct must match the code in cmd/compile/internal/gc/reflect.go:deferstruct
// and cmd/compile/internal/gc/ssa.go:(*state).call.
// Some defers will be allocated on the stack and some on the heap.
// All defers are logically part of the stack, so write barriers to
// initialize them are not required. All defers must be manually scanned,
// and for heap defers, marked.
type _defer struct {
	siz     int32 // includes both arguments and results
	started bool
	heap    bool
	// openDefer indicates that this _defer is for a frame with open-coded
	// defers. We have only one defer record for the entire frame (which may
	// currently have 0, 1, or more defers active).
	openDefer bool
	sp        uintptr  // sp at time of defer
	pc        uintptr  // pc at time of defer
	fn        *funcval // can be nil for open-coded defers
	_panic    *_panic  // panic that is running defer
	link      *_defer

	// If openDefer is true, the fields below record values about the stack
	// frame and associated function that has the open-coded defer(s). sp
	// above will be the sp for the frame, and pc will be address of the
	// deferreturn call in the function.
	fd   unsafe.Pointer // funcdata for the function associated with the frame
	varp uintptr        // value of varp for the stack frame
	// framepc is the current pc associated with the stack frame. Together,
	// with sp above (which is the sp associated with the stack frame),
	// framepc/sp can be used as pc/sp pair to continue a stack trace via
	// gentraceback().
	framepc uintptr
}
```

# 3. deferprocStack的执行模型
    函数声明为`func deferprocStack(d *_defer)`，此时_defer结构体在栈上分配，
    而_defer.pc/sp由函数deferprocStack内设定，调用者需要处理的是siz、fn以及紧跟在_defer指针后的参数，

    调用者的栈       地址        值
    |__fn___|   =>  sp+0x30     fn  //defer所执行函数的地址
    |__...__|
    |__*d___|   =>  sp+0x18     0x10    //openDefer+heap+started+siz共同占用第一个8字节，我们只需传入低位的4字节siz即0x10代表16字节的参数
    |__arg1_|   =>  sp+0x10
    |__arg2_|   =>  sp+0x8
    |___d___|   =>  sp:         sp+0x18 //此处是d的指针，解引用后的*d即上面的+0x18处，代表栈上_defer结构体的开始 
    <call runtime.deferprocStack>






