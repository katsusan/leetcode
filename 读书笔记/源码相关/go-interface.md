refer:
    https://research.swtch.com/interfaces
    https://zhuanlan.zhihu.com/p/27055513

# 1. go运行时将常用的interface用两种结构表示：
    - 空接口(empty interface)
        var any interface{}
    - 非空接口(non-empty interface)
        var speaker interface {
            speak()
        }

# 2. 空接口在runtime中用eface描述：
```Go
type eface struct {
    _type *_type        //类型信息
    data  unsafe.Pointer    //数据指针
}

type _type struct {
    size       uintptr
    ptrdata    uintptr // size of memory prefix holding all pointers
    hash       uint32
    tflag      tflag
    align      uint8
    fieldAlign uint8
    kind       uint8
    // function for comparing objects of this type
    // (ptr to object A, ptr to object B) -> ==?
    equal func(unsafe.Pointer, unsafe.Pointer) bool
    // gcdata stores the GC type data for the garbage collector.
    // If the KindGCProg bit is set in kind, gcdata is a GC program.
    // Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
    gcdata    *byte
    str       nameOff
    ptrToThis typeOff
}
```

# 3. 非空接口(non-empty interface)在runtime中用iface描述：
```Go
type iface struct {
    tab  *itab
    data unsafe.Pointer
}

type itab struct {
    inter *interfacetype
    _type *_type
    hash  uint32 // copy of _type.hash. Used for type switches.
    _     [4]byte
    fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}

type interfacetype struct {
	typ     _type
	pkgpath name
	mhdr    []imethod
}
```

# 4. 相关函数签名

```golang
// 由编译时的convFuncName函数决定, 其中T代表具体类型，I代表非空接口，E代表空接口
func convI2I(inter *interfacetype, i iface) (r iface)   // interface转interface
func convT16(val uint16) (x unsafe.Pointer)     // val为2字节大小且2字节对齐，比如uint16
func convT32(val uint32) (x unsafe.Pointer)     // val为4字节大小且4字节对齐，比如uint32
func convT64(val uint64) (x unsafe.Pointer)     // val为64字节大小且8字节对齐，比如uint64
func convTstring(val string) (x unsafe.Pointer) // val为string型或基于string型的自定义类型(type name string这种)
func convTslice(val []byte) (x unsafe.Pointer)  // val为切片或基于切片的自定义类型
// 非上述情况的话则视源变量是否含有指针以及转换目标为空接口还是非空接口调用以下函数
func convT2Enoptr(t *_type, elem unsafe.Pointer) (e eface)
func convT2E(t *_type, elem unsafe.Pointer) (e eface) 
func convT2I(tab *itab, elem unsafe.Pointer) (i iface) 
func convT2Inoptr(tab *itab, elem unsafe.Pointer) (i iface)

// 一般接口类型不能直接赋值给T或E，必须调用下面的断言函数转换之后方可使用
func assertI2I(inter *interfacetype, i iface) (r iface)             // dst := src.(worker), src为非空接口
func assertI2I2(inter *interfacetype, i iface) (r iface, b bool)    // dst, ok := src.work(worker)，src为非空接口
func assertE2I(inter *interfacetype, e eface) (r iface)             // dst := e.(worker)，e为interface{}变量
func assertE2I2(inter *interfacetype, e eface) (r iface, b bool)    // dst, ok := e.(worker),e为interf{}变量 

```


