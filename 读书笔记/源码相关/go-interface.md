refer:
    https://research.swtch.com/interfaces
    https://zhuanlan.zhihu.com/p/27055513

1. go运行时将常用的interface用两种结构表示：
    - 空接口(empty interface)
        var any interface{}
    - 非空接口(non-empty interface)
        var speaker interface {
            speak()
        }

2. 空接口在runtime中用eface描述：
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

3. 非空接口(non-empty interface)在runtime中用iface描述：
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
    ```