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

func TestRadix4(t *testing.T) {
	var tree Radix
	mm := make(map[string]int, 10000)
	for i:=0; i < 100; i ++ {
		val := rand.Int() % 1000000
		val_str := fmt.Sprintf("A%d", val)
		tree.Set(val_str, val)
		mm[val_str] = val
	}
	for i:=0; i < 100; i ++ {
		val := rand.Int() % 1000000
		val_str := fmt.Sprintf("A%d", val)
		//fmt.Println(val_str)
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