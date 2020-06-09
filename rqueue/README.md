# 环形队列

## 特性
- 线性安全
- 自动扩容,缩容

## 压测报告
- CPU2.60GHz
- 4核
- 内存8G


```go
goos: linux
goarch: amd64
pkg: rqueue
BenchmarkRQueue_Push-4   	 7429197	       140 ns/op
PASS
ok  	rqueue	3.244s
```


## 其他开源项目压测

> https://github.com/eapache/queue

```go
goos: linux
goarch: amd64
pkg: queue
BenchmarkQueue_Add-4   	26814978	       337 ns/op
PASS
ok  	queue	9.122s
```

> https://github.com/gammazero/deque

```go
goos: linux
goarch: amd64
pkg: github.com/gammazero/deque
BenchmarkPushFront-4       	12806047	       117 ns/op
BenchmarkPushBack-4        	10583301	       158 ns/op
BenchmarkSerial-4          	 8633694	       246 ns/op
BenchmarkSerialReverse-4   	 9675147	       171 ns/op
BenchmarkRotate-4          	   21948	     98061 ns/op
BenchmarkInsert-4          	   18344	    107520 ns/op
BenchmarkRemove-4          	   62634	    141247 ns/op
BenchmarkYoyo-4            	     297	   4078129 ns/op
BenchmarkYoyoFixed-4       	     501	   2475827 ns/op
PASS
ok  	github.com/gammazero/deque	24.961s
```















