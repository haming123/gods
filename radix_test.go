package gdata

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestRadixCRUD(t *testing.T) {
	var tree Radix
	tree.Set("A11", 1)
	t.Log(tree.Get("A11"))
	tree.Set("A1122", 2)
	t.Log(tree.Get("A1122"))
	tree.Set("A33", 3)
	t.Log(tree.Get("A33"))
	tree.Set("A334", 4)
	t.Log(tree.Get("A334"))

	t.Log(tree.GetItems())
	t.Log("\n" + tree.GetNodesInfo('A'))
	tree.Delete("A11")
	t.Log(tree.GetItems())
	t.Log("\n" + tree.GetNodesInfo('A'))
}

func TestRadixLT(t *testing.T) {
	var tree Radix
	tree.Set("Hello Go", 3)
	tree.Set("Hello", 2)
	tree.Set("H1", 1)
	t.Log("\n" + tree.GetNodesInfo('H'))
	t.Log(tree.Get("H1"))
	t.Log(tree.Get("Hello"))
	t.Log(tree.Get("Hello Go"))
}

func TestRadixEQ(t *testing.T) {
	var tree Radix
	tree.Set("Hello", 1)
	tree.Set("H1", 2)
	tree.Set("H1", 3)
	//t.Log("\n" + tree.GetNodesInfo('H'))
	t.Log(tree.Get("H1"))
	t.Log(tree.Get("Hello"))
}

func TestRadixDemo(t *testing.T) {
	var tree Radix
	tree.Set("god", 1)
	tree.Set("goto", 2)
	t.Log("\n" + tree.GetNodesInfo('g'))
}

func TestRadixGT(t *testing.T) {
	var tree Radix
	tree.Set("我", 1)
	tree.Set("我们", 2)
	tree.Set("我和你", 3)
	t.Log("\n" + tree.GetNodesInfo([]byte("我")[0]))
	t.Log(tree.Get("我"))
	t.Log(tree.Get("我们"))
	t.Log(tree.Get("我和你"))
}

/*
func TestRadixBatch(t *testing.T) {
	var tree Radix
	for i:=0; i < 10000000; i ++ {
		val := rand.Int() % 1000000
		val_str := fmt.Sprintf("A%d", val)
		tree.Set(val_str, val)
	}
}

func TestMapBatch(t *testing.T) {
	mm := make(map[string]int)
	for i:=0; i < 10000000; i ++ {
		val := rand.Int() % 1000000
		val_str := fmt.Sprintf("A%d", val)
		mm[val_str] = val
	}
}
*/

func TestRadixMap(t *testing.T) {
	var tree Radix
	mm := make(map[string]int, 10000)
	for i:=0; i < 10000; i ++ {
		val := rand.Int() % 1000000
		val_str := fmt.Sprintf("A%d", val)
		tree.Set(val_str, val)
		mm[val_str] = val
	}
	for i:=0; i < 10000; i ++ {
		val := rand.Int() % 1000000
		val_str := fmt.Sprintf("A%d", val)
		ret := tree.Get(val_str)
		ret1 := 0
		if ret != nil {
			ret1 = ret.(int)
		}
		ret2,_ := mm[val_str]
		if ret1 != ret2 {
			t.Log("error:" + val_str)
		}
	}
}

//go test -v -run=none -bench="BenchmarkMap" -benchmem
func BenchmarkMapSet(b *testing.B) {
	b.StartTimer()
	mm := make(map[string]int)
	for i := 0; i < b.N; i++ {
		val_str := fmt.Sprintf("A%d", i)
		mm[val_str] = i
	}
	b.StopTimer()
}

func BenchmarkMapGet(b *testing.B) {
	mm := make(map[string]int)
	for i:=0; i < 1000000; i ++ {
		val := rand.Int() % 1000000
		val_str := fmt.Sprintf("A%d", val)
		mm[val_str] = val
	}
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		val := i % 1000000
		val_str := fmt.Sprintf("A%d", val)
		_,_ = mm[val_str]
	}
	b.StopTimer()
}

//go test -v -run=none -bench="BenchmarkRadix" -benchmem
func BenchmarkRadixSet(b *testing.B) {
	b.StartTimer()
	var tree Radix
	for i := 0; i < b.N; i++ {
		val_str := fmt.Sprintf("A%d", i)
		tree.Set(val_str, i)
	}
	b.StopTimer()
}

func BenchmarkRadixGet(b *testing.B) {
	var tree Radix
	for i:=0; i < 1000000; i ++ {
		val := rand.Int() % 1000000
		val_str := fmt.Sprintf("A%d", val)
		tree.Set(val_str, val)
	}
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		val := i % 1000000
		val_str := fmt.Sprintf("A%d", val)
		tree.Get(val_str)
	}
	b.StopTimer()
}
