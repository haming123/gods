package gdata

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type RaxNode struct {
	bit_len int
	bit_val []byte
	left 	*RaxNode
	right 	*RaxNode
	val interface{}
}

type Radix struct {
	root [256]*RaxNode
}

func NewRaxNode() *RaxNode {
	nd := &RaxNode{}
	nd.bit_len = 0
	nd.bit_val = nil
	nd.left = nil
	nd.right = nil
	nd.val = nil
	return nd
}

func (this *RaxNode)pathCompare(data []byte, bbeg int) (bool, int) {
	bend := bbeg + int(this.bit_len)
	if bend > len(data) * 8 {
		return false, len(data) * 8
	}

	//起始和终止字节的位置
	cbeg := bbeg / 8; cend := bend / 8
	//起始和终止字节的偏移量
	obeg := bbeg % 8; oend := bend % 8
	for bb := bbeg; bb < bend; {
		//获取两个数组的当前字节位置
		dci := bb / 8
		nci := dci - cbeg

		//获取数据的当前字节以及循环步长
		step := 8
		byte_data := data[dci]
		if dci == cbeg && obeg > 0 {
			//清零不完整字节的低位
			byte_data = CLEAR_BITS_LOW(byte_data, obeg)
			step -= obeg
		}
		if dci == cend && oend > 0 {
			//清零不完整字节的高位
			byte_data = CLEAR_BITS_HIGH(byte_data, 8-oend)
			step -= 8-oend
		}

		//获取结点的当前字节，并与数据的当前字节比较
		byte_node := this.bit_val[nci]
		if byte_data != byte_node {
			return false, len(data)*8
		}

		bb += step
	}

	return true, bend
}

func (this *RaxNode)pathSplit(key []byte, key_pos int, val interface{}) int {
	//与path对应的key数据(去掉已经处理的公共字节)
	data := key[key_pos/8:]
	//key以bit为单位长度（包含开始字节的字节内bit位的偏移量）
	bit_end_key := len(data) * 8
	//path以bit为单位长度（包含开始字节的字节内bit位的偏移量）
	bit_end_path := key_pos % 8 + int(this.bit_len)
	//当前的bit偏移量，需要越过开始字节的字节内bit位的偏移量
	bpos := key_pos % 8
	for ; bpos < bit_end_key && bpos < bit_end_path; {
		ci := bpos / 8
		byte_path := this.bit_val[ci]
		byte_data := data[ci]

		//起始字节的内部偏移量
		beg := 0
		if ci == 0 {
			beg = key_pos % 8
		}
		//终止字节的内部偏移量
		end := 8
		if  ci == bit_end_path / 8 {
			end = bit_end_path % 8
			if end == 0 {
				end = 8
			}
		}

		if beg != 0 || end != 8 {
			//不完整字节的比较，若不等则跳出循环
			num := GetPrefixBitLength2(byte_data, byte_path, beg, end)
			bpos += num
			if num < end - beg {
				break
			}
		} else if byte_data != byte_path {
			//完整字节比较，若不想等，获取bit相同的长度，并跳出循环
			//若相等，则增长相等的bit长度，并继续比较下一个字节
			num := GetPrefixBitLength2(byte_data, byte_path, 0, 8)
			bpos += num
			break
		} else {
			//完整字节比较，相等，则继续比较下一个字节
			bpos += 8
		}
	}

	//当前字节的位置
	char_index := bpos / 8
	//当前字节的bit偏移量
	bit_offset := bpos % 8
	//剩余的path长度
	bit_last_path := bit_end_path - bpos
	//剩余的key长度
	bit_last_data := bit_end_key - bpos

	//key的数据有剩余
	//若path有子结点，则继续处理子结点
	//若path没有子结点，则创建一个key子结点
	var nd_data *RaxNode = nil
	var bval_data byte
	if bit_last_data > 0 {
		//若path有子结点，则退出本函数，并在子结点中进行处理
		byte_data := data[char_index]
		bval_data = GET_BIT(byte_data, bit_offset)
		if bit_last_path == 0 {
			if bval_data == 0 && this.left != nil {
				return key_pos + int(this.bit_len)
			} else if bval_data == 1 && this.right != nil {
				return key_pos + int(this.bit_len)
			}
		}

		//为剩余的key创建子结点
		nd_data = NewRaxNode()
		nd_data.left = nil
		nd_data.right = nil
		nd_data.val = val
		nd_data.bit_len = bit_last_data
		nd_data.bit_val = make([]byte, len(data[char_index:]))
		copy(nd_data.bit_val, data[char_index:])

		//若bit_offset!=0，说明不是完整字节，
		//将字节分裂，并将字节中的非公共部分分离出来,保存到子结点中
		if bit_offset != 0 {
			byte_tmp := CLEAR_BITS_LOW(byte_data, bit_offset)
			nd_data.bit_val[0] = byte_tmp
		}
	}

	//path的数据有剩余
	//创建子节点：nd_path结点
	//并将数据分开，公共部分保存this结点，其他保存到nd_path结点
	var nd_path *RaxNode = nil
	var bval_path byte
	if bit_last_path > 0 {
		byte_path := this.bit_val[char_index]
		bval_path = GET_BIT(byte_path, bit_offset)

		//为剩余的path创建子结点
		nd_path = NewRaxNode()
		nd_path.left = this.left
		nd_path.right = this.right
		nd_path.val = this.val
		nd_path.bit_len = bit_last_path
		nd_path.bit_val = make([]byte, len(this.bit_val[char_index:]))
		copy(nd_path.bit_val, this.bit_val[char_index:])

		//将byte_path字节中的非公共部分分离出来,保存到子结点中
		if bit_offset != 0 {
			byte_tmp := CLEAR_BITS_LOW(byte_path, bit_offset)
			nd_path.bit_val[0] = byte_tmp
		}

		//修改当前结点，作为nd_path结点、nd_data结点的父结点
		//多申请一个子节，用于存储可能出现的不完整字节
		bit_val_old := this.bit_val
		this.left = nil
		this.right = nil
		this.val = nil
		this.bit_len = this.bit_len - bit_last_path //=bpos - (key_pos % 8)
		this.bit_val = make([]byte, len(bit_val_old[0:char_index])+1)
		copy(this.bit_val, bit_val_old[0:char_index])
		this.bit_val = this.bit_val[0:len(this.bit_val)-1]

		//将byte_path字节中的公共部分分离出来,保存到父结点
		if bit_offset != 0 {
			byte_tmp := CLEAR_BITS_HIGH(byte_path, 8-bit_offset)
			this.bit_val = append(this.bit_val, byte_tmp)
		}
	}

	//若path包含key，则将val赋值给this结点
	if bit_last_data == 0 {
		this.val = val
	}
	if nd_data != nil {
		if bval_data == 0 {
			this.left  = nd_data
		} else {
			this.right = nd_data
		}
	}
	if nd_path != nil {
		if bval_path == 0 {
			this.left  = nd_path
		} else {
			this.right = nd_path
		}
	}
	return len(key) * 8
}

