package gdata

import (
	"encoding/json"
	"math/rand"
	"testing"
)

func TestBPT(t *testing.T) {
	bpt := NewBPTree(4)

	bpt.Set(10, 1)
	bpt.Set(23, 1)
	bpt.Set(33, 1)
	bpt.Set(35, 1)
	bpt.Set(15, 1)
	//bpt.Set(16, 1)
	//bpt.Set(17, 1)
	//bpt.Set(19, 1)
	//bpt.Set(20, 1)

	bpt.Remove(23)

	t.Log(bpt.Get(10))
	t.Log(bpt.Get(15))
	t.Log(bpt.Get(20))

	data, _ := json.MarshalIndent(bpt.GetData(), "", "    ")
	t.Log(string(data))
}

func TestBPTRand(t *testing.T) {
	bpt := NewBPTree(3)

	for i:=0; i < 12; i++ {
		key := rand.Int() % 20 +1
		t.Log(key)
		bpt.Set(int64(key), key)
	}

	data, _ := json.MarshalIndent(bpt.GetData(), "", "    ")
	t.Log(string(data))
}
