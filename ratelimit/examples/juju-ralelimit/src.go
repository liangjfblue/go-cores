// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3 with static-linking exception.
// See LICENCE file for details.

// Package ratelimit provides an efficient token bucket implementation
// that can be used to limit the rate of arbitrary things.
// See http://en.wikipedia.org/wiki/Token_bucket.
package main

import (
	"math"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// The algorithm that this implementation uses does computational work
// only when tokens are removed from the bucket, and that work completes
// in short, bounded-constant time (Bucket.Wait benchmarks at 175ns on
// my laptop).
//
// Time is measured in equal measured ticks, a given interval
// (fillInterval) apart. On each tick a number of tokens (quantum) are
// added to the bucket.
//
// When any of the methods are called the bucket updates the number of
// tokens that are in the bucket, and it records the current tick
// number too. Note that it doesn't record the current time - by
// keeping things in units of whole ticks, it's easy to dish out tokens
// at exactly the right intervals as measured from the start time.
//
// This allows us to calculate the number of tokens that will be
// available at some time in the future with a few simple arithmetic
// operations.
//
// The main reason for being able to transfer multiple tokens on each tick
// is so that we can represent rates greater than 1e9 (the resolution of the Go
// time package) tokens per second, but it's also useful because
// it means we can easily represent situations like "a person gets
// five tokens an hour, replenished on the hour".

/**
设计思路:
采用在取令牌时, 根据逻辑时间当前tick和上次取令牌的tick来计算生成这段时间的令牌, 是一种把生成令牌的操作延后到取令牌时刻.
好处: 不用维护一个后台线程来生成令牌加入令牌桶中, 只需要在取令牌时计算这段时间需要生成的令牌即可, 实现较简单
坏处: 对比后台线程定时生成令牌的方式, 这里会加重取令牌操作的逻辑,更耗时, 如果是频繁的取令牌,
	这种方式就相对多了不必要的判断和计算这段时间所需生成令牌的逻辑


该项目是一个令牌桶算法实现的限流器, 众所周知, 限流器可以用于控流, 限制访问频率, 保护后端服务.

令牌桶算法实现的限流器可以较好的应对突然流量, 也可以较好的平滑请求

juju/ratelimit的实现不是像一般的使用后台线程来生成令牌, 前台来取令牌的方案.



*/

// Bucket represents a token bucket that fills at a predetermined rate.
// Methods on Bucket may be called concurrently.
type Bucket struct {
	clock Clock

	// startTime holds the moment when the bucket was
	// first created and ticks began.
	//首次加令牌的时间点
	startTime time.Time

	// capacity holds the overall capacity of the bucket.
	//令牌桶的容量
	capacity int64

	// quantum holds how many tokens are added on
	// each tick.
	//每次滴答新增的令牌
	quantum int64

	// fillInterval holds the interval between each tick.
	//每次滴答的间隔
	fillInterval time.Duration

	// mu guards the fields below it.
	mu sync.Mutex

	// availableTokens holds the number of available
	// tokens as of the associated latestTick.
	// It will be negative when there are consumers
	// waiting for tokens.
	//令牌桶现有的令牌
	availableTokens int64

	// latestTick holds the latest tick for which
	// we know the number of tokens in the bucket.
	//最后一次滴答的值
	latestTick int64
}

// NewBucket returns a new token bucket that fills at the
// rate of one token every fillInterval, up to the given
// maximum capacity. Both arguments must be
// positive. The bucket is initially full.
//使用默认时钟
func NewBucket(fillInterval time.Duration, capacity int64) *Bucket {
	return NewBucketWithClock(fillInterval, capacity, nil)
}

// NewBucketWithClock is identical to NewBucket but injects a testable clock
// interface.
//注入自定义时钟
func NewBucketWithClock(fillInterval time.Duration, capacity int64, clock Clock) *Bucket {
	return NewBucketWithQuantumAndClock(fillInterval, capacity, 1, clock)
}

// rateMargin specifes the allowed variance of actual
// rate from specified rate. 1% seems reasonable.
const rateMargin = 0.01

// NewBucketWithRate returns a token bucket that fills the bucket
// at the rate of rate tokens per second up to the given
// maximum capacity. Because of limited clock resolution,
// at high rates, the actual rate may be up to 1% different from the
// specified rate.
func NewBucketWithRate(rate float64, capacity int64) *Bucket {
	return NewBucketWithRateAndClock(rate, capacity, nil)
}

// NewBucketWithRateAndClock is identical to NewBucketWithRate but injects a
// testable clock interface.
func NewBucketWithRateAndClock(rate float64, capacity int64, clock Clock) *Bucket {
	// Use the same bucket each time through the loop
	// to save allocations.
	tb := NewBucketWithQuantumAndClock(1, capacity, 1, clock)
	for quantum := int64(1); quantum < 1<<50; quantum = nextQuantum(quantum) {
		fillInterval := time.Duration(1e9 * float64(quantum) / rate)
		if fillInterval <= 0 {
			continue
		}
		tb.fillInterval = fillInterval
		tb.quantum = quantum
		//根据速率计算, 逼近最合适的fillInterval, quantum
		if diff := math.Abs(tb.Rate() - rate); diff/rate <= rateMargin {
			return tb
		}
	}
	panic("cannot find suitable quantum for " + strconv.FormatFloat(rate, 'g', -1, 64))
}

// nextQuantum returns the next quantum to try after q.
// We grow the quantum exponentially, but slowly, so we
// get a good fit in the lower numbers.
//慢速指数增长, 到达一定程度递增
func nextQuantum(q int64) int64 {
	q1 := q * 11 / 10
	if q1 == q {
		q1++
	}
	return q1
}

// NewBucketWithQuantum is similar to NewBucket, but allows
// the specification of the quantum size - quantum tokens
// are added every fillInterval.
func NewBucketWithQuantum(fillInterval time.Duration, capacity, quantum int64) *Bucket {
	return NewBucketWithQuantumAndClock(fillInterval, capacity, quantum, nil)
}

// NewBucketWithQuantumAndClock is like NewBucketWithQuantum, but
// also has a clock argument that allows clients to fake the passing
// of time. If clock is nil, the system clock will be used.
func NewBucketWithQuantumAndClock(fillInterval time.Duration, capacity, quantum int64, clock Clock) *Bucket {
	if clock == nil {
		clock = realClock{}
	}
	if fillInterval <= 0 {
		panic("token bucket fill interval is not > 0")
	}
	if capacity <= 0 {
		panic("token bucket capacity is not > 0")
	}
	if quantum <= 0 {
		panic("token bucket quantum is not > 0")
	}

	b := &Bucket{
		clock:        clock,
		startTime:    clock.Now(),
		latestTick:   0,
		fillInterval: fillInterval,
		capacity:     capacity,
		quantum:      quantum,
	}

	atomic.AddInt64(&b.availableTokens, capacity)

	return b
}

// Wait takes count tokens from the bucket, waiting until they are
// available.
//阻塞等待获取count个令牌
func (tb *Bucket) Wait(count int64) {
	if d := tb.Take(count); d > 0 {
		tb.clock.Sleep(d)
	}
}

// WaitMaxDuration is like Wait except that it will
// only take tokens from the bucket if it needs to wait
// for no greater than maxWait. It reports whether
// any tokens have been removed from the bucket
// If no tokens have been removed, it returns immediately.
//指定等待时间
func (tb *Bucket) WaitMaxDuration(count int64, maxWait time.Duration) bool {
	d, ok := tb.TakeMaxDuration(count, maxWait)
	if d > 0 {
		tb.clock.Sleep(d)
	}
	return ok
}

//阻塞等待的最大时间
const infinityDuration time.Duration = 0x7fffffffffffffff

// Take takes count tokens from the bucket without blocking. It returns
// the time that the caller should wait until the tokens are actually
// available.
//
// Note that if the request is irrevocable - there is no way to return
// tokens to the bucket once this method commits us to taking them.
func (tb *Bucket) Take(count int64) time.Duration {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	d, _ := tb.take(tb.clock.Now(), count, infinityDuration)
	return d
}

// TakeMaxDuration is like Take, except that
// it will only take tokens from the bucket if the wait
// time for the tokens is no greater than maxWait.
//
// If it would take longer than maxWait for the tokens
// to become available, it does nothing and reports false,
// otherwise it returns the time that the caller should
// wait until the tokens are actually available, and reports
// true.
func (tb *Bucket) TakeMaxDuration(count int64, maxWait time.Duration) (time.Duration, bool) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.take(tb.clock.Now(), count, maxWait)
}

