# 0


# 1. definition

```Golang
const (
	mutexLocked = 1 << iota // mutex is locked
	mutexWoken
	mutexStarving
	mutexWaiterShift = iota
	starvationThresholdNs = 1e6
)

// A Mutex is a mutual exclusion lock.
// The zero value for a Mutex is an unlocked mutex.
//
// A Mutex must not be copied after first use.
type Mutex struct {
	state int32
	sema  uint32
}
```
Mutex.state = 0								→		未上锁(初始化状态)
Mutex.state = mutexLocked(1)	→		已上锁
Mutex.state = mutexWoken(2)		→		
Mutex.state = mutexStarving(4)	→			

Mutex有两种模式：normal和starvation。

normal模式下所有waiter在FIFO的队列等待。新加入争抢锁的协程有一个优势，
就是它们已经获得了CPU的执行权，并且这样的争抢者可能有很多。因此唤醒的waiter
极有可能再次再次争抢失败，这种情况下它会被插入到等待队列头部。
如果一个waiter等待锁的时间超过1e6纳秒(1ms)，那么它会把锁置为starvation模式。

starvation模式下锁的持有权被直接从unlock的协程移交到队列头的协程，
新加入争抢锁的协程不会尝试获得锁，即使锁马上可能会被释放，同时它也不会自旋，
而是直接插入到等待队列尾部。

如果一个waiter获得了锁但发现(1)它在等待队列尾部或者(2)它的等待时间<1ms，
则会把锁切换回normal模式。

# 2. lock flow 

- Fastpath: if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked)	→	return
- slowpath: `m.lockSlow()`
```Go
func (m *Mutex) lockSlow() {
	var waitStartTime int64
	starving := false
	awoke := false
	iter := 0
	old := m.state
	for {
		// Don't spin in starvation mode, ownership is handed off to waiters
		// so we won't be able to acquire the mutex anyway.
		if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
			// old为末比特位为0?1，也就是 locked & not starving，且此时可以自旋
			// 当iter>=4时，runtime_canSpin会返回false，也就是runtime_doSpin自旋4次还没获得锁会进入最差情况即当前协程陷入休眠变成waiter
			// Active spinning makes sense.
			// Try to set mutexWoken flag to inform Unlock
			// to not wake other blocked goroutines.
			if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
				atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
				// awoke = false && mutexWoken位为0 
				// && old的高29位不为0 → 有waiter存在
				// && 把m.state从old设为old|mutexWoken成功	→ 没有其它竞争者修改到锁的状态
				// 说明自旋过程中获得了锁，设置awoke=true，此时m.state中mutexWoken位也被置为1
				awoke = true
			}
			runtime_doSpin()		// 自旋30次，procyield(active_spin_cnt=30)
			iter++
			old = m.state
			continue
		}
		new := old
		// Don't try to acquire starving mutex, new arriving goroutines must queue.
		// 如果是非starving状态，则为new设置locked标志
		// 而starving状态下由于没设置locked，故会进入之后的runtime_SemacquireMutex逻辑等待唤醒
		if old&mutexStarving == 0 {
			new |= mutexLocked
		}
		if old&(mutexLocked|mutexStarving) != 0 {		// old处于locked或starving状态，则waiter++
			new += 1 << mutexWaiterShift
		}
		// The current goroutine switches mutex to starvation mode.
		// But if the mutex is currently unlocked, don't do the switch.
		// Unlock expects that starving mutex has waiters, which will not
		// be true in this case.
		if starving && old&mutexLocked != 0 {		// old处于locked状态且starving为true,把锁切换为starving模式
			new |= mutexStarving
		}
		if awoke {
			// The goroutine has been woken from sleep,
			// so we need to reset the flag in either case.
			if new&mutexWoken == 0 {
				throw("sync: inconsistent mutex state")
			}
			new &^= mutexWoken			// 将mutexWoken位设为0，&^ = bit clear = and not
		}
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			if old&(mutexLocked|mutexStarving) == 0 {			// 如果此时m.state为unlocked且非starving状态,被CAS成功了则代表获得了锁
				break // locked the mutex with CAS
			}
			// If we were already waiting before, queue at the front of the queue.
			queueLifo := waitStartTime != 0
			if waitStartTime == 0 {
				waitStartTime = runtime_nanotime()
			}
			runtime_SemacquireMutex(&m.sema, queueLifo, 1)	// 对于已上锁的m，竞争者会阻塞在这里，等待锁被释放后唤醒
			starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
			old = m.state
			if old&mutexStarving != 0 {
				// If this goroutine was woken and mutex is in starvation mode,
				// ownership was handed off to us but mutex is in somewhat
				// inconsistent state: mutexLocked is not set and we are still
				// accounted as waiter. Fix that.
				if old&(mutexLocked|mutexWoken) != 0 || old>>mutexWaiterShift == 0 {
					throw("sync: inconsistent mutex state")
				}
				delta := int32(mutexLocked - 1<<mutexWaiterShift)		// delta = 1 - 1<<3
				if !starving || old>>mutexWaiterShift == 1 {
					// Exit starvation mode.
					// Critical to do it here and consider wait time.
					// Starvation mode is so inefficient, that two goroutines
					// can go lock-step infinitely once they switch mutex
					// to starvation mode.
					delta -= mutexStarving
				}
				atomic.AddInt32(&m.state, delta)
				break
			}
			awoke = true
			iter = 0
		} else {
			old = m.state
		}
	}
}
```


# 3 unlock flow

- fastpath: new := atomic.AddInt32(&m.state, -mutexLocked); new != 0 → slowpath，else unlock succeeds.
- slowpath: 核心是调用runtime_Semrelease唤醒等待列表上的g
```Go
func (m *Mutex) unlockSlow(new int32) {
	if (new+mutexLocked)&mutexLocked == 0 {			// new+mutexLocked = m.state, state&mutexLocked应该为1代表处于locked状态
		throw("sync: unlock of unlocked mutex")
	}
	if new&mutexStarving == 0 {
		// 不在starving状态，那么寻找waiter唤醒
		old := new
		for {
			// If there are no waiters or a goroutine has already
			// been woken or grabbed the lock, no need to wake anyone.
			// In starvation mode ownership is directly handed off from unlocking
			// goroutine to the next waiter. We are not part of this chain,
			// since we did not observe mutexStarving when we unlocked the mutex above.
			// So get off the way.
			if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
				// 没有waiter，或者有goroutine已经被唤醒/获得了锁
				return
			}
			// Grab the right to wake someone.
			new = (old - 1<<mutexWaiterShift) | mutexWoken
			if atomic.CompareAndSwapInt32(&m.state, old, new) {
				runtime_Semrelease(&m.sema, false, 1)
				return
			}
			old = m.state
		}
	} else {
		// Starving mode: handoff mutex ownership to the next waiter, and yield
		// our time slice so that the next waiter can start to run immediately.
		// Note: mutexLocked is not set, the waiter will set it after wakeup.
		// But mutex is still considered locked if mutexStarving is set,
		// so new coming goroutines won't acquire it.
		runtime_Semrelease(&m.sema, true, 1)
	}
}
```
