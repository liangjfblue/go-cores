/**
 *
 * @author liangjf
 * @create on 2020/5/29
 * @version 1.0
 */
package tokenRateLimit

import (
	"errors"
	"sync"
	"time"

	"github.com/liangjfblue/go-cores/ratelimit"
)

var (
	ErrCapRateLessZero = errors.New("cap/rate less zero")
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
	bucket       chan struct{} //令牌桶
	tokenChannel chan struct{} //令牌通道(生产者:生产令牌, 消费者:客户程序)
	stop         chan struct{}
}

func New(cap, rate int) ratelimit.IRateLimit {
	l := new(tokenRateLimit)

	l.cap = cap
	l.rate = rate
	l.lastTime = time.Now()
	l.stop = make(chan struct{}, 1)

	l.bucket = make(chan struct{}, cap)
	for i := 0; i < cap; i++ {
		l.bucket <- struct{}{}
	}

	return l
}

//开始限流
func (l *tokenRateLimit) Limiter() error {
	if !l.validate() {
		return ErrCapRateLessZero
	}

	//生产令牌goroutine
	go l.product()

	return nil
}

//Take 阻塞等待资源, 支持超时控制
func (l *tokenRateLimit) Take(timeout time.Duration) bool {
	t := time.After(timeout)
	for {
		select {
		case <-t:
			return false
		case <-l.bucket:
			return true
		}
	}
}

//TryTake 非阻塞获取资源
func (l *tokenRateLimit) TryTake() bool {
	l.RLock()
	if len(l.bucket) <= 0 {
		l.RUnlock()
		return false
	}
	l.RUnlock()

	<-l.bucket
	return true
}

//SetRate 控制速率
func (l *tokenRateLimit) SetRate(rate int) {
	l.rate = rate
}

//获取令牌数量
func (l *tokenRateLimit) GetToken() int {
	l.RLock()
	defer l.RUnlock()
	return len(l.bucket)
}

//validate 检查参数
func (l *tokenRateLimit) validate() bool {
	if l.rate <= 0 || l.cap <= 0 {
		return false
	}
	return true
}

//product 生产令牌
func (l *tokenRateLimit) product() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-l.stop:
			return
		case <-ticker.C:
			for i := 0; i < l.rate; i++ {
				l.bucket <- struct{}{}
			}
		}
	}
}
