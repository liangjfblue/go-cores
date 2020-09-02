/**
 *
 * @author liangjf
 * @create on 2020/8/24
 * @version 1.0
 */
package lru

import (
	"reflect"
	"testing"
)

type Int int

func (i Int) Len() int64 {
	return int64(reflect.TypeOf(1).Len())
}

type String string

func (s String) Len() int64 {
	return int64(len(s))
}

func TestNewCache(t *testing.T) {
	NewCache(10, nil)
}

func Test_cache_Del(t *testing.T) {
	c := NewCache(10, nil)
	c.Set("a", Int(1))
	if err := c.Del("a"); err != nil {
		t.Fatal(err)
	}
}

func Test_cache_Get(t *testing.T) {
	c := NewCache(10, nil)
	c.Set("a", Int(1))
	if v, ok := c.Get("a"); !ok || v.(*Entry).Value != Int(1) {
		t.Fatal("set fail")
	}
}

func Test_cache_Len(t *testing.T) {
	c := NewCache(32, nil)
	c.Set("a", Int(1))
	c.Set("c", Int(111))
	c.Set("b", String("123"))
	//l := int64(len("a")+len("b")+len("c")) + Int(1).Len() + Int(111).Len() + String("123").Len()
	if 3 != c.Len() {
		t.Fatal("get len fail, want 3 but ", c.Len())
	}
}

func Test_cache_RemoveOldest(t *testing.T) {
	k1, v1 := "a", Int(1)
	k2, v2 := "b", Int(2)
	k3, v3 := "c", Int(3)

	sum := int64(len(k1)+len(k2)) + v1.Len() + v2.Len()

	cc := NewCache(sum, nil)
	cc.Set(k1, v1)
	cc.Set(k2, v2)
	cc.Set(k3, v3)

	if _, ok := cc.Get("a"); ok || cc.Len() != 2 {
		t.Fatal("remove oldest fail")
	}
}

func Test_cache_Set(t *testing.T) {
	c := NewCache(10, nil)
	c.Set("a", Int(1))
	if v, ok := c.Get("a"); !ok || v.(*Entry).Value != Int(1) {
		t.Fatal("set fail")
	}
}
