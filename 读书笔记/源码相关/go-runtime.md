https://tonybai.com/2017/06/23/an-intro-about-goroutine-scheduler/
https://www.cnblogs.com/yjf512/archive/2012/07/17/2595689.html
https://medium.com/@ankur_anand/illustrated-tales-of-go-runtime-scheduler-74809ef6d19b
https://blog.learngoprogramming.com/a-visual-guide-to-golang-memory-allocator-from-ground-up-e132258453ed

//pre: debug调度信息
设置GODEBUG=schedtrace=DURATION来定时输出scheduler日志，如GODEBUG=schedtrace=5 go run main.go
或者 GODEBUG=scheddetail=1,schedtrace=1000 ./program
=> debug信息解析参考：https://www.ardanlabs.com/blog/2015/02/scheduler-tracing-in-go.html

//pre: runtime中的全局变量
```Go
//runtime2.go
var (
	var memstats mstats		//MemStats records statistics about the memory allocator , in mstats.go

	allglen    uintptr
	allm       *m
	allp       []*p  // len(allp) == gomaxprocs; may change at safe points, otherwise immutable
	allpLock   mutex // Protects P-less reads of allp and all writes
	gomaxprocs int32
	ncpu       int32
	forcegc    forcegcstate
	sched      schedt
	newprocs   int32

	// Information about what cpu features are available.
	// Packages outside the runtime should not use these
	// as they are not an external api.
	// Set on startup in asm_{386,amd64}.s
	processorVersionInfo uint32
	isIntel              bool
	lfenceBeforeRdtsc    bool

	goarm                uint8 // set by cmd/link on arm systems
	framepointer_enabled bool  // set by cmd/link
)

```


1. go的G-P-M模型定义

//D:\Dev\Go\src\runtime\runtime2.go
```Go
type g struct {
    // Stack parameters.
    // stack describes the actual stack memory: [stack.lo, stack.hi).
    // stackguard0 is the stack pointer compared in the Go stack growth prologue.
    // It is stack.lo+StackGuard normally, but can be StackPreempt to trigger a preemption.
    // stackguard1 is the stack pointer compared in the C stack growth prologue.
    // It is stack.lo+StackGuard on g0 and gsignal stacks.
    // It is ~0 on other goroutine stacks, to trigger a call to morestackc (and crash).
    stack       stack   // offset known to runtime/cgo
    stackguard0 uintptr // offset known to liblink
    stackguard1 uintptr // offset known to liblink

    _panic         *_panic // innermost panic - offset known to liblink
    _defer         *_defer // innermost defer
    m              *m      // current m; offset known to arm liblink
    sched          gobuf
    syscallsp      uintptr        // if status==Gsyscall, syscallsp = sched.sp to use during gc
    syscallpc      uintptr        // if status==Gsyscall, syscallpc = sched.pc to use during gc
    stktopsp       uintptr        // expected sp at top of stack, to check in traceback
    param          unsafe.Pointer // passed parameter on wakeup
    atomicstatus   uint32
    stackLock      uint32 // sigprof/scang lock; TODO: fold in to atomicstatus
    goid           int64
    schedlink      guintptr
    waitsince      int64      // approx time when the g become blocked
    waitreason     waitReason // if status==Gwaiting
    preempt        bool       // preemption signal, duplicates stackguard0 = stackpreempt
    paniconfault   bool       // panic (instead of crash) on unexpected fault address
    preemptscan    bool       // preempted g does scan for gc
    gcscandone     bool       // g has scanned stack; protected by _Gscan bit in status
    gcscanvalid    bool       // false at start of gc cycle, true if G has not run since last scan; TODO: remove?
    throwsplit     bool       // must not split stack
    raceignore     int8       // ignore race detection events
    sysblocktraced bool       // StartTrace has emitted EvGoInSyscall about this goroutine
    sysexitticks   int64      // cputicks when syscall has returned (for tracing)
    traceseq       uint64     // trace event sequencer
    tracelastp     puintptr   // last P emitted an event for this goroutine
    lockedm        muintptr
    sig            uint32
    writebuf       []byte
    sigcode0       uintptr
    sigcode1       uintptr
    sigpc          uintptr
    gopc           uintptr         // pc of go statement that created this goroutine
    ancestors      *[]ancestorInfo // ancestor information goroutine(s) that created this goroutine (only used if debug.tracebackancestors)
    startpc        uintptr         // pc of goroutine function
    racectx        uintptr
    waiting        *sudog         // sudog structures this g is waiting on (that have a valid elem ptr); in lock order
    cgoCtxt        []uintptr      // cgo traceback context
    labels         unsafe.Pointer // profiler labels
    timer          *timer         // cached timer for time.Sleep
    selectDone     uint32         // are we participating in a select and did someone win the race?

    // Per-G GC state

    // gcAssistBytes is this G's GC assist credit in terms of
    // bytes allocated. If this is positive, then the G has credit
    // to allocate gcAssistBytes bytes without assisting. If this
    // is negative, then the G must correct this by performing
    // scan work. We track this in bytes to make it fast to update
    // and check for debt in the malloc hot path. The assist ratio
    // determines how this corresponds to scan work debt.
    gcAssistBytes int64
}
```

