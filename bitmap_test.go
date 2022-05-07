package gdata

import (
	"fmt"
	"testing"
)

func TestBitMap(t *testing.T) {
	bm := NewBitMap(24)
	fmt.Printf("%08b\n", bm.bits)
	bm.Set(11)
	fmt.Printf("%08b\n", bm.bits)
	has := bm.Check(11)
	fmt.Println(has)
	bm.UnSet(11)
	fmt.Printf("%08b\n", bm.bits)
}

