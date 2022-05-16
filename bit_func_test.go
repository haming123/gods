package gdata

import (
	"fmt"
	"testing"
)

func TestSET_GET_BIT(t *testing.T) {
	var val byte = 0x00
	val = GET_BIT(val, 3)
	fmt.Printf("%08b\n", val)
	val = SET_BIT(val, 3)
	fmt.Printf("%08b\n", val)
	val = GET_BIT(val, 3)
	fmt.Printf("%08b\n", val)
}

func TestSET_BITS_LOW(t *testing.T) {
	var val byte = 0x00
	for i:=1; i <= 8; i++ {
		val = SET_BITS_LOW(val, i)
		fmt.Printf("%08b\n", val)
	}
}

func TestCLEAR_BITS_LOW(t *testing.T) {
	var val byte = 0xff
	for i:=1; i <= 8; i++ {
		val = CLEAR_BITS_LOW(val, i)
		fmt.Printf("%08b\n", val)
	}
}

func TesSET_BITS_HIGH(t *testing.T) {
	var val byte = 0x00
	for i:=1; i <= 8; i++ {
		val = SET_BITS_HIGH(val, i)
		fmt.Printf("%08b\n", val)
	}
}

func TestCLEAR_BITS_HIGH(t *testing.T) {
	var val byte = 0xff
	for i:=1; i <= 8; i++ {
		val = CLEAR_BITS_HIGH(val, i)
		fmt.Printf("%08b\n", val)
	}
}

func TestGetPrefixBitLength(t *testing.T) {
	var b1 byte = 0b00001111
	var b2 byte = 0b00000011
	pos := GetPrefixBitLength(b1, b2)
	fmt.Println(pos)

	b1 = 0b00001001
	b2 = 0b00011001
	pos = GetPrefixBitLength(b1, b2)
	fmt.Println(pos)

	b1 = 0b1101001
	b2 = 0b1101001
	pos = GetPrefixBitLength(b1, b2)
	fmt.Println(pos)

	b1 = 0b0000000
	b2 = 0b1111111
	pos = GetPrefixBitLength(b1, b2)
	fmt.Println(pos)
}

func TestGetPrefixBitLength2(t *testing.T) {
	var b1 byte
	var b2 byte
	b1 = 0b00001111
	b2 = 0b11111111
	pos := GetPrefixBitLength2(b1, b2, 2, 5)
	fmt.Println(pos)
}