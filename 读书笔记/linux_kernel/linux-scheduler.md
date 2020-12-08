refer:  
	https://www.cs.montana.edu/~chandrima.sarkar/AdvancedOS/CSCI560_Proj_main/index.html
	http://www.cs.montana.edu/~chandrima.sarkar/AdvancedOS/SchedulingLinux/index.html

# 1. major features
	- The 2.6 scheduler was designed and implemented by Ingo Molnar. His motivation in working on the new scheduler was 
		to create a completely O(1) scheduler for wakeup, context-switch, and timer interrupt overhead

	- One of the issues that triggered the need for a new scheduler was the use of Java virtual machines (JVMs). 
		The Java programming model uses many threads of execution, which results in lots of overhead for scheduling in an O(n) scheduler

	- Each CPU has a runqueue made up of 140 priority lists that are serviced in FIFO order. Tasks that are scheduled to 
		execute are added to the end of their respective runqueue's priority list
		// 每个CPU都有140个FIFO的优先级任务列表，被调度任务会被加到对应执行队列的末尾。

	- Each task has a time slice that determines how much time it's permitted to execute

	- The first 100 priority lists of the runqueue are reserved for real-time tasks, and the last 40 are used for 
		user tasks (MAX_RT_PRIO=100 and MAX_PRIO=140)
		// 前100的优先级执行队列是为实时任务保留，之后的40个用于用户任务。

	- In addition to the CPU's runqueue, which is called the active runqueue, there's also an expired runqueue
		// 除了CPU的执行队列(也被称为活动队列active queue)之外，相对应的还有过期队列expired queue。

	- When a task on the active runqueue uses all of its time slice, it's moved to the expired runqueue. 
		During the move, its time slice is recalculated (and so is its priority)
		// 当active队列的任务用完所有的时间片后它会被移到expired队列。
		// 移动时会重新计算它的时间片以及优先级、

	- If no tasks exist on the active runqueue for a given priority, the pointers for the active and expired runqueues are swapped,
		thus making the expired priority list the active one
		// 如果某个优先级的active队列里没有任务了，则会交换active队列和expired队列的指针，使得expired队列的任务能重新投入运行。


# 2. kernel 2.6 scheduler policies
	- SCHED_NORMAL: A conventional, time-shared process		//通常任务的类型，共享时间片(time-shared)，以前也叫SCHED_OTHER
		- Each task assigned a “Nice” value
		- PRIO = MAX_RT_PRIO + NICE + 20
		- Assigned a time slice
		- Tasks at the same prio(rity) are round-robined
		- Ensures Priority + Fairness

	- SCHED_FIFO：A First-In, First-Out real-time process		// 先进先出FIFO的实时调度类型
		- Run until they relinquish the CPU voluntarily		// 一直执行到自己主动放弃CPU
		- Priority levels maintained
		- Not pre-empted !!			// 不会被抢占

	- SCHED_RR：A Round Robin real-time process	// 轮训+实时型调度策略
		- Assigned a timeslice and run till the timeslice is exhausted.		// 分配其时间片并且一直运行到时间片用完
		- Once all RR tasks of a given prio(rity) level exhaust their timeslices, their timeslices are refilled and they continue running
		- Prio(rity) levels are maintained

	- SCHED_BATCH：for "batch" style execution of processes
		- For computing-intensive tasks			// 适用于计算密集型任务
		- Timeslices are long and processes are round robin scheduled
		- lowest priority tasks are batch-processed (nice +19)	// 极低优先级的任务是batch-processed

	- SCHED_IDLE：for running very low priority background job
		- nice value has no influence for this policy		// 调整nice值对这种任务无影响
		- extremely low priority (lower than +19 nice)	// 极其极其低的优先级(低于+19的nice值)


