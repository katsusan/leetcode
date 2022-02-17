# 0. some useful links:
    https://0xax.gitbooks.io/linux-insides
    https://linux-kernel-labs.github.io/
		

# 1. interrupts and exceptions
	- interrupt: An interrupt is an asynchronous event that is typically triggered by an I/O device
		// 中断为异步事件，通常由IO设备产生
	- exception：a synchronous event that is generated when the processor detects one or more predefined conditions while executing an instruction 
		// 异常为同步事件，在处理器检测到预先定义好的情况时产生
		 The IA-32 architecture specifies three classes of exceptions: faults, traps, and aborts.
		// IA-32架构定义了三种异常：faults, traps, aborts。

异常的来源有两种:
	- processor detected: faults/traps/aborts，当处理器执行指令时所检测出的异常
		- fault: 通常可以在指令执行前就检测到并且可以被纠正，保存的EIP是产生fault时的指令地址，因此fault纠正后可以从该指令继续执行，比如page fault。
		- trap：在指令执行后检测到，保存的EIP是产生trap时的下一条指令，比如debug trap。
	- programme：int n

# 1.1 硬件概念：
	Programmable Interrupt Controller(PIC)

																			NMI		->		|-----------|
	device0	->	IRO0	-> 	|-------|									|		CPU			|
	device1	->	IRQ1	->	|  PIC	|	->	INTR	-> 		|-----------|
	device2	-> 	IRQ2	->	|-------|	
					
当连接到PIC的硬件有事件需要CPU处理时，会发生以下流程：
	- device raises an interrupt on the corresponding IRQn pin
	- PIC converts the IRQ into a vector number and writes it to a port for CPU to read
	- PIC raises an interrupt on CPU INTR pin
	- PIC waits for CPU to acknowledge an interrupt
	- CPU handles the interrupt

可以看出PIC这种设计有一个好处就是当CPU处理完当前中断例程前PIC不会发出另一个中断，也就是中断会串行执行。

	Advanced Programmable Interrupt Controller(APIC)
	
	External Interrupts	->	I/O APIC	->	Interrupt Controller Communication Bus	-> CPU0.local-APIC, CPU1.local-APIC,...

	I/O APIC负责分发外部中断到具体CPU核心。
	以上是中断的硬件相关概念，下面来讨论软件上如何处理中断。

# 1.2 中断控制Interrupt Control
IDT(Interrupt Descriptor Table)描述了中断/异常标识符和与之对应的操作，前者可用vector number表示，而后者通常叫interrupt/exception handler。
它的特性有：
	- it is used as a jump table by the CPU when a given vector is triggered
	- it is an array of 256 x 8 bytes entries
	- may reside anywhere in physical memory
	- processor locates IDT by the means of IDTR

IDT的具体定义在arch\x86\include\asm\irq_vectors.h下：
	* Vectors   0 ...  31 : system traps and exceptions - hardcoded events
	* Vectors  32 ... 127 : device interrupts
	* Vector  128         : legacy int80 syscall interface
	* Vectors 129 ... INVALIDATE_TLB_VECTOR_START-1 except 204 : device interrupts
	* Vectors INVALIDATE_TLB_VECTOR_START ... 255 : special interrupts
	// 64-bit x86 has per CPU IDT tables, 32-bit has one shared IDT table.

x86上IDT为8字节，也可以说是named gate，一般有三种门(gate):
	- interrupt gate, 包含中断/异常处理例程的地址，执行例程时会禁用maskable interrupts，也就是清除IF标志
	- trap gate，类似于interrupt gate，但执行例程时不会禁用maskable interrupts
	- task gate，linux未使用

IDT条目中的关键信息：
	- 段选择子segment selector，
	- offset，代码段内的偏移量
	-	T，代表gate的类型
	-	DPL，访问段内容的最小优先级
	63---------------------47--------42-----------------------32
	|		offset(16-31)			 |P|DPL|	 |T|											|
	|		segment selector	 |			offset(0-15)								|
	31--------------------------------------------------------0


# 1.3 内部实现
//refer: http://sklinuxblog.blogspot.com/2015/12/linux-interrupt-handling-for-x86-systems.html

arch/x86/kernel/irq.c
	common_interrupt
		handle_irq
			run_irq_on_irqstack_cond	// x86_64
			__handle_irq	// otherwise

//TODO：use ftrace to show complete call chain



