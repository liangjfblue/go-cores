/**
 *
 * @author liangjf
 * @create on 2020/5/29
 * @version 1.0
 */
package ratelimit

import "time"

type IRateLimit interface {
	//开始限流
	Limiter() error
	//阻塞等待资源, 支持超时控制
	Take(time.Duration) bool
	//非阻塞获取资源
	TryTake() bool
	//控制速率
	SetRate(int)
	//获取令牌数量
	GetToken() int
}
