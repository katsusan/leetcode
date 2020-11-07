https://tiancaiamao.gitbooks.io/go-internals
https://draveness.me/golang/docs
https://lrita.github.io/2017/12/12/golang-asm

tip:
    - go:linkname -> call private functions of other package
        //go:linkname mallocgc runtime.mallocgc
        function mallocgc(~)
    -

1. go build -gcflags "-N -l -m -S": //go tool compule -help
    -N: disable optimizations
    -l: disable inlining
    -S: print assembly listing
    -m: print optimization decisions
    note: 上面的禁用优化形式只适用于main.go,想对其他库起作用需要加all，即"all=-N -l"

2. 逃逸分析
    https://docs.google.com/document/d/1CxgUBPlx9iJzkz9JWkb6tIpTe5q32QDmz8l0BouG0Cw/preview#
    