# 2. CFS

	The main idea behind the CFS is to maintain balance (fairness) in providing processor time to tasks. This means processes should be given a fair amount of the processor. When the time for tasks is out of balance (meaning that one or more tasks are not given a fair amount of time relative to others), then those out-of-balance tasks should be given time to execute.

	To determine the balance, the CFS maintains the amount of time provided to a given task in what's called the virtual runtime. The smaller a task's virtual runtime—meaning the smaller amount of time a task has been permitted access to the processor—the higher its need for the processor. The CFS also includes the concept of sleeper fairness to ensure that tasks that are not currently runnable (for example, waiting for I/O) receive a comparable share of the processor when they eventually need it.

	But rather than maintain the tasks in a run queue, as has been done in prior Linux schedulers, the CFS maintains a time-ordered red-black tree (see Figure below). A red-black tree is a tree with a couple of interesting and useful properties. First, it's self-balancing, which means that no path in the tree will ever be more than twice as long as any other. Second, operations on the tree occur in O(log n) time (where n is the number of nodes in the tree). This means that you can insert or delete a task quickly and efficiently.

# 3. Components
	- A Running Queue : A running queue (rq) is created for each processor (CPU). It is defined in kernel/sched.h as `struct rq`. 
		Each rq contains a list of runnable processes on a given processor. The struct rq is defined in sched.c notsched.h to 
		abstract the internal data structure of the scheduler.

	- Schedule Class : schedule class was introduced in 2.6.23. It is an extensible hierarchy of scheduler modules. 
		These modules encapsulate scheduling policy details and are called from the scheduler core without the core code assuming too much about them. 
		Scheduling classes are implemented through the sched_class structure, which contains hooks to functions that must be called 
		whenever an interesting event occurs. Tasks refer to their schedule policy through struct task_struct and sched_class.

		There are two schedule classes implemented in 2.6.32:
			- Completely Fair Schedule class: schedules tasks following Completely Fair Scheduler (CFS) algorithm. 
				Tasks which have policy set to SCHED_ NORMA L (SCHED_OTHER), SCHED_BATCH, SCHED_IDLE are scheduled by this schedule class. 
				The implementation of this class is in kernel /sched_fai r.c
			- RT schedule class: schedules tasks following real-time mechanism defined in POSIX standard. 
				Tasks which have policy set to SCHED_FIFO, SCHED_RR are scheduled using this schedule class. 
				The implementation of this class is kernel/sched_rt.c

		- Load balancer: In SMP environment, each CPU has its own rq. These queues might be unbalanced from time to time. 
			A running queue with empty task pushes its associated CPU to idle, which does not take full advantage of symmetric multiprocessor systems. 
			Load balancer is to address this issue. It is called every time the system requires scheduling tasks. 
			If running queues are unbalanced, load balancer will try to pull idle tasks from busiest processors to idle processor.


# 4. core implementation
schedule()是scheduler的对外暴露的调度API，代码如下：

```C++
asmlinkage __visible void __sched schedule(void)
{
	struct task_struct *tsk = current;

	sched_submit_work(tsk);
	do {
		preempt_disable();
		__schedule(false);
		sched_preempt_enable_no_resched();
	} while (need_resched());
}
EXPORT_SYMBOL(schedule);

struct thread_info {
	unsigned long		flags;		/* low level flags */
	u32			status;		/* thread synchronous flags */
};
```

可以看出当need_resched()返回true也就是线程的flags里TIF_NEED_RESCHED为1时进行调度，
核心是_schedule()函数，该函数的注释写明了进入schedule函数的时机有以下几种;
	- Explicit blocking: mutex, semaphore, waitqueue, etc.
	- TIF_NEED_RESCHED flag is checked on interrupt and userspace return
		paths. For example, see arch/x86/entry_64.S.
	- Wakeups don't really cause entry into schedule().
		They add a task to the run-queue and that's it.
当加入到runqueue的新task抢占了当前task时，wakeup会设置TIF_NEED_RESCHED标志位然后
schedule函数会在下述可能的场景下被调用：
	- 如果内核被设定为可抢占的，即编译选项CONFIG_PREEMPT=y：
		- 在系统调用syscall以及异常exception上下文中，下一个最外层的preempt_enable()返回时，
			很可能在wake_up()的spin_unlock()执行后。
		- 在IRQ上下文中，从中断处理程序interrupt handler返回到可抢占的上下文
	- 如果内核未被设定为可抢占，即CONFIG_PREEMPT未被设置，那么在下一次：
		- cond_resched() call					//cond_resched()调用
		- explicit schedule() call		//显式的schedule()调用
		- return from syscall or exception to user-space	//从系统调用/异常返回到用户态
		- return from interrupt-handler to user-space			//从中断处理程序返回到用户态

