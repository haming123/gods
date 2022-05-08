package gdata

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestBloom(t *testing.T) {
	bf := NewBloomFilter(1024)
	bf.Set("aaa")
	bf.Set("bbb")
	t.Log(bf.Check("aaa"))
	t.Log(bf.Check("bbb"))
	t.Log(bf.Check("ccc"))
}

func TestBloom2(t *testing.T) {
	bf := NewBloomFilter()
	mm := make(map[string]bool, 10000)
	for i:=0; i < 10000; i ++ {
		val := rand.Int() % 10000
		val_str := fmt.Sprintf("%d", val)
		bf.Set(val_str)
		mm[val_str] = true
	}
	for i:=0; i < 10000; i ++ {
		val := rand.Int() % 10000
		val_str := fmt.Sprintf("%d", val)
		ret1 := bf.Check(val_str)
		ret2,_ := mm[val_str]
		if ret1 != ret2 {
			t.Log("error:" + val_str)
		}
	}
}
