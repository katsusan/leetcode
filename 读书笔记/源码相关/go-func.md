refer:
	https://docs.google.com/document/d/1lyPIbmsYbXnpNj57a261hgOYVpNRcgydurVQIyZOz_o/pub

# 0.base
	FUNCDATA $2, $sym(SB)
	// declares that index 2 of the funcdata list should be a pointer to "sym"
	// 代表funcdata列表的索引2的位置为一个指向syn的指针

	PCDATA $3, $45
	// declares that the value with index 3 associated with the current program counter is 45
	// 代表位于相对当前PC的索引3位置的值为45

# 1. pclntab
The memory for the function symbol table, file name table, Func structures, and the data referred to by offset from those 
are all contiguous in memory and recorded as a single symbol, “pclntab”.
// 函数符号表、源文件名表、函数结构以及数据都存在一个连续的内存空间即pclntab，在ELF上体现为section .gopclntab。
// gopclntab以runtime.epclntab结尾，开头为一串magic number：0xfffffffb，之后为2字节的0x00，再之后是1字节的指令大小quantum(x86:1, ARM:4)，
// 再之后的1字节为给出以字节为单位的uintptr的大小。

也就是在AMD64架构下，gopclintab的形式大概如下：
[4] 0xfffffffb
[2] 0x0000
[1] 0x01
[1] 0x08
[8] N (size of function symbol table)
[8] pc0
[8] func0 offset
[8] pc1
[8] func1 offset
...
[8] pcN
[8] funcN offset
[K] data referred to by offset 

每个funcN的结构如下，储存在上述的相对于pclntab的offset处(即pclntab+offset)：
```go
struct Func
{
		uintptr	entry;		// start pc
		int32 name;    		// name (offset to C string)
		int32 args;				// size of arguments passed to function
		int32 frame;			// size of function frame, including saved caller PC
		int32	pcsp;				// pcsp table (offset to pcvalue table)
		int32 pcfile;			// pcfile table (offset to pcvalue table)
		int32 pcln;				// pcln table (offset to pcvalue table)
		int32 nfuncdata;	// number of entries in funcdata list
		int32 npcdata;		// number of entries in pcdata list
};
```



