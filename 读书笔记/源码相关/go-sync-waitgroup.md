
1. 结构定义
```Go
type WaitGroup struct {
	noCopy noCopy   //struct{}

	// 64-bit value: high 32 bits are counter, low 32 bits are waiter count.
	// 64-bit atomic operations require 64-bit alignment, but 32-bit
	// compilers do not ensure it. So we allocate 12 bytes and then use
	// the aligned 8 bytes in them as state, and the other 4 as storage
    // for the sema.
    //为了保证存储counter和waiter count的8字节对齐，因此多分配4字节
	state1 [3]uint32    
}
```

2. 核心逻辑
    ```Go
    func (wg *WaitGroup) state() (statep *uint64, semap *uint32) {
        //保证对齐的8字节存储于statep中，其中高4字节为counter，低4字节为waiter count
        if uintptr(unsafe.Pointer(&wg.state1))%8 == 0 {
            return (*uint64)(unsafe.Pointer(&wg.state1)), &wg.state1[2]
        } else {
            return (*uint64)(unsafe.Pointer(&wg.state1[1])), &wg.state1[0]
        }
    }

    func (wg *WaitGroup) Add(delta int) {
        statep, semap := wg.state()
        state := atomic.AddUint64(statep, uint64(delta)<<32)    //加到状态变量的高32位上

        // - Adds must not happen concurrently with Wait,   -> Add与Wait不可以并发执行
    	// - Wait does not increment waiters if it sees counter == 0. -> 如果counter为0，Wait并不会增加waiter count
    } 

    func (wg *WaitGroup) Done() {
        wg.Add(-1)
    }

    func (wg *WaitGroup) Wait() {
        statep, semap := wg.state()
        for {
            state := atomic.LoadUint64(statep)
            v := int32(state >> 32)     //取高32位v，v == 0则return
            w := uint32(state)
            ...
        }

    ```






