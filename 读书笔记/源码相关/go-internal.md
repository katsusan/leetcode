https://tiancaiamao.gitbooks.io/go-internals
https://draveness.me/golang/docs
https://lrita.github.io/2017/12/12/golang-asm

tip:
    - go:linkname -> call private functions of other package
        //go:linkname mallocgc runtime.mallocgc
        func mallocgc(~)
    -

# 0. clutter
    preamble in Go, in format of '//go:xxxx'.
    //go:generate
    //go:linkname
    //go:name
    //go:noescape
    //go:norace
    //go:nosplit
    //go:nowritebarrier
    //go:systemstack

    CGO:
    //go:cgo_dynamic_linker
    //go:cgo_export_dynamic
    //go:cgo_export_static
    //go:cgo_import_dynamic
    //go:cgo_import_static
    //go:cgo_ldflag

# 1. go build -gcflags "-N -l -m -S": //go tool compule -help
    -N: disable optimizations
    -l: disable inlining
    -S: print assembly listing
    -m: print optimization decisions
    note: 上面的禁用优化形式只适用于main.go,想对其他库起作用需要加all，即"all=-N -l"

# 2. 逃逸分析
    https://docs.google.com/document/d/1CxgUBPlx9iJzkz9JWkb6tIpTe5q32QDmz8l0BouG0Cw/preview#
    
# 3. SSA
    cmd/compile/internal/ssa/README.md
    https://docs.google.com/document/d/1szwabPJJc4J-igUZU4ZKprOrNRNJug2JPD8OYi3i1K0/edit
    https://pp.info.uni-karlsruhe.de/uploads/publikationen/braun13cc.pdf
    http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.8.1979&rep=rep1&type=pdf
