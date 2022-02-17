
# 0. nashi


# 1. basic structure

Take `net.Tcpconn` as example,

// some global instance(in epoll case), in runtime/netpoll_epoll.go
```Go
var (
	epfd int32 = -1 // epoll descriptor, initialized by epollcreate()

	netpollBreakRd, netpollBreakWr uintptr // for netpollBreak

	netpollWakeSig uint32 // used to avoid duplicate calls of netpollBreak
)
```

// connection descriptor
```Go
type TCPConn struct {
    conn
}
  ↓
type conn struct {
    fd *netFD
}
  ↓
// Network file descriptor.
type netFD struct {
	pfd poll.FD   // from internal/poll

	// immutable until Close
	family      int
	sotype      int
	isConnected bool // handshake completed or use of association with peer
	net         string
	laddr       Addr
	raddr       Addr
}
  ↓
// FD is a file descriptor. The net and os packages use this type as a
// field of a larger type representing a network connection or OS file.
type FD struct {
	// Lock sysfd and serialize access to Read and Write methods.
	fdmu fdMutex
	// System file descriptor. Immutable until Close.
	Sysfd int
	// I/O poller.
	pd pollDesc
	// Writev cache.
	iovecs *[]syscall.Iovec
	// Semaphore signaled when file is closed.
	csema uint32
	// Non-zero if this file has been set to blocking mode.
	isBlocking uint32
	// Whether this is a streaming descriptor, as opposed to a
	// packet-based descriptor like a UDP socket. Immutable.
	IsStream bool
	// Whether a zero byte read indicates EOF. This is false for a
	// message based socket connection.
	ZeroReadIsEOF bool
	// Whether this is a file rather than a network socket.
	isFile bool
}
  ↓
type pollDesc struct {
	runtimeCtx uintptr
}
  ↓
// go to runtime_netpoll.go

```



# 2. netpoller
核心逻辑：
	通过runtime.netpoll与netpoller通信获取待执行的协程列表gList来

相关srcfile(以epoll为例)：   
	/runtime/netpoll.go
	/runtime/netpoll_epoll.go
	/internal/poll/

# 2.1 implement definition
// 所有poller(epoll/kqueue/...)必须实现下列接口
- func netpollinit()	// 初始化poller，只执行一次

- func netpollopen(fd uintptr, pd *pollDesc) int32
	+ edge-triggered notifications for fd
	+ pd argument is to pass back to netpollready when fd is ready.
	+ return an errno value

- func netpoll(delta int64) gList
	+ delay < 0		→		block indefinitely
	+ delay = 0		→		poll without blocking
	+ delay > 0		→		block for up to [delay] nanoseconds
	+ return a list of goroutines built by calling netpollready.

- func netpollBreak()
	+ wake up the network poller

- func netpollIsPollDescriptor(fd uintptr) bool
	+ reports whether fd is a file descriptor used by the poller.

# 2.2 call chain

# 2.2.1 net.DialTCP

// 着重与netpoller交互的部分
```Go
func DialTCP(network string, laddr, raddr *TCPAddr) (*TCPConn, error)			// net/tcpsock.go
	sd := &sysDialer{network: network, address: raddr.String()}
	c, err := sd.dialTCP(context.Background(), laddr, raddr)
	↓
func (sd *sysDialer) dialTCP(ctx context.Context, laddr, raddr *TCPAddr) (*TCPConn, error)		// net/tcpsock_posix.go
	return sd.doDialTCP(ctx, laddr, raddr)
	↓
func (sd *sysDialer) doDialTCP(ctx context.Context, laddr, raddr *TCPAddr) (*TCPConn, error)		// net/tcpsock_posix.go
	fd, err := internetSocket(ctx, sd.network, laddr, raddr, syscall.SOCK_STREAM, 0, "dial", sd.Dialer.Control)
	↓
func internetSocket(ctx context.Context, net string, laddr, raddr sockaddr, sotype, proto int, mode string, ctrlFn func(string, string, syscall.RawConn) error) (fd *netFD, err error)
	// net/ipsock_posix.go
	return socket(ctx, net, family, sotype, proto, ipv6only, laddr, raddr, ctrlFn)
	↓
func socket(ctx context.Context, net string, family, sotype, proto int, ipv6only bool, laddr, raddr sockaddr, ctrlFn func(string, string, syscall.RawConn) error) (fd *netFD, err error)
	// net/sock_posix.go
	if err := fd.dial(ctx, laddr, raddr, ctrlFn); err != nil
	↓
func (fd *netFD) dial(ctx context.Context, laddr, raddr sockaddr, ctrlFn func(string, string, syscall.RawConn) error) error 	// net/sock_posix.go
	if crsa, err = fd.connect(ctx, lsa, rsa); err != nil {
	↓
func (fd *netFD) connect(ctx context.Context, la, ra syscall.Sockaddr) (rsa syscall.Sockaddr, ret error)		// net/fd_unix.go
	if err := fd.pfd.Init(fd.net, true); err != nil
	↓
func (fd *FD) Init(net string, pollable bool) error			// internal/poll/fd_unix.go
	err := fd.pd.init(fd)
	↓
func (pd *pollDesc) init(fd *FD) error				// internal/poll/fd_poll_runtime.go
	serverInit.Do(runtime_pollServerInit)
	ctx, errno := runtime_pollOpen(uintptr(fd.Sysfd))
	↓
//go:linkname poll_runtime_pollOpen internal/poll.runtime_pollOpen
func poll_runtime_pollOpen(fd uintptr) (*pollDesc, int) {				// runtime/netpoll.go
	errno = netpollopen(fd, pd)
	↓
func netpollopen(fd uintptr, pd *pollDesc) int32	{			// runtime/netpoll_epoll.go, 此处poller接口的实现视平台相关
	var ev epollevent
	ev.events = _EPOLLIN | _EPOLLOUT | _EPOLLRDHUP | _EPOLLET
	*(**pollDesc)(unsafe.Pointer(&ev.data)) = pd
	return -epollctl(epfd, _EPOLL_CTL_ADD, int32(fd), &ev)
}	
```

