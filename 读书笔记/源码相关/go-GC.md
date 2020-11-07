阅读摘要from：https://juejin.im/post/5d78b3276fb9a06b1829e691

1. go1.5引入三色标记，go1.8引入写屏障。

2. 三色标记
    - 黑色集合中确保没有任何指针指向白色集合
    - 白色集合中的对象允许有指针指向黑色集合    ->> 待GC的对象
    - 灰色集合可能会有指针指向白色集合里的对象

    GC开始时，
        - 所有对象标记为白色。
        - 遍历所有根对象并标记为灰色。
        - 遍历灰色对象，将引用到的白色对象标记为灰色，同时把遍历过的灰色对象标记为黑色。
        - 重复上述第三个步骤，直到没有对象为灰色
        - 此时白色对象就是不可达对象，可以GC掉
    
    写屏障(write barrier)
        遍历灰色对象时灰色对象可能会改变，导致某些对象扫描不到被误认为白色垃圾而被GC。
        如: A(灰) -> B,C
            D(灰) -> E
            扫描完A后将A标记为黑，B、C标记为灰，扫描到D时有goroutine将D->E改成A->E,
            则E虽然还在被引用但会GC掉。
        
        go1.7以前，采用的是Dijkstra-style insertion write barrier，将被指向的对象变灰，
		这样新的对象创建或者黑色对象指向白色的时候，目标会变灰从而满足黑色不会指向白色。

		```
		writePointer(slot, ptr):
    		shade(ptr)	//将ptr标记为灰色
    		*slot = ptr //slot指向ptr
		```

		该写屏障只能shade堆上的内存引用，栈上实现写屏障有较大难度，因此需要STW来rescan一次。


		为了缩小STW时间，go1.8开始引入混合写屏障来去掉rescan，
		```
		writePointer(slot, ptr):
			shade(*slot)
			if current stack is grey:
				shade(ptr)
			*slot = ptr
		```



3. GC流程
    - Sweep Termination: 对未清扫的span进行清扫, 只有上一轮的GC的清扫工作完成才可以开始新一轮的GC
    - Mark: 扫描所有根对象, 和根对象可以到达的所有对象, 标记它们不被回收
    - Mark Termination: 完成标记工作, 重新扫描部分根对象(要求STW)
    - Sweep: 按标记结果清扫span

	0. 开始时所有对象都被标为白色(可被回收)
	1. STW，将所有栈及global变量标为灰色 (Stack scan)
	2. 从灰色对象开始并发标记所有可达对象，将可达对象标为灰色，当一个对象所有引用对象都被标为灰色后，该对象就被置为黑色。(Mark)
	3. 当所有对象都被标为黑色后，再进行一次STW，重新扫描栈和global未扫描过的对象。(Mark termination)
	4. 并发的清除所有白色对象。(Sweep)
	参考assets/GO-GC-Algorithm-Phases.png


4. 部分关键源码

三色标记    //mgcmark.go::gcDrain()
```Go
func gcDrain(gcw *gcWork, flags gcDrainFlags) {
	if !writeBarrier.needed {
		throw("gcDrain phase incorrect")
	}

	gp := getg().m.curg         //获取当前g
	preemptible := flags&gcDrainUntilPreempt != 0   //是否可抢占,恒为false
	flushBgCredit := flags&gcDrainFlushBgCredit != 0
	idle := flags&gcDrainIdle != 0

	initScanWork := gcw.scanWork

	// checkWork is the scan work before performing the next
	// self-preempt check.
	checkWork := int64(1<<63 - 1)
	var check func() bool
	if flags&(gcDrainIdle|gcDrainFractional) != 0 {
		checkWork = initScanWork + drainCheckThreshold
		if idle {
			check = pollWork
		} else if flags&gcDrainFractional != 0 {
			check = pollFractionalWorkerExit
		}
	}

	// Drain root marking jobs.
	if work.markrootNext < work.markrootJobs {
		for !(preemptible && gp.preempt) {
			job := atomic.Xadd(&work.markrootNext, +1) - 1
			if job >= work.markrootJobs {
				break
			}
			markroot(gcw, job)
			if check != nil && check() {
				goto done
			}
		}
	}

	// Drain heap marking jobs.
	for !(preemptible && gp.preempt) {
		// Try to keep work available on the global queue. We used to
		// check if there were waiting workers, but it's better to
		// just keep work available than to make workers wait. In the
		// worst case, we'll do O(log(_WorkbufSize)) unnecessary
		// balances.
		if work.full == 0 {
			gcw.balance()
		}

		b := gcw.tryGetFast()
		if b == 0 {
			b = gcw.tryGet()
			if b == 0 {
				// Flush the write barrier
				// buffer; this may create
				// more work.
				wbBufFlush(nil, 0)
				b = gcw.tryGet()
			}
		}
		if b == 0 {
			// Unable to get work.
			break
		}
		scanobject(b, gcw)

		// Flush background scan work credit to the global
		// account if we've accumulated enough locally so
		// mutator assists can draw on it.
		if gcw.scanWork >= gcCreditSlack {
			atomic.Xaddint64(&gcController.scanWork, gcw.scanWork)
			if flushBgCredit {
				gcFlushBgCredit(gcw.scanWork - initScanWork)
				initScanWork = 0
			}
			checkWork -= gcw.scanWork
			gcw.scanWork = 0

			if checkWork <= 0 {
				checkWork += drainCheckThreshold
				if check != nil && check() {
					break
				}
			}
		}
	}

done:
	// Flush remaining scan work credit.
	if gcw.scanWork > 0 {
		atomic.Xaddint64(&gcController.scanWork, gcw.scanWork)
		if flushBgCredit {
			gcFlushBgCredit(gcw.scanWork - initScanWork)
		}
		gcw.scanWork = 0
	}
}
```




