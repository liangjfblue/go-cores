/**
 *
 * @author liangjf
 * @create on 2020/6/9
 * @version 1.0
 */
package rqueue

import "testing"

func TestNewRQueue(t *testing.T) {
	q := NewRQueue(16)
	t.Log(q.Size())
}

func TestRQueue_Push(t *testing.T) {
	q := NewRQueue(16)

	for i := 0; i < 27; i++ {
		_ = q.Push(i)
	}

	t.Log(q.Size())
	t.Log(q.Cap())
}

func TestRQueue_Get(t *testing.T) {
	q := NewRQueue(16)

	for i := 0; i < 27; i++ {
		_ = q.Push(i)
	}

	for i := 0; i < 5; i++ {
		t.Log(q.Get(i))
	}
}

func TestRQueue_Peek(t *testing.T) {
	q := NewRQueue(16)

	for i := 0; i < 27; i++ {
		_ = q.Push(i)
	}

	t.Log(q.Size())
	for i := 0; i < 5; i++ {
		t.Log(q.Peek())
	}
	t.Log(q.Size())
}

func TestRQueue_Pop(t *testing.T) {
	q := NewRQueue(16)

	for i := 0; i < 33; i++ {
		_ = q.Push(i)
	}

	t.Log(q.Size(), q.Cap())
	for i := 0; i < 28; i++ {
		t.Log(q.Pop())
	}
	t.Log(q.Size(), q.Cap())
}

func TestRQueue_Size(t *testing.T) {
	q := NewRQueue(16)

	for i := 0; i < 33; i++ {
		_ = q.Push(i)
	}
	t.Log(q.Size())
}

func TestRQueue_Cap(t *testing.T) {
	q := NewRQueue(16)

	for i := 0; i < 33; i++ {
		_ = q.Push(i)
	}
	t.Log(q.Cap())
}

//go test -test.bench=".*"
func BenchmarkRQueue_Push(b *testing.B) {
	q := NewRQueue(16)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = q.Push(i)
	}

	b.StopTimer()
}