// TakeAvailable takes up to count immediately available tokens from the
// bucket. It returns the number of tokens removed, or zero if there are
// no available tokens. It does not block.
func (tb *Bucket) TakeAvailable(count int64) int64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.takeAvailable(tb.clock.Now(), count)
}

// takeAvailable is the internal version of TakeAvailable - it takes the
// current time as an argument to enable easy testing.
//统一底层入口
func (tb *Bucket) takeAvailable(now time.Time, count int64) int64 {
	if count <= 0 {
		return 0
	}
	//计算当前tick的令牌数
	tb.adjustavailableTokens(tb.currentTick(now))
	if tb.availableTokens <= 0 {
		return 0
	}
	//判断是否足够
	if count > tb.availableTokens {
		count = tb.availableTokens
	}
	//消耗的令牌
	tb.availableTokens -= count
	return count
}

// Available returns the number of available tokens. It will be negative
// when there are consumers waiting for tokens. Note that if this
// returns greater than zero, it does not guarantee that calls that take
// tokens from the buffer will succeed, as the number of available
// tokens could have changed in the meantime. This method is intended
// primarily for metrics reporting and debugging.
func (tb *Bucket) Available() int64 {
	return tb.available(tb.clock.Now())
}

// available is the internal version of available - it takes the current time as
// an argument to enable easy testing.
func (tb *Bucket) available(now time.Time) int64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.adjustavailableTokens(tb.currentTick(now))
	return tb.availableTokens
}

