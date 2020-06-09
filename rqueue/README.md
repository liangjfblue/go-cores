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


















