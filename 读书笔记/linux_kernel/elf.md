refer:
  https://lwn.net/Articles/631631/			// ELF load
  https://www.technovelty.org/linux/plt-and-got-the-key-to-code-sharing-and-dynamic-libraries.html  // PLT&GOT
  https://stevens.netmeister.org/631/elf.html   // ELF overview
  

# 0. base
there’s two types of binaries on any system：
  - statically linked：静态可执行文件里包含自己执行所需的全部代码段并且不依赖任何外部库。
  - dynamically linked：动态链接的可执行文件在运行时需要依赖外部库，比如运行时调用printf函数需要链接libc.so.6。。

如果把printf函数的地址hardcoding如可执行文件的话会有以下问题：
  - 只要动态库一更新就要重新编译所有运行时需要链接它的可执行文件
  - 如果OS启用ASLR加载可执行文件会使动态库的虚拟地址随机化

因此引入了relocation重定位这一操作，在静态编译时由链接完成，而对于动态链接的可执行文件则在运行时链接器dynamic linker
(比如x64下的/lib64/ld-linux-x86-64.so)来解析符号地址，linker在ELF中由.interp段中指定。

对于重定位有几个关键的section：
  - .got: For dynamic binaries, this Global Offset Table holds the addresses of variables which are relocated upon loading.
  - .plt: For dynamic binaries, this Procedure Linkage Table holds the trampoline/linkage code.

# 1. 概念concept
ELF(Executable and linkable Format，再以前叫Extensible Linking Format)是一种通用的可执行文件/目标文件/共享库/coredump文件格式。
首次出现在System V Release 4(SVR4)的ABI规范里，之后成为Unix系OS的二进制文件标准。

# 2. 构成layout
- ELF header：ELF32/ELF64分别为52/64字节
  - magic number(4 bytes): 0x7f 45 4c 46
  - binary class(1 byte): 0x01→ELF32, 0x02→ELF64
  - binary endianness(1 byte): 0x01→little endianness, 0x02→big endianness
  - ELF version(1 byte): usually set to 0x01
  - OS ABI(1 byte): 0x00→System V, 0x01→HP-UX, ...
  - ABI Version(1 byte): depends on target ABI
  - pad(7 byte): 0x00 00 00 00 00 00 00
  - // 16 byte divding line

  - ELF type(2 bytes): 0x00→ET_NONE, 0x01→ET_REL(object file), 0x02→ET_EXEC(static-linked ELF), 0x03→ET_DYN(dynamically-linked ELF), 0x04→ET_CORE, 0xfe00/feff/ff00/ffff...
  - machine type(2 bytes): target instruction set architecture. 0x00→None, 0x03→x86, 0x3e→amd64,0xb7→ARM64,0xf3→RISC-V,...
  - version(4 bytes): usually 0x1, 0x0→invalid version / 0x1→current version.
  - ★entry point(4/8 bytes): address of entry point, for example:0x44e990. no entry point
  - ★program header table offset(4/8 bytes): the start of the program header table, usually followes ELF header immediately, ELF32/ELF64→0x34/0x40.
  - section header table offset(4/8 bytes): the start of the section header table
  - flags(4 bytes): interpretation depends on the target architecture
  - header size(2 bytes): the size of this header, ELF32/ELF64→52/64 bytes.
  - ★program header entry size(2 bytes): the size of a program header table entry.
  - ★program header entry number(2 bytes): the number of entries in the program header table.
  - section header entry size(2 bytes):  the size of a section header table entry.
  - section header entry number(2 bytes):  the number of entries in the section header table.
  - section header string table index(2 bytes): index of the section header table entry that contains the section names.
  - // ELF file header over

