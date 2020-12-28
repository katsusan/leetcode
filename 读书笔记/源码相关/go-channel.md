refer： https://dave.cheney.net/2014/03/19/channel-axioms

source:
	runtime/chan.go
	runtime/select.go

# 1. channel底层结构

```Go
type hchan struct {
	qcount   uint           // total data in the queue
	dataqsiz uint           // size of the circular queue
	buf      unsafe.Pointer // points to an array of dataqsiz elements
	elemsize uint16
	closed   uint32
	elemtype *_type // element type
	sendx    uint   // send index
	recvx    uint   // receive index
	recvq    waitq  // list of recv waiters
	sendq    waitq  // list of send waiters

	// lock protects all fields in hchan, as well as several
	// fields in sudogs blocked on this channel.
	//
	// Do not change another G's status while holding this lock
	// (in particular, do not ready a G), as this can deadlock
	// with stack shrinking.
	lock mutex
}

type waitq struct {
	first *sudog
	last  *sudog
}

// sudog represents a g in a wait list, such as for sending/receiving
// on a channel.
//
// sudog is necessary because the g ↔ synchronization object relation
// is many-to-many. A g can be on many wait lists, so there may be
// many sudogs for one g; and many gs may be waiting on the same
// synchronization object, so there may be many sudogs for one object.
//
// sudogs are allocated from a special pool. Use acquireSudog and
// releaseSudog to allocate and free them.
type sudog struct {
	// The following fields are protected by the hchan.lock of the
	// channel this sudog is blocking on. shrinkstack depends on
	// this for sudogs involved in channel ops.

	g *g

	next *sudog
	prev *sudog
	elem unsafe.Pointer // data element (may point to stack)

	// The following fields are never accessed concurrently.
	// For channels, waitlink is only accessed by g.
	// For semaphores, all fields (including the ones above)
	// are only accessed when holding a semaRoot lock.

	acquiretime int64
	releasetime int64
	ticket      uint32

	// isSelect indicates g is participating in a select, so
	// g.selectDone must be CAS'd to win the wake-up race.
	isSelect bool

	parent   *sudog // semaRoot binary tree
	waitlink *sudog // g.waiting list or semaRoot
	waittail *sudog // semaRoot
	c        *hchan // channel
}

type chantype struct {
	typ  _type
	elem *_type
	dir  uintptr
}
```

# 2. 所有chan操作函数签名

```Go
func makechan(t *chantype, size int) *hchan

func chanrecv1(c *hchan, elem unsafe.Pointer)       // <-c
func chanrecv2(c *hchan, elem unsafe.Pointer) (received bool)   // _, ok := <-c || for range c

// chanrecv receives on channel c and writes the received data to ep.
// ep may be nil, in which case received data is ignored.
// If block == false and no elements are available, returns (false, false).
// Otherwise, if c is closed, zeros *ep and returns (true, false).
// Otherwise, fills in *ep with an element and returns (true, true).
// A non-nil ep must point to the heap or the caller's stack.
func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool)    //上述两者的底层实现


func chansend1(c *hchan, elem unsafe.Pointer)   // c <- x

/*
 * generic single channel send/recv
 * If block is not nil,
 * then the protocol will not
 * sleep but return if it could
 * not complete.
 *
 * sleep can wake up with g.param == nil
 * when a channel involved in the sleep has
 * been closed.  it is easiest to loop and re-run
 * the operation; we'll see that it's now closed.
 */
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool

/*
select {				if selectnbsend(c, x) {
case c <- x:				foo
	foo			==>		} else {
default:					bar
	bar					}
}
*/
func selectnbsend(c *hchan, elem unsafe.Pointer) (selected bool)
	 chansend(c, elem, false, getcallerpc())

/*
select {				if selectnbrecv(v, c) {
case v := <- c:				foo
	foo			==>		} else {
default:					bar
	bar					}
}
*/
func selectnbrecv(elem unsafe.Pointer, c *hchan) (selected bool)

/*
select {				if c != nil && selectnbrecv2(&v, &ok, c) {
case v, ok := <- c:			foo
	foo			==>		} else {
default:					bar
	bar					}
}
*/
func selectnbrecv2(elem unsafe.Pointer, received *bool, c *hchan) (selected bool)

/*
select {
case <- c1:
	foo
case <- c2:
	bar
}

cas0为select语句里所有case的数组，order0为遍历case所用的原始顺序数组，ncases为case的数量，
返回对应recv/send/default的case数组的位置，以若为recv操作时是否接受到值(正常receive到值为true)
*/
func selectgo(cas0 *scase, order0 *uint16, ncases int) (int, bool)
```

# 3. chan行为一览

--------------------------------------------------------------
操作		nil channel		正常channel			closed channel
--------------------------------------------------------------
<-ch		阻塞			成功或阻塞			读到零值
--------------------------------------------------------------
ch<-		阻塞			成功或阻塞			panic
--------------------------------------------------------------
close(ch)	panic			成功				panic
--------------------------------------------------------------