//添加元素
func (this *Radix)Set(key string, val interface{}) {
	data := []byte(key)
	root := this.root[data[0]]
	if root == nil {
		//没有根节点，则创建一个根节点
		root = NewRaxNode()
		root.val = val
		root.bit_val = make([]byte, len(data))
		copy(root.bit_val, data)
		root.bit_len = len(data) * 8
		this.root[data[0]] = root
		return
	} else if root.val == nil &&  root.left == nil && root.right == nil {
		//只有一个根节点，并且是一个空的根节点，则直接赋值
		root.val = val
		root.bit_val = make([]byte, len(data))
		copy(root.bit_val, data)
		root.bit_len = len(data) * 8
		this.root[data[0]] = root
		return
	}

	cur := root
	blen := len(data) * 8
	for bpos := 0; bpos < blen && cur != nil; {
		bpos = cur.pathSplit(data, bpos, val)
		if bpos >= blen {
			return
		}

		ci := bpos / 8
		co := bpos % 8
		byte_data := data[ci]
		bit_pos := GET_BIT(byte_data, co)
		if bit_pos == 0 {
			cur = cur.left
		} else {
			cur = cur.right
		}
	}
}

//将当前结点的子结点进行合并
//若当前结点只有一个子结点，并且当前结点是空结点，才可以进行合并操作
func (this *RaxNode)pathMerge(bpos int) bool {
	//若当前结点存在值，则不能合并
	if this.val != nil {
		return false
	}

	//若当前结点有2个子结点，则不能合并
	if this.left != nil && this.right != nil {
		return false
	}

	//若当前结点没有子结点，则不能合并
	if this.left != nil && this.right != nil {
		return false
	}

	//获取当前结点的子结点
	child := this.left
	if this.right != nil {
		child = this.right
	}

	//判断当前结点最后一个字节是否是完整的字节
	//若不是完整字节，需要与子结点的第一个字节进行合并
	if bpos % 8 != 0 {
		char_len := len(this.bit_val)
		char_last := this.bit_val[char_len-1]
		char_0000 := child.bit_val[0]
		child.bit_val = child.bit_val[1:]
		this.bit_val[char_len-1] = char_last | char_0000
	}

	//合并当前结点以及子结点
	this.val = child.val
	this.bit_val = append(this.bit_val, child.bit_val...)
	this.bit_len += child.bit_len
	this.left = child.left
	this.right = child.right

	return true
}

