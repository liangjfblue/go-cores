/**
 *
 * @author liangjf
 * @create on 2020/5/29
 * @version 1.0
 */
package tokenRateLimit

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/liangjfblue/go-cores/ratelimit"
)

var (
	ErrCapRateLessZero = errors.New("cap/rate less zero")
)

var (
	maxTimeOut = 0x7fffffffffffffff
)

/**
令牌桶算法:
令牌桶算法的原理是系统会以一个恒定的速度往桶里放入令牌，而如果请求需要被处理，则需要先从桶里获取一个令牌，当桶里没有令牌可取时，则拒绝服务。 当桶满时，新添加的令牌被丢弃或拒绝。
令牌桶算法是一个存放固定容量令牌（token）的桶，按照固定速率往桶里添加令牌。令牌桶算法基本可以用下面的几个概念来描述：

令牌将按照固定的速率被放入令牌桶中。比如每秒放10个。
桶中最多存放b个令牌，当桶满时，新添加的令牌被丢弃或拒绝。

当一个请求到达，将从桶中删除1个令牌，接着请求继续进行下一步处理。
如果桶中的令牌不足1个，则不会删除令牌，且该请求将被限流（丢弃/队列等待）
*/
type tokenRateLimit struct {
	cap      int       //最大容量(令牌桶的令牌个数)
	rate     int       //生成令牌速率(每秒生成令牌的个数)
	lastTime time.Time //上次拿令牌时间(用于动态调整rate)
	sync.RWMutex
	tokens int           //令牌(生产者:生产令牌, 消费者:客户程序)
	stop   chan struct{} //停止流程
}

func New(cap, rate int) ratelimit.IRateLimit {
	l := new(tokenRateLimit)

	l.cap = cap
	l.rate = rate
	l.lastTime = time.Now()
	l.tokens = cap
	l.stop = make(chan struct{}, 1)

	if !l.validate() {
		panic(ErrCapRateLessZero)
	}

	go l.producer()

	return l
}

//Wait 阻塞等待资源
func (l *tokenRateLimit) Wait(count int) bool {
	if l.take(count) {
		return true
	}

	t := time.After(time.Duration(maxTimeOut))
	select {
	case <-t:
		return false
	}
}

//WaitWithTimeout 阻塞等待资源, 支持超时控制
func (l *tokenRateLimit) WaitWithTimeout(count int, timeout time.Duration) bool {
	if l.take(count) {
		return true
	}

	if timeout < 0 {
		timeout = 0
	}
	t := time.After(timeout)
	select {
	case <-t:
		return false
	}
}

//TryWait 非阻塞获取资源
func (l *tokenRateLimit) TryWait(count int) bool {
	return l.takeAvailable(count)
}

//SetRate 控制速率
func (l *tokenRateLimit) SetRate(rate int) {
	l.Lock()
	defer l.Unlock()
	l.rate = rate
}

//SetRate 控制速率
func (l *tokenRateLimit) GetRate() int {
	l.RLock()
	defer l.RUnlock()
	return l.rate
}

//GetToken 获取令牌数量
func (l *tokenRateLimit) GetToken() int {
	l.RLock()
	defer l.RUnlock()
	return l.tokens
}

//Stop 停止
func (l *tokenRateLimit) Stop() {
	l.Lock()
	defer l.Unlock()
	if len(l.stop) > 0 {
		return
	}
	l.stop <- struct{}{}
}

//take 阻塞等待资源, 支持超时控制
func (l *tokenRateLimit) take(count int) bool {
	l.Lock()
	defer l.Unlock()

	if len(l.stop) > 0 {
		return true
	}

	if l.tokens >= count {
		l.tokens -= count
		if l.tokens < 0 {
			l.tokens = 0
		}
		fmt.Println("222 [ok ] tokens:", l.tokens, " cap:", l.cap, " rate:", l.rate, "time:", time.Now().Format("2006-01-02 15:04:05"))
		return true
	}
	return false
}

//takeAvailable 非阻塞获取资源内部函数
func (l *tokenRateLimit) takeAvailable(count int) bool {
	l.Lock()
	defer l.Unlock()

	if len(l.stop) > 0 {
		return true
	}

	if l.tokens >= count {
		l.tokens -= count
		if l.tokens < 0 {
			l.tokens = 0
		}
		return true
	}

	return false
}

//validate 检查参数
func (l *tokenRateLimit) validate() bool {
	if l.rate <= 0 || l.cap <= 0 {
		return false
	}
	return true
}

func (l *tokenRateLimit) producer() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-l.stop:
			ticker.Stop()
			return
		case <-ticker.C:
			fmt.Printf("000 producer tokens:%d, rate:%d\n", l.tokens, l.rate)
			l.Lock()
			fmt.Printf("111 producer tokens:%d, rate:%d\n", l.tokens, l.rate)
			l.tokens += l.rate
			if l.tokens > l.cap {
				l.tokens = l.cap
			}
			fmt.Printf("222 producer tokens:%d, rate:%d\n", l.tokens, l.rate)
			l.Unlock()
		}
	}
}
