package gdata

import (
	"fmt"
	"testing"
)

func TestPathParamSet(t *testing.T) {
	list := NewSkipListInt()
	for i:=0; i < 20; i++ {
		t.Log(list.randomLevel())
	}
}

func TestBasicInt(t *testing.T) {
	list := NewSkipListInt()

	list.Set(4, 4)
	list.Set(1, 1)
	list.Set(5, 5)
	list.Set(2, 2)
	list.Set(6, 6)
	list.Set(3, 3)

	list.Remove(0)
	list.Remove(5)

	t.Log(list.Get(1))
	t.Log(list.Get(2))
	t.Log(list.Get(3))
	t.Log(list.Get(4))
	t.Log(list.Get(5))
	t.Log(list.Get(6))
}

func TestBasicString(t *testing.T) {
	list := NewSkipListString()

	list.Set("d", "d")
	list.Set("a", "a")
	list.Set("b", "b")
	list.Set("e", "e")
	list.Set("c", "c")

	list.Remove(" ")
	list.Remove("e")

	t.Log(list.Get("a"))
	t.Log(list.Get("b"))
	t.Log(list.Get("c"))
	t.Log(list.Get("d"))
	t.Log(list.Get("e"))
}

func initListInt()*SkipListInt {
	list := NewSkipListInt()
	var i int64 = 0
	for ; i <= 1000000; i++ {
		list.Set(i, [1]byte{})
	}
	return list
}

func initListString()*SkipListString {
	list := NewSkipListString()
	var i int64 = 0
	for ; i <= 1000000; i++ {
		list.Set(fmt.Sprintf("%d", i), [1]byte{})
	}
	return list
}

//go test -v -run=none -bench="BenchmarkSet" -benchmem
func BenchmarkSetInt(b *testing.B) {
	b.ReportAllocs()
	list := NewSkipListInt()
	for i := 0; i < b.N; i++ {
		list.Set(int64(i), [1]byte{})
	}
}

func BenchmarkSetString(b *testing.B) {
	b.ReportAllocs()
	list := NewSkipListString()
	for i := 0; i < b.N; i++ {
		list.Set(fmt.Sprintf("%d", i), [1]byte{})
	}
}

//go test -v -run=none -bench="BenchmarkGet" -benchmem
func BenchmarkGetInt(b *testing.B) {
	b.ReportAllocs()
	benchList := initListInt()
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		benchList.Get(int64(i))
	}
	b.StopTimer()
}

func BenchmarkGetString(b *testing.B) {
	b.ReportAllocs()
	benchList := initListString()
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		benchList.Get(fmt.Sprintf("%d", i))
	}
	b.StopTimer()
}
