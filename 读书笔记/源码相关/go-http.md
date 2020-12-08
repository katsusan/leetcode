# 1. http Request flow


```Go
func (c *Client) Do(req *Request) (*Response, error)
    c.do(req)
↓
func (c *Client) do(req *Request) (retres *Response, reterr error)
    resp, didTimeout, err = c.send(req, deadline)
↓
func (c *Client) send(req *Request, deadline time.Time) (resp *Response, didTimeout func() bool, err error)
    resp, didTimeout, err = send(req, c.transport(), deadline)
↓
func send(ireq *Request, rt RoundTripper, deadline time.Time) (resp *Response, didTimeout func() bool, err error)
    resp, err = rt.RoundTrip(req)
↓
func (t *Transport) RoundTrip(req *Request) (*Response, error)
    t.roundTrip(req)
↓
func (t *Transport) roundTrip(req *Request) (*Response, error) 
    pconn, err := t.getConn(treq, cm)
↓
func (t *Transport) getConn(treq *transportRequest, cm connectMethod) (pc *persistConn, err error)
    t.queueForDial(w)
↓
func (t *Transport) queueForDial(w *wantConn)
    go t.dialConnFor(w)
↓
func (t *Transport) dialConnFor(w *wantConn)
    pc, err := t.dialConn(w.ctx, w.cm)
↓
func (t *Transport) dialConn(ctx context.Context, cm connectMethod) (pconn *persistConn, err error)
	go pconn.readLoop()
	go pconn.writeLoop()
```



