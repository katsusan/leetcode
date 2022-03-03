package ratelimit

import (
	"sync"
	"time"
)

/*
常见有漏桶和令牌桶算法：
- 漏桶(leaky bucket)：每过固定时间向向外漏一滴水，接到水就可以继续请求API，否则要等待下一次漏水。
  refer: https://en.wikipedia.org/wiki/Leaky_bucket

- 令牌桶(token bucket)：匀速向固定容量的桶中放置令牌，服务请求时从桶中消耗令牌，没有令牌则等待。
  refer： https://en.wikipedia.org/wiki/Token_bucket
  + A token is added to the bucket every 1/r seconds.
  + The bucket can hold at the most b tokens. If a token arrives when the bucket is full, it is discarded.
  + When a packet (network layer PDU) of n bytes arrives：
	* if at least n tokens are in the bucket, n tokens are removed from the bucket, and the packet is sent to the network.
	* if fewer than n tokens are available, no tokens are removed from the bucket, and the packet is considered to be non-conformant.
  + Let M be the maximum possible transmission rate in bytes/second:
	* Maximum Busrt time is Tmax = b/(M-r) if M > r, else infinite. Tmax is the time for which the rate M is fully utilized.
	* Maximum Busrt Size is Bmax = Tmax * M.

区别；漏桶流出的速率恒定，而令牌桶允许一定程序的并发。业界的限流通常是令牌桶。
example：
	github.com/juju/ratelimit (token bucket)
	go.uber.org/ratelimit (leaky bucket)
*/

type Bucket struct {
	bucketSize   uint64        // capacity of bucket
	fillInterval time.Duration // interval to fill bucket
	quantum      uint64        // count to fill bucket

	mu              sync.Mutex // protects the following fields
	availableTokens uint64     // available tokens, negative when there are waiters
	lastTick        int64      // last tick holding when we know how many available tokens, in nanoseconds
}
