/**
 *
 * @author liangjf
 * @create on 2020/5/29
 * @version 1.0
 */
package main

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/liangjfblue/go-cores/ratelimit/tokenRateLimit"
)

func testTake() {
	//令牌桶容量为20, 每秒放入5个令牌;
	//令牌桶预先有20个令牌, 20个令牌可以被领取, 取完后每秒放入5个, 即每秒最多接收5个请求
	l := tokenRateLimit.New(10000, 1000)
	if err := l.Limiter(); err != nil {
		log.Fatal(err)
		return
	}

	var (
		sum uint64
		wg  sync.WaitGroup
	)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		ii := i
		go func(ii int) {
			for j := 0; j < 10001; j++ {
				//等待令牌, 超时1s
				if l.Take(time.Second) {
					atomic.AddUint64(&sum, 1)
					fmt.Printf("go%d get token %d, now:%s\n",
						ii, j, time.Now().Format("2006-01-02 15:04:05"))
				}
			}
			wg.Done()
		}(ii)
	}

	wg.Wait()

	//get token num:30003
	fmt.Printf("get token num:%d\n", sum)
}

func testTryTake() {
	//令牌桶容量为20, 每秒放入5个令牌;
	//令牌桶预先有20个令牌, 20个令牌可以被领取, 取完后每秒放入5个, 即每秒最多接收5个请求
	l := tokenRateLimit.New(1000, 500)
	if err := l.Limiter(); err != nil {
		log.Fatal(err)
		return
	}

	var (
		sum uint64
		wg  sync.WaitGroup
	)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		ii := i
		go func(ii int) {
			for j := 0; j < 50; j++ {
				//非阻塞
				if l.TryTake() {
					atomic.AddUint64(&sum, 1)
					fmt.Printf("go%d get token %d, now:%s\n",
						ii, j, time.Now().Format("2006-01-02 15:04:05"))
				}
			}
			wg.Done()
		}(ii)
	}

	wg.Wait()

	//get token num:150
	fmt.Printf("get token num:%d\n", sum)
}

func main() {
	//testTake()
	testTryTake()
}
