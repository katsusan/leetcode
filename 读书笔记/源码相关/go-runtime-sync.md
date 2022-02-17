
# 0. 引言
golang里有诸多同步原语的实现，比如channel通信时用的是runtime内部的mutex，   
而标准库里也有sync.Mutex、sync.RWMutex等各种形式的锁，下面就其用途与实现   
稍作挖掘来深入了解runtime中的同步原语背后的设计用意。

# 1 runtime.mutex

// From runtime/HACKING.md   
The simplest is `mutex`, which is manipulated using `lock` and
`unlock`. This should be used to protect shared structures for short
periods. Blocking on a `mutex` directly blocks the M, without
interacting with the Go scheduler. This means it is safe to use from
the lowest levels of the runtime, but also prevents any associated G
and P from being rescheduled. `rwmutex` is similar.

# 1.1 定义：
```Go
// Mutual exclusion locks.  In the uncontended case,
// as fast as spin locks (just a few user-level instructions),
// but on the contention path they sleep in the kernel.
// A zeroed Mutex is unlocked (no need to initialize each lock).
// Initialization is helpful for static lock ranking, but not required.
type mutex struct {
	// Empty struct if lock ranking is disabled, otherwise includes the lock rank
	lockRankStruct
	// Futex-based impl treats it as uint32 key,
	// while sema-based impl as M* waitm.
	// Used to be a union, but unions break precise GC.
	key uintptr
}

const (
	mutex_unlocked = 0
	mutex_locked   = 1
	mutex_sleeping = 2

	active_spin     = 4
	active_spin_cnt = 30
	passive_spin    = 1
)
```

# 1.2 用途：
// 参考lockrank.go里列出的lockName，大概60+处使用了runtime.mutex。
- hchan.lock：保护channel内数据成员
- cpuProfile.lock: 保护全局变量`var cpuprof cpuProfile`
- schedt.lock: 保护全局调度结构`var sched schedt`
- itabLock: 保护itab表，`var itabLock mutex`
- globalAlloc.mutex: 保护全局内存结构globalAlloc的persistentAlloc(用途为function/type/debug-related persistent data，参考persistentalloc函数)
- debugPtrmask.mutex: (disabled)
- mheap.lock: 保护全局堆结构`var mheap_ mheap`
- mcentral.lock: 保护mcentral(比如mcache里没有空闲mspan需要从mcentral里寻找时)
- finlock: 保护finalizer的全局统计数据，见mfinal.go
- work.wbufSpans.lock:
- work.assistQueue.lock: 保护assistQueue
- work.sweepWaiters.lock: 保护sweepwaiters(mark termination转向sweep时需唤醒的被阻塞协程)   
- ...
 

# 1.3 内部实现
# 1.3.1 call chain:   
lock(l *mutex)   
↓   
lockWithRank(l, getLockRank(l))   
↓   
lock2(l)

unlock(l *mutex)
↓
unlockWithRank(l)
↓
unlock2(l)

# 1.3.2 lock2
- gp.m.locks++      // 锁定当前线程m
- // Speculative grab for lock
  - v := atomic.Xchg(key32(&l.key), mutex_locked) // 将l.key与mutex_locked(1)交换，并返回原本的l.key值写到v里
  - 若v为mutex_unlocked则代表原来的l.key为unlocked(0)即未上锁状态，则xchgl执行后l.key变为mutex_locked(1),代表上锁成功，可直接返回
    - return
- // 否则进入循环请求锁阶段，此时l.key要不就是mutex_locked(1)-锁被别人持有，要不就是mutex_sleeping(2)-有休眠线程在等待锁
- wait := v   // 保存原本的l.key值，原因在于若l.key为mutex_sleeping时，为了能让休眠中的进程被唤醒后面必须还原l.key为mutex_sleeping
- spin := 0   // 单核cpu不自旋 
- if ncpu > 1
  - spin = active_spin  // 多核情况下自旋active_spin(4)次
- for {   // 加锁循环
  - // Try for lock, spinning.
  - for i := 0; i < spin; i++ {   // 自旋spin(4)次
    - for l.key == mutex_unlocked {
      - if atomic.Cas(key32(&l.key), mutex_unlocked, wait)   // 比较l.key与mutex_unlocked若相等代表其它线程释放了锁，置为wait上锁
        - return  // Cas返回true代表把l.key从mutex_unlocked变为wait成功，即加锁成功因此直接返回
      - procyield(active_spin_cnt)  // 自旋active_spin_cnt(30)次，在x86下procyield底层用PAUSE指令实现，具体就是循环执行30次PAUSE指令

  - // Try for lock, rescheduling.
  - for i := 0; i < passive_spin; i++ { // 让出CPU执行权passive_spin(1)次
    - for l.key == mutex_unlocked {
      - if atomic.Cas(key32(&l.key), mutex_unlocked, wait) { // 逻辑同上
        - return
      - osyield() // 不同的是争取锁失败的话这里会调用osyield()让出执行权，runtime.osyield底层使用了系统调用sched_yield
  
  - // Sleep.
  - v = atomic.Xchg(key32(&l.key), mutex_sleeping)    // 经过4次自旋和1次主动yield后仍未获得锁的情况下，将l.key置为mutex_sleeping，表明自己将进入睡眠
  - if v == mutex_unlocked {    // 若xchg前l.key为unlocked则代表可以加锁
    - return
  - wait = mutex_sleeping   // ★ 这里令wait为sleeping可以保证后续sleeping的线程能被unlock时futexwakeup唤醒
  - futexsleep(key32(&l.key), mutex_sleeping, -1)


简而言之，lock2函数里改变l.key有如下几处：
- v := atomic.Xchg(key32(&l.key), mutex_locked)
  * 1) 若l.key = unlocked, √
  * 2) 若l.key = locked, 进入forloop
  * 3) 若l.key = sleeping, 进入forloop
- for 
  + atomic.Cas(key32(&l.key), mutex_unlocked, wait)  // wait:=v, 为locked或sleeping状态
  + v = atomic.Xchg(key32(&l.key), mutex_sleeping)  // 改变l.key为sleeping，然后自己futexsleep

example：
  T1.lock(), T2.lock(), T3.lock(), T1.unlock()
    - T1.lock: l.key → locked
    - T2.lock: l.key → sleeping, futexsleep(&l.key, sleeping, -1)   // wait = sleeping
    - T3.lock: l.key → sleeping，futexsleep(&l.key, sleeping, -1)   // wait = sleeping
    - T1.unlock: l.key → unlocked，futexwakeup(&l.key, 1) 
      // 唤醒T2或T3中的一个，T2/T3被唤醒后通过forloop里的CAS加锁成功
      // CAS后把l.key变为wait(sleeping)，因此T2被唤醒后lock然后unlock时可以继续futexwakeup唤醒T3



# 1.3.3 unlock2
- v := atomic.Xchg(key32(&l.key), mutex_unlocked)
- if v == mutex_unlocked  → throw("unlock of unlocked lock")  // 若l.key为unlocked状态代表锁之前并未被lock，抛出异常
- if v == mutex_sleeping  → futexwakeup(key32(&l.key), 1)   // 若l.key为sleeping状态，代表有线程futexsleep在当前锁上，需要futexwakeup唤醒
- gp.m.locks--    // 解除与当前m的绑定