```Go
type m struct {
	g0      *g     // goroutine with scheduling stack
	morebuf gobuf  // gobuf arg to morestack
	divmod  uint32 // div/mod denominator for arm - known to liblink

	// Fields not known to debuggers.
	procid        uint64       // for debuggers, but offset not hard-coded
	gsignal       *g           // signal-handling g
	goSigStack    gsignalStack // Go-allocated signal handling stack
	sigmask       sigset       // storage for saved signal mask
	tls           [6]uintptr   // thread-local storage (for x86 extern register)
	mstartfn      func()
	curg          *g       // current running goroutine
	caughtsig     guintptr // goroutine running during fatal signal
	p             puintptr // attached p for executing go code (nil if not executing go code)
	nextp         puintptr
	oldp          puintptr // the p that was attached before executing a syscall
	id            int64
	mallocing     int32
	throwing      int32
	preemptoff    string // if != "", keep curg running on this m
	locks         int32
	dying         int32
	profilehz     int32
	spinning      bool // m is out of work and is actively looking for work
	blocked       bool // m is blocked on a note
	newSigstack   bool // minit on C thread called sigaltstack
	printlock     int8
	incgo         bool   // m is executing a cgo call
	freeWait      uint32 // if == 0, safe to free g0 and delete m (atomic)
	fastrand      [2]uint32
	needextram    bool
	traceback     uint8
	ncgocall      uint64      // number of cgo calls in total
	ncgo          int32       // number of cgo calls currently in progress
	cgoCallersUse uint32      // if non-zero, cgoCallers in use temporarily
	cgoCallers    *cgoCallers // cgo traceback if crashing in cgo call
	park          note
	alllink       *m // on allm
	schedlink     muintptr
	mcache        *mcache
	lockedg       guintptr
	createstack   [32]uintptr // stack that created this thread.
	lockedExt     uint32      // tracking for external LockOSThread
	lockedInt     uint32      // tracking for internal lockOSThread
	nextwaitm     muintptr    // next m waiting for lock
	waitunlockf   func(*g, unsafe.Pointer) bool
	waitlock      unsafe.Pointer
	waittraceev   byte
	waittraceskip int
	startingtrace bool
	syscalltick   uint32
	thread        uintptr // thread handle
	freelink      *m      // on sched.freem

	// these are here because they are too large to be on the stack
	// of low-level NOSPLIT functions.
	libcall   libcall
	libcallpc uintptr // for cpu profiler
	libcallsp uintptr
	libcallg  guintptr
	syscall   libcall // stores syscall parameters on windows

	vdsoSP uintptr // SP for traceback while in VDSO call (0 if not in call)
	vdsoPC uintptr // PC for traceback while in VDSO call

	dlogPerM

	mOS
}
```