- Program Header Table(Phdr)
The program header table tells the system how to create a process image.  
// 程序头给出了OS加载ELF的所需信息，只对可执行文件以及目标文件有效。通常起始于program header table offset。
每个Phdr条目的结构如下(ELF32/ELF64→32/56 bytes)：
  - segment type(4 bytes): 0x00→PT_NULL(unused entry), 0x01→PT_LOAD, 0x02→PT_DYNAMIC, 0x03→PT_INTERP, 0x04→PT_NOTE(auxiliary information),
      0x06→PT_PHDR(specifies the location and size of phdr itself, must precede any loadable segment)
  - flags(4 bytes): 64-bit only. PF_X(1)/PF_W(2)/PF_R(4)→executable/writable/readable.
  - offset(4/8 bytes): the offset from the beginning of the file at which the first byte of the segment resides.
  - vaddr(4/8 bytes): the virtual address at which the first byte of the segment resides in memory.
  - paddr(4/8 bytes): on systems where physical address is relevant, reserved for segment's physical address.
  - fsize(4/8 bytes): size in bytes of the segment in the file image. may be 0.
  - memsize(4/8 bytes): size in bytes of the segment in memory. may be 0.
  - flags(4 bytes): 32-bit only. PF_X(1)/PF_W(2)/PF_R(4)→executable/writable/readable.
  - align(4/8 bytes): 0/1 means no alignment. otherwise should be a positive, integral power of 2, with p_vaddr equating p_offset modulus p_align.

- Section Header Table(Shdr)
A file's section header table lets one locate all the file's sections.
// 节头给出了所有节Section的信息，起始于ELF header中的section header table offset
// 大小为ELF32/ELF64→40/64 bytes.
  - sh_name(4 bytes):	An offset to a string in the .shstrtab section that represents the name of this section.
  - sh_type(4 bytes): 0x00→SHT_NULL(unused section), 0x01→SHT_PROGBITS(program data), 0x02→SHT_SYMTAB(symbol table), 0x03→SHT_STRTAB(string table),
      0x04→SHT_RELA(relocation entries with addends), 0x05→SHT_HASH(symbol hash table), 0x06→SHT_DYNAMIC(dynamic linking info),
      0x08→SHT_NOBITS(program space with no data->bss)
  - sh_flags(4/8 bytes): attributes of the section. 0x1→SHF_WRITE(writable), 0x02→SHF_ALLOC(occupies memory during exection),
      0x04→SHF_EXECINSTR(executable machine instructions),...
  - sh_addr(4/8 bytes): virtual address of (loaded) section in memory.
  - sh_offset(4/8 bytes): offset of section in file image.
  - sh_size(4/8 bytes): section size in file image. may be 0.
  - sh_link(4 bytes): contains the section index of an associated section. 
  - sh_info(4 bytes): contains extra information about the section.
  - sh_addralign(4/8 bytes): required alignment of the section. must be a power of two.
  - sh_entsize(4/8 bytes): some sections hold a table of fixed-sized entries, such as a symbol table. gives the size in bytes for each entry. otherwise contains 0.

- Data

以ELF64为例，通常情况下具体分布为：
64字节ELF头 + 56字节 x m个程序头Phdr + Code + Data + 64字节 x n个节头Shdr

# 2. 实践practice

# 2.1 ELF layout
from:  http://www.muppetlabs.com/~breadbox/software/tiny/teensy.html

首先使用gcc编译一个极小的c程序tiny.c，执行O3优化、去除符号表、禁用pie后得到6016字节的ELF64。

```c
// tiny.c
int main()
{
  return 42;
}
```

```
➜  cground gcc -Wall -O3 -s -no-pie tiny.c
➜  cground ./a.out; echo $?
42
➜  cground wc -c a.out
6016 a.out
➜  cground size a.out
   text    data     bss     dec     hex filename
   1004     472       8    1484     5cc a.out
```

objdump反汇编后可以看出这样直接编译出的ELF里由__libc_start_main进行程序初始化以及结束后的善后操作，
因此可以进一步精简。   
<br>   
<hr>
<br>
我们可以直接越过libc直接进行exit系统调用退出程序，如下：
(这里只为探究ELF程序的内部结构，实际编程中不推荐直接调用system call)

