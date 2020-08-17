package main

import (
	"fmt"
	"time"
)

type takeAvailableReq struct {
	time   time.Duration
	count  int64
	expect int64
}

func main() {
	reqs := []takeAvailableReq{
		{
			time:   0,
			count:  5,
			expect: 5,
		}, {
			time:   60 * time.Millisecond,
			count:  5,
			expect: 5,
		}, {
			time:   70 * time.Millisecond,
			count:  1,
			expect: 1,
		}}

	tb := NewBucket(10*time.Millisecond, 5)

	for _, req := range reqs {
		d := tb.Take(req.count)
		fmt.Println(d)
	}

	//for j, req := range reqs {
	//	d := tb.TakeAvailable(req.count)
	//	if d != req.expect {
	//		fmt.Printf("test %d, got %v want %v\n", j, d, req.expect)
	//	}
	//}
}