```Go
type p struct {
	id          int32
	status      uint32 // one of pidle/prunning/...
	link        puintptr
	schedtick   uint32     // incremented on every scheduler call
	syscalltick uint32     // incremented on every system call
	sysmontick  sysmontick // last tick observed by sysmon
	m           muintptr   // back-link to associated m (nil if idle)
	mcache      *mcache
	raceprocctx uintptr

	deferpool    [5][]*_defer // pool of available defer structs of different sizes (see panic.go)
	deferpoolbuf [5][32]*_defer

	// Cache of goroutine ids, amortizes accesses to runtime·sched.goidgen.
	goidcache    uint64
	goidcacheend uint64

	// Queue of runnable goroutines. Accessed without lock.
	runqhead uint32
	runqtail uint32
	runq     [256]guintptr
	// runnext, if non-nil, is a runnable G that was ready'd by
	// the current G and should be run next instead of what's in
	// runq if there's time remaining in the running G's time
	// slice. It will inherit the time left in the current time
	// slice. If a set of goroutines is locked in a
	// communicate-and-wait pattern, this schedules that set as a
	// unit and eliminates the (potentially large) scheduling
	// latency that otherwise arises from adding the ready'd
	// goroutines to the end of the run queue.
	runnext guintptr

	// Available G's (status == Gdead)
	gFree struct {
		gList
		n int32
	}

	sudogcache []*sudog
	sudogbuf   [128]*sudog

	tracebuf traceBufPtr

	// traceSweep indicates the sweep events should be traced.
	// This is used to defer the sweep start event until a span
	// has actually been swept.
	traceSweep bool
	// traceSwept and traceReclaimed track the number of bytes
	// swept and reclaimed by sweeping in the current sweep loop.
	traceSwept, traceReclaimed uintptr

	palloc persistentAlloc // per-P to avoid mutex

	_ uint32 // Alignment for atomic fields below

	// Per-P GC state
	gcAssistTime         int64    // Nanoseconds in assistAlloc
	gcFractionalMarkTime int64    // Nanoseconds in fractional mark worker (atomic)
	gcBgMarkWorker       guintptr // (atomic)
	gcMarkWorkerMode     gcMarkWorkerMode

	// gcMarkWorkerStartTime is the nanotime() at which this mark
	// worker started.
	gcMarkWorkerStartTime int64

	// gcw is this P's GC work buffer cache. The work buffer is
	// filled by write barriers, drained by mutator assists, and
	// disposed on certain GC state transitions.
	gcw gcWork

	// wbBuf is this P's GC write barrier buffer.
	//
	// TODO: Consider caching this in the running G.
	wbBuf wbBuf

	runSafePointFn uint32 // if 1, run sched.safePointFn at next safe point

	// Lock for timers. We normally access the timers while running
	// on this P, but the scheduler can also do it from a different P.
	timersLock mutex

	// Actions to take at some time. This is used to implement the
	// standard library's time package.
	// Must hold timersLock to access.
	timers []*timer

	// Number of timers in P's heap.
	// Modified using atomic instructions.
	numTimers uint32

	// Number of timerModifiedEarlier timers on P's heap.
	// This should only be modified while holding timersLock,
	// or while the timer status is in a transient state
	// such as timerModifying.
	adjustTimers uint32

	// Number of timerDeleted timers in P's heap.
	// Modified using atomic instructions.
	deletedTimers uint32

	// Race context used while executing timer functions.
	timerRaceCtx uintptr

	// preempt is set to indicate that this P should be enter the
	// scheduler ASAP (regardless of what G is running on it).
	preempt bool

	pad cpu.CacheLinePad
}
```


2. 启动流程
	程序启动时会先根据GOMAXPROCS来启动一定数量的P,其中假设P0执行main
	M0 -> P0 -> G0 -> main().
	idle P: P1->P2->P3->...->Pn	//其余空闲的P放于链表

3. runtime调度中的公平性fairness
    - Local run queue：运行大于10ms的goroutine被标记为preemptible(可抢占),抢占点有函数调用，内存分配等。
        //Go1.13以前如果没有进行上述抢占点的操作可以独占CPU资源
    - Global run queue: 每61次schedule会检测一次grq来确保
    
    ```Go
    //proc.go:schedule()
    func schedule() {
        ...
        if gp == nil {
			//schedtick：每调度一次会++
    		if _g_.m.p.ptr().schedtick%61 == 0 && sched.runqsize > 0 {
                lock(&sched.lock)
                gp = globrunqget(_g_.m.p.ptr(), 1)
                unlock(&sched.lock)
            }
        }
    }
    ```