```
// rax存放系统调用号，x86_64下exit调用号为60，参考：https://filippo.io/linux-syscall-table/
➜  cground more tiny.asm
BITS 64
GLOBAL _start
SECTION .text
_start:
        mov rax, 60
        mov rdi, 42
        syscall

➜  cground nasm -f elf64 tiny.asm
➜  cground gcc -Wall -s -nostdlib -no-pie tiny.o
➜  cground ./a.out; echo $?
42
➜  cground wc -c a.out
528 a.out
➜  cground size a.out
   text    data     bss     dec     hex filename
     48       0       0      48      30 a.out
```

<br>
<hr>
<br>

用readelf可以看出我们编译出的ELF64里还有个.note.gnu.build-id，这个section对于运行时是非必须的，
因此可以再精简：

```
➜  cground gcc -Wall -s -nostdlib -no-pie -Wl,--build-id=none -o a.out-noid tiny.o
➜  cground wc -c a.out-noid
352 a.out-noid
➜  cground ./a.out-noid; echo $?
42
➜  cground size a.out-noid
   text    data     bss     dec     hex filename
     12       0       0      12       c a.out-noid
```

可以看出text段为12字节，包含了两个5字节的mov指令以及一个2字节的syscall指令，已经压缩到极限了。

<br>
<hr>
<br>

想进一步压缩只有从ELF本身的定义格式上做文章，根据#2的ELF layout可以看出，ELF64最少包含一个64字节的ELF头
以及若干56字节的Phdr+code/data。而section虽然包含了一些重要的信息但却非运行必须，因此可以设法去除section以及在
一些填充字段里做文章。由于amd64下系统调用以及指令长度的关系，这里以x86为例进行优化。

```
// 这里用xor+inc实现x86下编号为1的系统调用exit以及操作寄存器低位来缩短指令长度,从12字节缩短到7字节
mov bl, 42    →   b3 2a     ||    bb 2a 00 00 00  ←   mov ebx, 42
xor eax, eax  →   31 c0     ||    b8 01 00 00 00  ←   mov eax, 1
inc eax       →   40            
int 0x80      →   cd 80     ||    cd 80           ←   int 0x80 
```

通过进一步实验可以发现linux加载时并非对ELF头的所有字段都会校验，见下：

```tiny32.asm
  BITS 32
  
                org     0x00200000
  
                db      0x7F, "ELF"             ; e_ident                   ; 7f 45 4c 46
                dd      1                                       ; p_type    ; 01(ELF32) 00(invalid endian) 00(invalid version) 00(SysV)
                dd      0                                       ; p_offset  ; 00(?) 00 00 00
                dd      $$                                      ; p_vaddr   ; 00 00 20 00 (0x00200000)
                dw      2                       ; e_type        ; p_paddr   ; 02 00 (0x0002->ET_EXEC)
                dw      3                       ; e_machine                 ; 03 00 (0x0003->x86)
                dd      _start                  ; e_version     ; p_filesz  ; 20 00 20 00 (0x00200020->invalid version)
                dd      _start                  ; e_entry       ; p_memsz   ; 20 00 20 00 (0x00200020->entry point)
                dd      4                       ; e_phoff       ; p_flags   ; 04 00 00 00 (0x00000004->start of Phdr)
  _start:
                mov     bl, 42                  ; e_shoff       ; p_align   ; start of Shdr
                xor     eax, eax
                inc     eax                     ; e_flags
                int     0x80
                db      0
                dw      0x34                    ; e_ehsize
                dw      0x20                    ; e_phentsize
                db      1                       ; e_phnum
                                                ; e_shentsize
                                                ; e_shnum
                                                ; e_shstrndx
  
  filesize      equ     $ - $$
```

