refer:
  https://go.dev/blog/context-and-structs

# 1. introduction

Context provides a means of transmitting deadlines, caller cancellations, and other 
request-scoped values across API boundaries and between processes.

Contexts should not be stored inside a struct type, but instead passed to each function
that needs it.

some rules:
  - Prefer contexts passed as arguments
  - Storing context in structs leads to confusion and often unnecessary allocations.(http.Request is an exception)
  - Contexts should be passed to functions that may block
  - Contexts should be shared across goroutines
  - Contexts should be canceled
  - Contexts should be canceled when the program exits


# 2. API&source

definition:
```Go
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

type emptyCtx int
var (
	background = new(emptyCtx)
	todo       = new(emptyCtx)
)
```

# 2.1 context.Background()

```Go
func Background() Context {
	return background
}
```

返回non-nil&empty context,不会被cancel，没有value也没有deadline,
通常用在main function/initiation/tests/top level context.

# 2.2 context.TODO()

```Go
func TODO() Context {
	return todo
}
```

返回non-nil&empty context, 通常在不知道该用何种context时使用.

# 2.3 context.WithValue(ctx Context, key, value interface{}) Context

```Go
func WithValue(parent Context, key, val interface{}) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if key == nil {
		panic("nil key")
	}
	if !reflectlite.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}
	return &valueCtx{parent, key, val}
}

type valueCtx struct {
    Context
    key, val interface{}
}
```

usage:
  type name string
  var tom name = "Tom"
  ctx := context.WithValue(ctx, tom, "21") // 注意key最好不要内置类型，防止与其它package的key冲突，同时key必须是Comparable的。
  fmt.Println(ctx.Value(tom)) // "21" // 这里ctx.Value只要判断tom==ctx.key就行了

# 2.4 context.WithCancel(parent Context) (ctx Context, cancel CancelFunc)

```Go
func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	c := newCancelCtx(parent)
	propagateCancel(parent, &c)
	return &c, func() { c.cancel(true, Canceled) }
}

func newCancelCtx(parent Context) cancelCtx {
	return cancelCtx{Context: parent}
}

type cancelCtx struct {
	Context

	mu       sync.Mutex            // protects following fields
	done     atomic.Value          // of chan struct{}, created lazily, closed by first cancel call
	children map[canceler]struct{} // set to nil by the first cancel call
	err      error                 // set to non-nil by the first cancel call
}


```

WithCancel返回一个context，该context可以被cancel，并且可以获取到cancel函数。
当cancel()时该context的err会被设置为Canceled，所有的子context也会被cancel。

usage:
  ctx, cancel := context.WithCancel(ctx)
  go func() {
    time.Sleep(time.Second)
    cancel()
  }()
  <-ctx.Done() // 等待ctx被cancel

ctx, cancel := context.WithCancel(ctx)
  ↓
cancel()
  ↓
c[type:cancelCtx].cancel(true, Canceled)  // true means remove from parent's children
  ↓
close(c.done)   // to unblock select which listens to c.Done()
child := range c.children => c[type:cancelCtx].cancel(true, Canceled) // cancels all children
removeFromParent => removeChild(c.Context, c)


# 2.5 context.WithDeadline(parent Context, deadline time.Time) (ctx Context, cancel CancelFunc)

```Go
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if cur, ok := parent.Deadline(); ok && cur.Before(d) {
		// The current deadline is already sooner than the new one.
		return WithCancel(parent)
	}
	c := &timerCtx{
		cancelCtx: newCancelCtx(parent),
		deadline:  d,
	}
	propagateCancel(parent, c)
	dur := time.Until(d)
	if dur <= 0 {
		c.cancel(true, DeadlineExceeded) // deadline has already passed
		return c, func() { c.cancel(false, Canceled) }
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err == nil {
		c.timer = time.AfterFunc(dur, func() {
			c.cancel(true, DeadlineExceeded)
		})
	}
	return c, func() { c.cancel(true, Canceled) }
}

type timerCtx struct {
	cancelCtx
	timer *time.Timer // Under cancelCtx.mu.

	deadline time.Time
}
```

核心是`c.timer = time.AfterFunc(dur, func() {c.cancel(true, DeadlineExceeded)})`,
经过dur时间后调用c.cancel(true, DeadlineExceed).



# 2.6 context.WithTimeout(parent Context, timeout time.Duration) (ctx Context, cancel CancelFunc)

```Go
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}
```

same behavior as WithDeadline.