4. runtime启动流程	
   //以GO1.13.10 linux-amd64为例

	1. _rt0_amd64_linux	//rt0_linux_amd64.s
	 	编译出的ELF文件的entrypoint即此函数的地址，内部直接调用了_rt0_amd64。
		TEXT _rt0_amd64_linux(SB),NOSPLIT,$-8
			JMP	_rt0_amd64(SB)

	2. _rt0_amd64	//asm_amd64.s
		大多数AMD64架构的启动程序，主要处理程序入参。此时RDI存argc，RSI存指向参数数组的指针
		TEXT _rt0_amd64(SB),NOSPLIT,$-8
			MOVQ	0(SP), DI	// argc
			LEAQ	8(SP), SI	// argv
			JMP	runtime·rt0_go(SB)
	
	3. runtime·rt0_go	//ams_amd64.s
		进行了栈参数复制、runtim.g0的栈参数设定、获取CPU信息、设定tls、
		调用的函数：
			- runtime·settls
			- runtime·check
			- runtime·args
			- runtime·osinit
			- runtime·schedinit
				+ tracebackinit()	
				+ moduledataverify()	//验证模块信息
				+ stackinit()			//初始化全局stack池(runtime.stackpool/runtime.stackLarge)
				+ mallocinit()			//内存分配器初始化，包括验证一些相关参数是否正确，然后调用_mheap.init(), TODO:弄清相关参数
				+ fastrandinit()		//初始化全局随机因子runtime.fastrandseed
				+ mcommoninit(_g_.m)	//初始化m的属性，m.id/m.fastrand/m.gsignal/m.alllink/m.cgoCallers(optional)
				+ cpuinit()				//解析GODEBUG环境变量
				+ alginit()				//386/amd64/arm64且支持AES的情况下初始化runtime.hashkey和runtime.aeskeysched
				+ modulesinit()			//读取所有加载的模块并存到runtime.modulesSlice
				+ typelinksinit()		//遍历所有模块并且建立typemap来对类型去重
				+ itabsinit()			//将所有模块的类型信息加入全局itab table
				+ msigsave(_g_.m)		//将当前线程的信号掩码设定为空并将之前的掩码存于m.sigmask
				+ goargs()				//利用runtime.argc/argv设定全局参数变量runtime.argslice([]string类型)
				+ goenvs()				//初始化全局环境变量runtime.envs([]string类型)
				+ parsedebugvars()		//从GODEBUG环境变量中读取相应值
				+ gcinit()				//初始化GC参数，包括触发百分比memstats.triggerRatio（默认87.5%，从GOGC中获取）以及一些标记状态值
				+ proc					//用GOMAXPROCS设定runtime.procs,然后调用procresize()来调整runtime.allp为相应长度

			- runtime.newproc(0, runtime.mainPC)	//runtime.mainPC为指向runtime.main的指针，此处相当于go runtime.main()
			- runtime·mstart			//m的入口点函数，设定当前g的栈边界stackguard0和stackguard1后调用runtime.mstart1
				+ mstart1()				
					* save(getcallerpc(), getcallersp())	//
					* asminit()			//amd64下为空函数
					* minit()			//所有m都需要调用的初始化操作,主要是信号相关
						- minitSignals
							+ minitSignalStack()	//sigaltstack, 设定线程的信号栈
							+ minitSignalMask()		//sigpromask, 设定线程的信号掩码
					* mstartm0()<optional>	//如果是m0的话，初始化信号处理函数
						- initsig(bool)		//初始化进程的信号处理handler
							+ setsig(i, funcPC(sighandler))
								- sigaction(i, &sa, nil)	//i代表信号编号,从0到64，sa.fn为信号处理函数
					* _g_.m.mstartfn()	//如果当前m的mstartfn存在则执行它
					* acquirep(_g_.m.nextp.ptr())	//将_g_.m.nextp与当前m绑定
					* schedule()		//选取适当的g执行
				+ mexit(osStack)		//退出当前线程m，不可直接调用

		TEXT runtime·rt0_go(SB),NOSPLIT,$0
		// copy arguments forward on an even stack
		// 此处复制参数可能是因为后面函数调用需要利用RSP、RSP+0x8?
		MOVQ	DI, AX		// argc
		MOVQ	SI, BX		// argv
		SUBQ	$(4*8+7), SP		// 2args 2auto
		ANDQ	$~15, SP
		MOVQ	AX, 16(SP)	//argc复制到RSP+0x10
		MOVQ	BX, 24(SP)	//argv复制到RSP+0x18

		// create istack out of the given (operating system) stack.
		// _cgo_init may update stackguard.
		MOVQ	$runtime·g0(SB), DI		//复制runtime.g0(*g指针)到RDI
		LEAQ	(-64*1024+104)(SP), BX	//RSP-0xff98 -> RBX
		MOVQ	BX, g_stackguard0(DI)	//runtime.g0.g_stackguard0 = RBX
		MOVQ	BX, g_stackguard1(DI)	//runtime.g0.g_stackguard1 = RBX
		MOVQ	BX, (g_stack+stack_lo)(DI)	//runtime.g0.stack.lo = RBX
		MOVQ	SP, (g_stack+stack_hi)(DI)	//runtime.g0.stack.hi = RSP

		// find out information about the processor we're on
		MOVL	$0, AX
		CPUID
		MOVL	AX, SI	//执行后，BX+DX+CX = "Genu" + "ineI" + "ntel" = "GenuineIntel"
		CMPL	AX, $0
		JE	nocpuinfo

		// Figure out how to serialize RDTSC.
		// On Intel processors LFENCE is enough. AMD requires MFENCE.
		// Don't know about the rest, so let's do MFENCE.
		CMPL	BX, $0x756E6547  // "Genu"
		JNE	notintel
		CMPL	DX, $0x49656E69  // "ineI"
		JNE	notintel
		CMPL	CX, $0x6C65746E  // "ntel"
		JNE	notintel
		MOVB	$1, runtime·isIntel(SB)		//经过上面的CPUID信息校对可以在此处判断为Intel CPU
		MOVB	$1, runtime·lfenceBeforeRdtsc(SB)	//设定runtime.lfenceBeforeRdtsc标志为true
	notintel:

		// Load EAX=1 cpuid flags
		MOVL	$1, AX
		CPUID
		MOVL	AX, runtime·processorVersionInfo(SB)	//AX=1时返回CPU的Family ID, Model, Stepping ID
	nocpuinfo:
		// if there is an _cgo_init, call it.
		MOVQ	_cgo_init(SB), AX
		TESTQ	AX, AX
		JZ	needtls
		// arg 1: g0, already in DI
		MOVQ	$setg_gcc<>(SB), SI // arg 2: setg_gcc
		...
	needtls:
		LEAQ	runtime·m0+m_tls(SB), DI	//把runtime.m0.tls的地址复制到RDI
		CALL	runtime·settls(SB)	//针对传入的地址进行系统调用arch_prctl(ARCH_SET_FS, RDI+0x8)

		// store through it, to make sure it works
		get_tls(BX)
		MOVQ	$0x123, g(BX)
		MOVQ	runtime·m0+m_tls(SB), AX
		CMPQ	AX, $0x123
		JEQ 2(PC)
		CALL	runtime·abort(SB)
	ok:
		// set the per-goroutine and per-mach "registers"
		get_tls(BX)
		LEAQ	runtime·g0(SB), CX
		MOVQ	CX, g(BX)		//runtime.m0.tls[0] = g0
		LEAQ	runtime·m0(SB), AX

		// save m->g0 = g0
		MOVQ	CX, m_g0(AX)	//设定runtime.m0.g0 = runtime.g0
		// save m0 to g0->m
		MOVQ	AX, g_m(CX)		//设定runtime.g0.m = runtime.m0

		CLD				// convention is D is always left cleared
		CALL	runtime·check(SB)

		MOVL	16(SP), AX		// copy argc
		MOVL	AX, 0(SP)
		MOVQ	24(SP), AX		// copy argv
		MOVQ	AX, 8(SP)
		CALL	runtime·args(SB)	//设定runtime.argc/argv,从auxv中获得random data和pagesize，vdso信息初始化(vdso_linux.go:vdsoauxv)
		CALL	runtime·osinit(SB)	//获取cpu核心数(通过sched_getaffinity)以及huge page大小
		CALL	runtime·schedinit(SB)	//runtime.schedinit函数执行了大部分初始化操作

		// create a new goroutine to start program
		MOVQ	$runtime·mainPC(SB), AX		// entry
		PUSHQ	AX
		PUSHQ	$0			// arg size
		CALL	runtime·newproc(SB)
		POPQ	AX
		POPQ	AX

		// start this M
		CALL	runtime·mstart(SB)

		CALL	runtime·abort(SB)	// mstart should never return
		RET

		// Prevent dead-code elimination of debugCallV1, which is
		// intended to be called by debuggers.
		MOVQ	$runtime·debugCallV1(SB), AX
		RET