```
➜  cground wc -c a.out
45 a.out
➜  cground ./a.out; echo $?
42
➜  cground strace -t ./a.out
19:00:33 execve("./a.out", ["./a.out"], 0x7ffe9fc2c248 /* 18 vars */) = 0
strace: [ Process PID=2304 runs in 32 bit mode. ]
19:00:33 exit(42)                       = ?
19:00:33 +++ exited with 42 +++
```


这里省略了不必要的section header相关字段(e_shentsize/e_shnum/e_shstrndx)，并且混合了Program header与ELF header，
有几点需要保证：(见https://elixir.bootlin.com/linux/latest/source/fs/binfmt_elf.c#L457)
  - 第5个字节为1，代表自己是ELF32
  - 第17-18字节为0x2，代表自己是static linked ELF
  - 第19-20字节为0x3，代表程序要在x86架构下运行
  - 第45字节phnum代表了Program header的数量，这里应为0x1
  - Program header offset要为0x4，即从ELF的magic number之后开始为Phdr的字段
  - phentsize即每个程序头大小要为固定的0x20
  
# 2.2 GOT/PLT

# 2.2.1 GOT

```
➜  plt_got more test2.c
extern int foo;

int function(void) {
        return foo;
}

➜  plt_got gcc -shared -fPIC -o libtest2.so test2.c
```

查看function函数的指令后可以看出foo的地址为当前rip的偏移0x200a3b处，此例中暂为0x200ff0,
再用readelf可以看出该地址位于.got节中。当dynamic linker加载libtest2.so时进行runtime relocation，
由于是R_X86_64_GLOB_DAT类型，会按照以下规则计算虚拟地址：
  Base Address + Symbol Value + Addend
其中Base Address就是包含这个符号的segment的基地址(启用了prelink的话会在动态库中写入desired address，
否则会在运行时调用mmap由内核返回)。

对于下面.rela.dyn里的foo而言，假设Base Address为0x7fffff400000，则先计算relocation值为
  0x7fffff400000 + 0 + 0 = 0x7fffff400000
然后将这个值写入地址(0x7fffff400000+0x200ff0)处。

```
➜  plt_got objdump -d libtest2.so
...
00000000000005aa <function>:
 5aa:   55                      push   %rbp
 5ab:   48 89 e5                mov    %rsp,%rbp
 5ae:   48 8b 05 3b 0a 20 00    mov    0x200a3b(%rip),%rax        # 200ff0 <foo>
 5b5:   8b 00                   mov    (%rax),%eax
 5b7:   5d                      pop    %rbp
 5b8:   c3                      retq

➜  plt_got readelf --sections libtest2.so
节头：
  [号] 名称              类型             地址              偏移量
       大小              全体大小          旗标   链接   信息   对齐
...
  [16] .got              PROGBITS         0000000000200fd8  00000fd8
       0000000000000028  0000000000000008  WA       0     0     8

➜  plt_got readelf --relocs libtest2.so

重定位节 '.rela.dyn' at offset 0x3d8 contains 8 entries:
  偏移量          信息           类型           符号值        符号名称 + 加数
...
000000200ff0  000400000006 R_X86_64_GLOB_DAT 0000000000000000 foo + 0
```

# 2.2.2 PLT

```
➜  plt_got more test3.c
int foo(void);

int function(void) {
        return foo();
}

➜  plt_got gcc -shared -fPIC -o libtest3.so test3.c
```

```
➜  plt_got objdump -d libtest3.so
[...]
00000000000004c0 <foo@plt>:
 4c0:   ff 25 52 0b 20 00       jmpq   *0x200b52(%rip)        # 201018 <foo>
 4c6:   68 00 00 00 00          pushq  $0x0
 4cb:   e9 e0 ff ff ff          jmpq   4b0 <.plt>
[...] // _GLOBAL_OFFSET_TABLE_为0x201000，也即.got.plt节的起始地址
00000000000004b0 <.plt>:
 4b0:   ff 35 52 0b 20 00       pushq  0x200b52(%rip)        # 201008 <_GLOBAL_OFFSET_TABLE_+0x8>
 4b6:   ff 25 54 0b 20 00       jmpq   *0x200b54(%rip)        # 201010 <_GLOBAL_OFFSET_TABLE_+0x10>
 4bc:   0f 1f 40 00             nopl   0x0(%rax)
[...]
Disassembly of section .got.plt:
0000000000201000 <_GLOBAL_OFFSET_TABLE_>:
  201000:       60                      (bad)
  201001:       0e                      (bad)
  201002:       20 00                   and    %al,(%rax)
        ...
  201018:       c6 04 00 00             movb   $0x0,(%rax,%rax,1)
  20101c:       00 00                   add    %al,(%rax)
        ...
[...]
00000000000005ba <function>:
 5ba:   55                      push   %rbp
 5bb:   48 89 e5                mov    %rsp,%rbp
 5be:   e8 fd fe ff ff          callq  4c0 <foo@plt>
 5c3:   5d                      pop    %rbp
 5c4:   c3                      retq

➜  plt_got readelf --relocs libtest3.so
[...]
重定位节 '.rela.plt' at offset 0x480 contains 1 entry:
  偏移量          信息           类型           符号值        符号名称 + 加数
000000201018  000400000007 R_X86_64_JUMP_SLO 0000000000000000 foo + 0
```

运行时调用链：// 假设映射的段基址segment base address为0x7ffff7bd1000
funtion: call 4c0 <foo@plt>   // 实际应该是segment base address + 4c0，即0x7ffff7bd1000+4c0=0x7ffff7bd14c0
 ↓
foo@plt: jmpq   *0x200b52(%rip) # 201018 <foo> // base address + 0x202018 = 0x7ffff7dd2018(GOT+0x18), 
  而0x7ffff7dd2018处的值正是0x00007ffff7bd14c6，即jmq的下一条指令pushq 0x0。
 ↓
foo@plt: jmpq   4b0 <.plt>    // 跳转到.plt
 ↓
<.plt>: pushq  0x200b52(%rip)   // 即把GOT+0x8指向的值压栈，即pc+0x200b52=0x7ffff7dd2008(GOT+0x8),把0x7ffff7dd2008处的值0x00007ffff7ff6000压栈
        jmpq   *0x200b54(%rip)  // 跳转到pc+0x200b54=0x7ffff7dd2010(GOT+0x10)指向的地址即0x00007ffff7dea7a0，也就是_dl_runtime_resolve_fxsave，动态库的解析函数
      
注：// 下述是_GLOBAL_OFFSET_TABLE_(GOT)的值
>>> x /4gx 0x7ffff7dd2000
0x7ffff7dd2000: 0x0000000000200e60      0x00007ffff7ff6000
0x7ffff7dd2010: 0x00007ffff7dea7a0      0x00007ffff7bd14c6
按照http://users.eecs.northwestern.edu/~kch479/docs/notes/linking.html所述：
  - GOT[0]是当前程序的.dynamic节的_DYNAMIC符号的地址(此处存疑)
  - GOT[1]指向程序链接的所有动态库的符号表构成的链表，当动态库解析时会遍历该链表寻找对应的符号，用LD_PRELOAD可以把动态库加到链表头 
  - GOT[2]是_dl_runtime_resolve的地址，本例中即_dl_runtime_resolve_fxsave in section .text of /lib64/ld-linux-x86-64.so.2
    _dl_runtime_resolve里面有个dl_fixup函数，由它来负责函数解析以及更新GOT条目。
    当找到foo函数的地址后，它会更新上面的0x7ffff7dd2018处的值为真正的地址，这样当调用foo@plt时，第一次jump就能直接跳转到目的函数地址。
    本例中0x7ffff7dd2018处的值会由原来的0x00007ffff7bd14c6(下一条push指令)更新为0x000055555555477a(目的函数foo的地址)。