```C++
static void __sched notrace __schedule(bool preempt)
{
	struct task_struct *prev, *next;
	unsigned long *switch_count;
	struct rq_flags rf;
	struct rq *rq;
	int cpu;

	cpu = smp_processor_id();
	rq = cpu_rq(cpu);		// 获取当前cpu的runqueue 
	prev = rq->curr;		// 保存调度前的任务于prev

	schedule_debug(prev);

	if (sched_feat(HRTICK))
		hrtick_clear(rq);

	local_irq_disable();
	rcu_note_context_switch(preempt);

	/*
	 * Make sure that signal_pending_state()->signal_pending() below
	 * can't be reordered with __set_current_state(TASK_INTERRUPTIBLE)
	 * done by the caller to avoid the race with signal_wake_up().
	 *
	 * The membarrier system call requires a full memory barrier
	 * after coming from user-space, before storing to rq->curr.
	 */
	rq_lock(rq, &rf);
	smp_mb__after_spinlock();

	/* Promote REQ to ACT */
	rq->clock_update_flags <<= 1;
	update_rq_clock(rq);

	switch_count = &prev->nivcsw;
	if (!preempt && prev->state) {
		if (signal_pending_state(prev->state, prev)) {			// 如果prev任务有信号待处理则继续设置其为TASK_RUNNING状态
			prev->state = TASK_RUNNING;
		} else {
			deactivate_task(rq, prev, DEQUEUE_SLEEP | DEQUEUE_NOCLOCK);		// 否则将prev任务从runqueue中移出
			prev->on_rq = 0;

			if (prev->in_iowait) {
				atomic_inc(&rq->nr_iowait);
				delayacct_blkio_start();
			}

			/*
			 * If a worker went to sleep, notify and ask workqueue
			 * whether it wants to wake up a task to maintain
			 * concurrency.
			 */
			if (prev->flags & PF_WQ_WORKER) {
				struct task_struct *to_wakeup;

				to_wakeup = wq_worker_sleeping(prev);
				if (to_wakeup)
					try_to_wake_up_local(to_wakeup, &rf);
			}
		}
		switch_count = &prev->nvcsw;
	}

	next = pick_next_task(rq, prev, &rf);		// 找到最高优先级的下一个任务
	clear_tsk_need_resched(prev);
	clear_preempt_need_resched();

	if (likely(prev != next)) {
		rq->nr_switches++;
		rq->curr = next;
		/*
		 * The membarrier system call requires each architecture
		 * to have a full memory barrier after updating
		 * rq->curr, before returning to user-space.
		 *
		 * Here are the schemes providing that barrier on the
		 * various architectures:
		 * - mm ? switch_mm() : mmdrop() for x86, s390, sparc, PowerPC.
		 *   switch_mm() rely on membarrier_arch_switch_mm() on PowerPC.
		 * - finish_lock_switch() for weakly-ordered
		 *   architectures where spin_unlock is a full barrier,
		 * - switch_to() for arm64 (weakly-ordered, spin_unlock
		 *   is a RELEASE barrier),
		 */
		++*switch_count;

		trace_sched_switch(preempt, prev, next);

		/* Also unlocks the rq: */
		rq = context_switch(rq, prev, next, &rf);			// 切换任务上下文，包括mm和寄存器上下文
	} else {
		rq->clock_update_flags &= ~(RQCF_ACT_SKIP|RQCF_REQ_SKIP);
		rq_unlock_irq(rq, &rf);
	}

	balance_callback(rq);
}
```

可以看出真正关键的逻辑在于以下三行：
	- deactivate_task(rq, prev, DEQUEUE_SLEEP | DEQUEUE_NOCLOCK);
	- next = pick_next_task(rq, prev, &rf);
	- rq = context_switch(rq, prev, next, &rf);	


//TODO: pick_next_task和context_switch的实现