5. Go 1.14的非协作式抢占
	GODEBUG=asyncpreemptoff=1可以关闭
	运行超过10ms的goroutine会被标记为Preemptable，并且在恰当的时机(isAsyncSafePoint返回true)时由sysmon进程发出抢占信号SIGURG。
	调用链：

	sysmon
	↓
	retake
	↓
	preemptone
	↓
	preemptM
	↓
	signalM(mp, sigPreempt)：发送SIGURG信号给指定线程
	↓
	sigtramp(sig uint32, info *siginfo, ctx unsafe.Pointer)：直接信号处理函数
	↓
	sigtrampgo(sig uint32, info *siginfo, ctx unsafe.Pointer)
	↓
	sighandler(sig, info, ctx, g)
	↓
	doSigPreempt(signal handler)
	↓
	ctxt.pushCall(funcPC(asyncPreempt))
	```Go
	func (c *sigctxt) pushCall(targetPC uintptr) {
		// Make it look like the signaled instruction called target.
		pc := uintptr(c.rip())		//get 中断时的rip寄存器
		sp := uintptr(c.rsp())		//get 中断时的rsp寄存器
		sp -= sys.PtrSize			//sp -= 8
		*(*uintptr)(unsafe.Pointer(sp)) = pc	//压入返回地址，看起来像是直接调用asyncPreempt
		c.set_rsp(uint64(sp))		//设定新rsp
		c.set_rip(uint64(targetPC))	//设定下一步执行的指令为asyncPreempt
	}
	```
	↓
	asyncPreempt
	备份所有寄存器值到栈上，然后调用asyncPreempt2，返回后再从栈上恢复寄存器上下文。
	↓ 
	asyncPreempt2
	根据g.preemptStop标志决定调用mcall(preemptPark)还是mcall(gopreempt_m)，
	- preemptPark：做了3件事
		(1)调用dropg()将g和m分离,
		(2)将g的状态设定为_Gpreempted,
		(3)调用schedule()来为M找到新的可运行g
	- gopreempt_m：调用goschedImpl，具体为
		(1)将当前g从_Grunning改为_Grunnable,
		(2)调用dropg()将g和m分离,
		(3)调用globrunqput将g放入全局队列sched.runq


