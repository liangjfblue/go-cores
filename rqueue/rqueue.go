/**
 *
 * @author liangjf
 * @create on 2020/6/9
 * @version 1.0
 */
package rqueue

import (
	"errors"
	"sync"
)

var (
	ErrRQueueNil   = errors.New("ring queue is nil")
	ErrRQueueEmpty = errors.New("ring queue is empty")
	ErrOutOfRQueue = errors.New("out of ring queue")
)

var (
	//循环队列最小大小
	minQueueSize uint32 = 8
)

//IQueue 固定接口
type IQueue interface {
	//插入元素到循环队列
	Push(interface{}) error
	//获取循环队列头元素
	Peek() (interface{}, error)
	//从循环队列弹出一个元素
	Pop() (interface{}, error)
	//获取指定位置元素
	Get(int) (interface{}, error)
	//重设循环队列的大小(2最小倍数)
	Resize() error
	//循环队列大小
	Size() int
	//容量
	Cap() int
}

type rQueue struct {
	em                []interface{}
	head, tail, count int
	mu                sync.RWMutex
}

//NewRQueue 创建循环队列
func NewRQueue(initSize uint32) IQueue {
	if initSize < minQueueSize {
		initSize = minQueueSize
	}
	return &rQueue{
		em:    make([]interface{}, min2(initSize)),
		head:  0,
		tail:  0,
		count: 0,
	}
}

//Push 插入元素到循环队列
func (r *rQueue) Push(e interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.count == len(r.em) {
		_ = r.resize(true)
	}

	r.em[r.tail] = e
	r.tail = (r.tail + 1) & (len(r.em) - 1)
	r.count++

	return nil
}

//Peek 获取循环队列头元素
func (r *rQueue) Peek() (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.em == nil {
		return nil, ErrRQueueNil
	}

	if len(r.em) == 0 {
		return nil, ErrRQueueEmpty
	}

	return r.em[r.head], nil
}

//Pop 从循环队列弹出一个元素
func (r *rQueue) Pop() (interface{}, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	e := r.em[r.head]
	r.em[r.head] = nil
	r.head = (r.head + 1) & (len(r.em) - 1)
	r.count--

	//只利用1/4时触发缩容
	if len(r.em) > int(minQueueSize) && r.count<<2 <= len(r.em) {
		_ = r.resize(false)
	}

	return e, nil
}

//Get 获取指定位置元素
func (r *rQueue) Get(index int) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.em == nil {
		return nil, ErrRQueueNil
	}

	if len(r.em) == 0 {
		return nil, ErrRQueueEmpty
	}

	if index > r.count {
		return nil, ErrOutOfRQueue
	}

	return r.em[(r.head+index)&(len(r.em)-1)], nil
}

//Resize 重设循环队列的大小(2最小倍数)
func (r *rQueue) Resize() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.resize(true)
}

//Size 循环队列大小
func (r *rQueue) Size() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.count
}

//容量
func (r *rQueue) Cap() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return cap(r.em)
}

//resize 调整queue
func (r *rQueue) resize(scala bool) error {
	count := uint32(r.count)
	if scala {
		count = count + 1
	}
	nEm := make([]interface{}, min2(count))

	if r.em == nil {
		return ErrRQueueNil
	}

	if len(r.em) == 0 {
		return ErrRQueueEmpty
	}

	if r.tail > r.head {
		copy(nEm, r.em[r.head:r.tail])
	} else {
		n := copy(nEm, r.em[r.head:])
		copy(nEm[n:], r.em[:r.tail])
	}

	r.head = 0
	r.tail = r.count
	r.em = nEm
	return nil
}

//min2 最近的2的倍数
func min2(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}
