/*
* 参考：https://github.com/golang/proposal/blob/master/design/34481-opencoded-defers.md
*/
# 1 Defer
# 1.1 golang的defer实现按照演化进程分为3个阶段：
    在defer语句处插入下述三种函数
    - runtime.deferproc：Go1.13以前的实现，所有_defer结构体都在堆上分配，即newdefer(siz)，它依次从当前p的deferpool，全局deferpool，堆上来尝试申请_defer的空间。
    - runtime.deferprocStack：Go1.13开始在栈上分配_defer的空间，用deferprocStack来将gp._defer指向栈上的该地址，不仅cache友好还免去了可能现分配内存的迟滞。
    - inline code：Go1.14开始如果满足一定条件会在返回语句之前直接插入defer代码，条件如下：
        + 编译选项里没有加入-N(即禁止优化)
        + defer语句不超过8个，且defer语句数与return数乘积不超过15
        + defer不在循环体内

    然后如果不是open-coded的话在每个return之后执行runtime.deferreturn

# 1.2 非open-coded情况下的defer执行
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

# 1.3 deferprocStack的执行模型
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


# 2 Recover

panic(x) → r := recove(), r = x
运行时rewrite为runtime.gorecover，如下：

```Go
//go:nosplit
func gorecover(argp uintptr) interface{} {
	// Must be in a function running as part of a deferred call during the panic.
	// Must be called from the topmost function of the call
	// (the function used in the defer statement).
	// p.argp is the argument pointer of that topmost deferred function call.
	// Compare against argp reported by caller.
	// If they match, the caller is the one who can recover.
	gp := getg()
	p := gp._panic
	if p != nil && !p.goexit && !p.recovered && argp == uintptr(p.argp) {
		p.recovered = true
		return p.arg
	}
	return nil
}
```

# 3 Panic
// panic descriptor
```Go
// A _panic holds information about an active panic.
//
// A _panic value must only ever live on the stack.
//
// The argp and link fields are stack pointers, but don't need special
// handling during stack growth: because they are pointer-typed and
// _panic values only live on the stack, regular stack pointer
// adjustment takes care of them.
type _panic struct {
	argp      unsafe.Pointer // pointer to arguments of deferred call run during panic; cannot move - known to liblink
	arg       interface{}    // argument to panic
	link      *_panic        // link to earlier panic
	pc        uintptr        // where to return to in runtime if this panic is bypassed
	sp        unsafe.Pointer // where to return to in runtime if this panic is bypassed
	recovered bool           // whether this panic is over
	aborted   bool           // the panic was aborted
	goexit    bool
}
```

panic关键字编译时会被rewrite成runtime.gopanic。
代码比较长，逻辑如下：
- 函数入口：`func gopanic(e interface{})`
- 首先判断如下4个条件，不满足的话抛出unrecoverable异常(throw)。
	+ gp.m.curg != gp					// 不能在g0栈上panic
	+ gp.m.mallocing != 0			// 不能在内存分配时panic
	+ gp.m.preemptoff	!= ""		// 禁止preempt时不能panic
	+ gp.m.locks != 0					// m被绑定时不能panic，如acquirem函数

- 把传入的参数e作为arg的_panic结构插入当前协程gp的_panic链表头，p必定在栈上
```Go
	var p _panic
	p.arg = e				
	p.link = gp._panic
	gp._panic = (*_panic)(noescape(unsafe.Pointer(&p)))
```

- 执行当前协程的defer task
```Go
	for {
		d := gp._defer
		if d == nil {			// 没有defer的话无法执行recover，直接break出去执行fatapanic
			break
		}

		if d.started {		// defer was started by earlier panic or Goexit
			if d._panic != nil {
				d._panic.aborted = true
			}
			d._panic = nil
			if !d.openDefer {
				d.fn = nil
				gp._defer = d.link
				freedefer(d)
				continue
			}
		}
		
		d.started = true	// Mark defer as started, but keep on list
		d._panic = (*_panic)(noescape(unsafe.Pointer(&p)))	// Record the panic that is running the defer, as panic may happen in defer.

		// do defer call
		// defer里带有recover的话会把gp._panic.recovered设为true，参考gorecover函数实现
		done := true
		if d.openDefer {
			done = runOpenDeferFrame(gp, d)
			if done && !d._panic.recovered {
				addOneOpenDeferFrame(gp, 0, nil)
			}
		} else {
			p.argp = unsafe.Pointer(getargp(0))
			reflectcall(nil, unsafe.Pointer(d.fn), deferArgs(d), uint32(d.siz), uint32(d.siz))
		}

		p.argp = nil		// p.argp已经被recover取得，赋为nil让其GC

		// reflectcall did not panic.
		if gp._defer != d {
			throw("bad defer entry in panic")
		}

		// 头个defer已经执行完毕可以free
		d._panic = nil

		pc := d.pc
		sp := unsafe.Pointer(d.sp) // must be pointer so it gets adjusted during stack copy
		if done {
			d.fn = nil
			gp._defer = d.link	// 指向次个defer
			freedefer(d)
		}		
		
		if p.recovered {		// 如果defer里执行了recover
			gp._panic = p.link
			if gp._panic != nil && gp._panic.goexit && gp._panic.aborted {
				// A normal recover would bypass/abort the Goexit.  Instead,
				// we return to the processing loop of the Goexit.
				gp.sigcode0 = uintptr(gp._panic.sp)
				gp.sigcode1 = uintptr(gp._panic.pc)
				mcall(recovery)
				throw("bypassed recovery failed") // mcall should not return
			}
			atomic.Xadd(&runningPanicDefers, -1)

			// remove any remaining non-started, open-coded defer
			d := gp._defer
			var prev *_defer
			if !done {
				// Skip our current frame, if not done. It is
				// needed to complete any remaining defers in
				// deferreturn()
				prev = d
				d = d.link
			}
			for d != nil {
				if d.started {
					// This defer is started but we
					// are in the middle of a
					// defer-panic-recover inside of
					// it, so don't remove it or any
					// further defer entries
					break		// defer is started
				}
				if d.openDefer {
					if prev == nil {
						gp._defer = d.link
					} else {
						prev.link = d.link
					}
					newd := d.link
					freedefer(d)
					d = newd
				} else {
					prev = d
					d = d.link
				}
			}

			gp._panic = p.link
			// Aborted panics are marked but remain on the g.panic list.
			// Remove them from the list.
			for gp._panic != nil && gp._panic.aborted {
				gp._panic = gp._panic.link
			}
			if gp._panic == nil { // must be done with signal
				gp.sig = 0
			}
			// Pass information about recovering frame to recovery.
			gp.sigcode0 = uintptr(sp)
			gp.sigcode1 = pc
			mcall(recovery)		// 切换到sp&pc的上下文
			throw("recovery failed") // mcall should not return
		}
	}


```

- 没有recover的话直接
```Go
	preprintpanics(gp._panic)		// 把gp._panic.arg具体化到Error或String方法(如果实现了的话)

	fatalpanic(gp._panic) // should not return，没有recover的话
	*(*int)(nil) = 0      // not reached
```