6. m初始化
	全局的runtime.allm指向了所有m串成的单向链表，m中的alllink指向了下一个m，
	每次新生成的m插入链表头，即 m5 -> m4 -> m3 -> m2 -> m1 -> runtime.m0


	mstart
	↓
	mstart1	/ needm
	↓
	minit
	↓
	minitSignals
	↓
	minitSignalStack / minitSignalMask


	mstartm0
	↓
	initsig
	↓
	setsig(i, funcPC(sighandler))
	↓
	sigaction

7. schedule时机

	- startm -> newm -> newm1 -> newosproc -> clone -> mstart -> mstart1	->	schedule
	  // mstart is the entry-point for new Ms

	- gopark	->	park_m	->	schedule
	  // gopark puts the current goroutine into a waiting state and calls unlockf. 
	  // most synchronization events will enter this.(channel, sleep, mutex etc.)

	- Gosched		->		gosched_m	-> 		|
												|-> goschedImpl	-> schedule
	  asyncPreempt2	->	|						|
	  					|	->	gopreempt_m ->	|
	  newstack		->	|

	- asyncPreempt2	->	|
						|->	preemptPark	-> schedule
	  newstack		->	|

	  // preemptPark parks gp and puts it in _Gpreempted.
	  // asyncPreempt2 is used for asynchronous preemption in Go 1.14 later.
	  // newstack is called by runtime.morestack

	- goyield	-> mcall(goyield_m)	-> schedule
	  // put into local runqueue
	
	- goexit1	->	mcall(goexit0)	->	schedule

	- exitsyscall0	->	schedule



8. mcall
	runtime.mcall分为以下几个步骤
		- 保存caller的goroutine上下文到其sched成员中
			+ MOV caller_PC g.sched.pc	//注意这里的pc是caller中mcall的下一个指令地址
			+ MOV caller_SP g.sched.sp
			+ MOV caller_g g.sched.g
			+ MOV caller_BP g.sched.bp
		- 切换到当前m的g0栈，执行传入的fn函数，具体：
			+ MOV g0.sp SP	//栈帧寄存器指向g0栈的栈帧位置, SP = m->g0.sched.sp
			+ PUSH g		//将caller的g入栈
			+ call fn		//执行传入的fn	//fn理论上必须是永不返回,比如schedule函数这种
			+ POP 			//理论上不会执行这步
			+ ...

9. runtime.newproc
本节着重关注协程创建时参数的传递策略。
```go
func newproc(siz int32, fn *funcval)
```

压栈时如下：
	rsp -> siz	// 4字节，MOVL指令
	rsp+0x8	->	fn	// 8字节
	rsp+0x10,0x18,.. -> 闭包变量





