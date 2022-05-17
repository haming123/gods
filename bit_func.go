package gdata

//提取字节某一位
func GET_BIT(x byte, bit int) byte {
	return (x & (1 << bit)) >> bit
}

//将字节某一位置1
func SET_BIT(x byte, bit int) byte {
	x |= (1 << bit)
	return x
}

//将字节低位置1
func SET_BITS_LOW(x byte, bit int) byte {
	switch bit {
	case 1:
		x |= 0x1
	case 2:
		x |= 0x3
	case 3:
		x |= 0x7
	case 4:
		x |= 0xF
	case 5:
		x |= 0x1F
	case 6:
		x |= 0x3F
	case 7:
		x |= 0x7F
	case 8:
		x |= 0xFF
	}
	return x
}

//将字节高位置1
func SET_BITS_HIGH(x byte, bit int) byte {
	switch bit {
	case 1:
		x |= 0x80
	case 2:
		x |= 0xC0
	case 3:
		x |= 0xE0
	case 4:
		x |= 0xF0
	case 5:
		x |= 0xF8
	case 6:
		x |= 0xFC
	case 7:
		x |= 0xFE
	case 8:
		x |= 0xFF
	}
	return x
}

//清零字节某一位
func CLEAR_BIT(x byte, bit int) byte {
	x &^= (1 << bit)
	return x
}

//清零字节低位
func CLEAR_BITS_LOW(x byte, bit int) byte {
	//x = x >> bit
	//x = x << bit
	switch bit {
	case 1:
		x &^= 0x1
	case 2:
		x &^= 0x3
	case 3:
		x &^= 0x7
	case 4:
		x &^= 0xF
	case 5:
		x &^= 0x1F
	case 6:
		x &^= 0x3F
	case 7:
		x &^= 0x7F
	case 8:
		x &^= 0xFF
	}
	return x
}

//清零字节高位
func CLEAR_BITS_HIGH(x byte, bit int) byte {
	//x = x << bit
	//x = x >> bit
	switch bit {
	case 1:
		x &^= 0x80
	case 2:
		x &^= 0xC0
	case 3:
		x &^= 0xE0
	case 4:
		x &^= 0xF0
	case 5:
		x &^= 0xF8
	case 6:
		x &^= 0xFC
	case 7:
		x &^= 0xFE
	case 8:
		x &^= 0xFF
	}
	return x
}

//两个字节对比，返回前缀的长度
//返回：0 完全不同
//返回：8 完全相同
func GetPrefixBitLength(b1 byte, b2 byte) int {
	var bb byte = 1
	for i:=0; i < 8; i++ {
		if b1 & bb != b2 & bb {
			return i
		}
		bb = bb << 1
	}
	return 8
}

//两个字节对比，返回前缀的长度
//从beg位开始对比，对比end-beg的长度
func GetPrefixBitLength2(b1 byte, b2 byte, beg int, end int) int {
	pp := 0
	for i:= beg; i < end; i++ {
		if GET_BIT(b1, i) != GET_BIT(b2, i) {
			return pp
		}
		pp += 1
	}
	return pp
}