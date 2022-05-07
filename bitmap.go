package gdata

type BitMap struct {
	bits []byte
	vmax uint
}

func NewBitMap(max_val ...uint) *BitMap {
	var max uint = 8192
	if len(max_val) > 0 && max_val[0] > 0 {
		max = max_val[0]
	}

	bm := &BitMap{}
	bm.vmax = max
	sz := (max + 7) / 8
	bm.bits = make([]byte, sz, sz)
	return bm
}

func (bm *BitMap)Set(num uint) {
	if num > bm.vmax {
		bm.vmax += 1024
		if bm.vmax < num {
			bm.vmax = num
		}

		dd := int(num+7)/8 - len(bm.bits)
		if dd > 0 {
			tmp_arr := make([]byte, dd, dd)
			bm.bits = append(bm.bits, tmp_arr...)
		}
	}

	//将1左移num%8后，然后和以前的数据做|，这样就替换成1了
	bm.bits[num/8] |= 1 << (num%8)
}

func (bm *BitMap)UnSet(num uint) {
	if num > bm.vmax {
		return
	}
	//&^:将1左移num%8后，然后进行与非运算，将运算符左边数据相异的位保留，相同位清零
	bm.bits[num/8] &^= 1 << (num%8)
}

func (bm *BitMap)Check(num uint) bool {
	if num > bm.vmax {
		return false
	}
	//&:与运算符，两个都是1，结果为1
	return bm.bits[num/8] & (1 << (num%8)) != 0
}

