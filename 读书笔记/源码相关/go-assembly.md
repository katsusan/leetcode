## 汇编常识
1. 寄存器宽度
    - Quad         => %rax(64bits)
    - DoubleWord   => %eax(32bits)
    - Word         => %ax(16bits)
    - Byte         => %al(8bits)

2. 常用操作
    - $0x15     => 代表常数21
    - (%rsp)    => 代表rsp寄存器地址解引用
    - 0x8(%rsp) => 代表rsp+0x8的地址解引用
    - 4(%rbx, %rcx, 8) => 代表%rbx+8*%rcx+4的地址解引用   
特殊情况：
    lea 0x8(%rsp) %rax => rsp+0x8的值直接写入rax，并不解引用

3. flag寄存器
    - CF(Carry Flag): Arithmetic carry/borrow，算术移位/借位
    - PF(Parity Flag): Odd or even number of bits set
    - ZF(Zero Flag): result was zero
    - SF(Sign Flag): Most significant bit was set
    - OF(Overflow Flag): Result doesn't fit into the location

4. 比较与跳转
    - cmp b, a: equals to computing a - b
    - test b, a: equals to computing a & b   
例：
    cmp %rax, %rdx  =>  if %rdx > %rax   
    jg 0x406330     =>  jump to 0x406330   
常用跳转：
    + jmp:  always jump
    + je/jz: jumps if euqals to 0
    + jne/jnz: jumps if not equals to 0
    + jg: jumps if >
    + jge: jumps if >=
    + jl: jumps if <
    + jle: jumps if <=
    + ja: jumps if above (unsigned >)
    + jae: jumps if above or euqal
    + jb: jumps if below (unsigned <)
    + jbe: jumps if below or equal
    + js: jumps if negative
    + jns: jumps if non-negative