//删除元素
//只能删除叶子结点，不能删除根节点
//删除叶子结点后，若parent结点只有一个子结点，则将parent结点与子结点合并
func (this *Radix)Delete(key string) {
	data := []byte(key)
	blen := len(data) * 8
	cur := this.root[data[0]]
	var parent *RaxNode = nil
	for bpos := 0; bpos < blen && cur != nil; {
		flag, part_end := cur.pathCompare(data, bpos)
		if flag == false {
			return
		}

		bpos = part_end
		if bpos >= blen {
			//将当前结点修改为空结点
			//若parent是根节点，不能删除
			cur.val = nil
			if parent == nil {
				return
			}

			//当前结点是叶子结点，先将当前结点删除，并将当前结点指向父结点
			if cur.left == nil && cur.right == nil {
				if parent.left == cur {
					parent.left = nil
				} else if parent.right == cur {
					parent.right = nil
				}
				bpos -= int(cur.bit_len)
				cur = parent
			}

			//尝试将当前结点与当前结点的子节点进行合并
			cur.pathMerge(bpos)
			return
		}

		ci := bpos / 8
		co := bpos % 8
		byte_data := data[ci]
		bit_pos := GET_BIT(byte_data, co)
		if bit_pos == 0 {
			parent =cur
			cur = cur.left
		} else {
			parent =cur
			cur = cur.right
		}
	}
}

//查找元素
func (this *Radix)Get(key string) interface{}{
	data := []byte(key)
	blen := len(data) * 8
	cur := this.root[data[0]]
	for bpos := 0; bpos < blen && cur != nil; {
		flag, part_end := cur.pathCompare(data, bpos)
		if flag == false {
			return nil
		}

		bpos = part_end
		if bpos >= blen {
			return cur.val
		}

		ci := bpos / 8
		co := bpos % 8
		byte_data := data[ci]
		bit_pos := GET_BIT(byte_data, co)
		if bit_pos == 0 {
			cur = cur.left
		} else {
			cur = cur.right
		}
	}
	return nil
}

//递归获取数据，用于调试
func (this *Radix)getItems(cur *RaxNode, bpos int, key []byte, items []StringKV) []StringKV {
	//备份key数据
	key_len := len(key)
	var key_last byte
	if key_len > 0 {
		key_last = key[key_len-1]
	}

	//合并key数据
	if bpos % 8 != 0 {
		key = key[0:key_len-1]
		key = append(key, key_last | cur.bit_val[0])
		key = append(key, cur.bit_val[1:]...)
	} else {
		key = append(key, cur.bit_val...)
	}
	bpos += int(cur.bit_len)

	//将value以及可以加入结果集
	if cur.val != nil {
		item := StringKV{string(key), cur.val}
		items = append(items, item)
	}

	if cur.left != nil {
		items = this.getItems(cur.left, bpos, key, items)
	}
	if cur.right != nil {
		items = this.getItems(cur.right, bpos, key, items)
	}

	//恢复key数据
	key = key[0:key_len]
	if key_len > 0 {
		key[key_len-1] = key_last
	}
	return items
}

//获取数据，用于调试
func (this *Radix)GetItems() []StringKV {
	items := make([]StringKV, 0)
	key := make([]byte, 0)
	for i:=0; i < 255; i++ {
		cur := this.root[i]
		if cur == nil {
			continue
		}
		items = this.getItems(cur, 0, key, items)
	}
	return items
}

////打印结点信息，用于调试
func (this *RaxNode)GetNodeInfo(bbeg int) string {
	buff := new(bytes.Buffer)

	bend := bbeg + int(this.bit_len)
	//起始和终止字节的位置
	cbeg := bbeg / 8; cend := bend / 8
	//起始和终止字节的偏移量
	obeg := bbeg % 8; oend := bend % 8
	for bb := bbeg; bb < bend; {
		//获取两个数组的当前字节位置
		dci := bb / 8
		nci := dci - cbeg
		byte_node := this.bit_val[nci]

		//获取数据的当前字节以及循环步长
		step := 8
		if nci == 0 && obeg > 0 {
			step = 8-obeg
		}
		if dci == cend && oend > 0 {
			step = oend
		}
		if cbeg == cend {
			step = int(this.bit_len)
		}

		if step != 8 {
			buff.WriteString(fmt.Sprintf("(%08b:%d)", byte_node, byte_node))
		} else {
			buff.WriteByte(byte_node)
		}
		bb += step
	}

	if this.val != nil {
		buff.WriteString(fmt.Sprintf("=%v", this.val))
	}

	return buff.String()
}

//递归打印结点信息，用于调试
func (this *Radix)getNodesInfo(cur *RaxNode, pos int, data map[string]interface{}) {
	data["info"] = cur.GetNodeInfo(pos)
	pos += int(cur.bit_len)

	if cur.left != nil {
		tmp := make(map[string]interface{})
		data["left"] = tmp
		this.getNodesInfo(cur.left, pos, tmp)
	}

	if cur.right != nil {
		tmp := make(map[string]interface{})
		data["right"] = tmp
		this.getNodesInfo(cur.right, pos, tmp)
	}
}

//打印结点信息，用于调试
func (this *Radix)GetNodesInfo(cc byte) string {
	cur := this.root[cc]
	data_root := make(map[string]interface{})
	this.getNodesInfo(cur, 0, data_root)
	ret, _ := json.MarshalIndent(data_root, "", "    ")
	return string(ret)
}