向epfd的监听列表中加入指定的fd，注意ev.data是关联当前fd的pollDesc指针，
(注：ev类型为epoll_event，包含成员events→储存event类型，data→储存数据用于就绪时回传给epoll_wait的caller)
以便之后epoll_wait返回时从ev.data中取出对应的goroutine唤醒之(netpoll_epoll.go::netpoll)。


# 2.2.2 net.TCPConn.Read

```Golang
func (c *conn) Read(b []byte) (int, error)		// net/net.go
	n, err := c.fd.Read(b)
	↓
func (fd *netFD) Read(p []byte) (n int, err error)	// net/fd_posix.go
	n, err = fd.pfd.Read(p)
	↓
func (fd *FD) Read(p []byte) (int, error)			// interal/poll/fd_unix.go
	for {
			// first try read from fd
			n, err := ignoringEINTRIO(syscall.Read, fd.Sysfd, p)
			// if err == syscall.EAGAIN && fd.pd.pollable(): 若系统调用read返回了EAGAIN且该fd是可poll的
			// 说明fd暂时无数据可读，调用waitRead等待通知
			err = fd.pd.waitRead(fd.isFile)		// err == nil则continue，否则代表出错或者读完，直接返回
	↓
func (pd *pollDesc) waitRead(isFile bool) error		// internal/poll/fd_poll_runtime.go
	return pd.wait('r', isFile)
	↓
func (pd *pollDesc) wait(mode int, isFile bool) error		// internal/poll/fd_poll_runtime.go
	res := runtime_pollWait(pd.runtimeCtx, mode)
	↓
//go:linkname poll_runtime_pollWait internal/poll.runtime_pollWait
func poll_runtime_pollWait(pd *pollDesc, mode int) int		// runtime/netpoll.go
	for !netpollblock(pd, int32(mode), false)
	↓
func netpollblock(pd *pollDesc, mode int32, waitio bool) bool
	if waitio || netpollcheckerr(pd, mode) == 0
		gopark(netpollblockcommit, unsafe.Pointer(gpp), waitReasonIOWait, traceEvGoBlockNet, 5)
		// gopark will suspend current goroutine and call netpollblockcommit on g0,
		// if netpollblockcommit returns false, then resumes current goroutine(see park_m);	← means data ready for reading
		// else there are no data arrived，go to other goroutine
```



# 2.3 netpollBreak

netpoller初始化时(netpollinit函数)会用pipe2系统调用创建半双工管道(linux上，详见APUE 15.1)，
netpollBreakRd为读端的fd，netpollBreakWr为写端的fd。

管道创建后用epoll_ctl向epfd的监听列表里加入netpollBreakRd，这样当向netpollBreakWr写入数据
(netpoll_epoll.go::netpollBreak)时会使poller轮训(netpoll_epoll.go::netpoll)返回避免阻塞
太久时间(findrunnable函数)。
// see https://github.com/golang/go/issues/27707


```Go
// netpoll_epoll.go
func netpollinit() {
	/*...*/
	*(**uintptr)(unsafe.Pointer(&ev.data)) = &netpollBreakRd
	errno = epollctl(epfd, _EPOLL_CTL_ADD, r, &ev)

	netpollBreakRd = uintptr(r)
	netpollBreakWr = uintptr(w)
}

// netpoll_epoll.go
func netpollBreak() {
	if atomic.Cas(&netpollWakeSig, 0, 1) {
		for {
			var b byte
			n := write(netpollBreakWr, unsafe.Pointer(&b), 1)		// 写入任意1字节内容
			if n == 1 {
				break
			}
	/*...*/

// netpoll_epoll.go
func netpoll(delay int64) gList {
	n := epollwait(epfd, &events[0], int32(len(events)), waitms)
	for i := int32(0); i < n; i++ {
		/*...*/
		if *(**uintptr)(unsafe.Pointer(&ev.data)) == &netpollBreakRd {
			if ev.events != _EPOLLIN {
				println("runtime: netpoll: break fd ready for", ev.events)
				throw("runtime: netpoll: break fd ready for something unexpected")
			}
			if delay != 0 {
				// netpollBreak could be picked up by a
				// nonblocking poll. Only read the byte
				// if blocking.
				var tmp [16]byte
				read(int32(netpollBreakRd), noescape(unsafe.Pointer(&tmp[0])), int32(len(tmp)))
				atomic.Store(&netpollWakeSig, 0)
			}
			continue
		}
```


// refer: https://go-review.googlesource.com/c/go/+/171824
netpollBreak函数的caller：
- func syscall_runtime_doAllThreadsSyscall(fn func(bool) bool)
- func findrunnable() (gp *g, inheritTime bool)
- func wakeNetPoller(when int64) 


