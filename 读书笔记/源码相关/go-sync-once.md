
sync.Once

# 0. basis
	refer：
		https://colobu.com/2021/05/05/triple-gates-of-sync-Once/
		https://groups.google.com/g/golang-nuts/c/XTltSRQccmY/m/1fFi9mnCSi0J
		https://codereview.appspot.com/4641066

		https://github.com/golang/go/issues/5045		// atomic vs memory model

# 1. source 

```Go
type Once struct {
	done uint32     //done代表f是否完成，1->完成，0->未完成
	m    Mutex      //避免多个并行的调用导致f执行多次
}
func (o *Once) Do(f func()) {
	// Note: Here is an incorrect implementation of Do:
	//
	//	if atomic.CompareAndSwapUint32(&o.done, 0, 1) {
	//		f()
	//	}
	//
	// Do guarantees that when it returns, f has finished.
	// This implementation would not implement that guarantee:
	// given two simultaneous calls, the winner of the cas would
	// call f, and the second would return immediately, without
	// waiting for the first's call to f to complete.
	// This is why the slow path falls back to a mutex, and why
	// the atomic.StoreUint32 must be delayed until after f returns.

	if atomic.LoadUint32(&o.done) == 0 {
		// Outlined slow-path to allow inlining of the fast-path.
		o.doSlow(f)
	}
}
func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}
```

# 2. 分析

```Go
func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 {			// 如注释所说，保证sync.Once返回时f一定执行完毕，用CAS方式竞争失败的一方会直接返回而非等待竞争成功者执行完返回
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {		// o.done在o.m的保护下因此不需要atomic保证内存序
		defer atomic.StoreUint32(&o.done, 1)		// 用defer保证即使f函数panic也会设置o.done, 对o.done用atomic则是保证store操作必定会在f返回后进行的memory order
		f()
	}
}


```


