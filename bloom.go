package gdata

type BloomFilter struct {
	bset *BitMap
	size uint
}

func NewBloomFilter(size_val ...uint) *BloomFilter {
	var size uint = 1024*1024
	if len(size_val) > 0 && size_val[0] > 0 {
		size = size_val[0]
	}

	bf := &BloomFilter{}
	bf.bset = NewBitMap(size)
	bf.size = size
	return bf
}

//hash函数
var seeds = []uint{3011, 3017, 3031}
func (bf *BloomFilter)hashFun(seed uint, value string) uint64 {
	hash := uint64(seed)
	for i := 0; i < len(value); i++ {
		hash = hash*33 + uint64(value[i])
	}
	return hash
}

//添加元素
func (bf *BloomFilter)Set(value string) {
	for _, seed := range seeds {
		hash := bf.hashFun(seed, value)
		hash = hash % uint64(bf.size)
		bf.bset.Set(uint(hash))
	}
}

//判断元素是否存在
func (bf *BloomFilter)Check(value string) bool {
	for _, seed := range seeds {
		hash := bf.hashFun(seed, value)
		hash = hash % uint64(bf.size)
		ret := bf.bset.Check(uint(hash))
		if !ret {
			return false
		}
	}
	return true
}