// Capacity returns the capacity that the bucket was created with.
func (tb *Bucket) Capacity() int64 {
	return tb.capacity
}

// Rate returns the fill rate of the bucket, in tokens per second.
func (tb *Bucket) Rate() float64 {
	//计算生成令牌的速率
	return 1e9 * float64(tb.quantum) / float64(tb.fillInterval)
}

// take is the internal version of Take - it takes the current time as
// an argument to enable easy testing.
func (tb *Bucket) take(now time.Time, count int64, maxWait time.Duration) (time.Duration, bool) {
	if count <= 0 {
		return 0, true
	}

	//逻辑时钟, 当前时刻比初始化时过了多少个tick
	tick := tb.currentTick(now)
	//计算当前合适的可用令牌
	tb.adjustavailableTokens(tick)
	//消耗count个令牌
	avail := tb.availableTokens - count
	//可用令牌的个数
	if avail >= 0 {
		tb.availableTokens = avail
		return 0, true
	}
	// Round up the missing tokens to the nearest multiple
	// of quantum - the tokens won't be available until
	// that tick.

	// endTick holds the tick when all the requested tokens will
	// become available.
	//通过计算还欠多少个令牌计算出需要等多少tick, 然后计算出这么多个tick的下个时刻是多少, 然后返回等待/睡眠这么长的是时间
	//-avail, 欠的令牌的个数
	//-avail+tb.quantum-1	欠的令牌+一次生产的令牌个数-1
	//(-avail+tb.quantum-1)/tb.quantum	需要多少tick才生产出欠的令牌
	endTick := tick + (-avail+tb.quantum-1)/tb.quantum
	//计算等待生产足够欠的令牌的下个即刻
	endTime := tb.startTime.Add(time.Duration(endTick) * tb.fillInterval)
	//需要等待的时间
	waitTime := endTime.Sub(now)
	if waitTime > maxWait {
		return 0, false
	}
	//当前的令牌(欠的照样算)
	tb.availableTokens = avail
	return waitTime, true
}

// currentTick returns the current time tick, measured
// from tb.startTime.
func (tb *Bucket) currentTick(now time.Time) int64 {
	//当前距离初始化时的时间差/滴答时间间隔, 就是当前的tick
	return int64(now.Sub(tb.startTime) / tb.fillInterval)
}

// adjustavailableTokens adjusts the current number of tokens
// available in the bucket at the given time, which must
// be in the future (positive) with respect to tb.latestTick.
func (tb *Bucket) adjustavailableTokens(tick int64) {
	//当前可用令牌已大于令牌桶容量, 不做处理, 直接返回
	if tb.availableTokens >= tb.capacity {
		return
	}
	//计算滴答时间差 * quantum 就是这段时间需要新增的令牌
	tb.availableTokens += (tick - tb.latestTick) * tb.quantum
	//超过容量需抛弃多余的
	if tb.availableTokens > tb.capacity {
		tb.availableTokens = tb.capacity
	}
	//更新tick
	tb.latestTick = tick
	return
}

// Clock represents the passage of time in a way that
// can be faked out for tests.
type Clock interface {
	// Now returns the current time.
	Now() time.Time
	// Sleep sleeps for at least the given duration.
	Sleep(d time.Duration)
}

// realClock implements Clock in terms of standard time functions.
type realClock struct{}

// Now implements Clock.Now by calling time.Now.
func (realClock) Now() time.Time {
	return time.Now()
}

// Now implements Clock.Sleep by calling time.Sleep.
func (realClock) Sleep(d time.Duration) {
	time.Sleep(d)
}
