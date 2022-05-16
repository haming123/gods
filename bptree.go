package gdata

import (
	"sync"
)

type BPItem struct {
	Key 	int64
	Val 	interface{}
}

type StringKV struct {
	Key 	string
	Val 	interface{}
}

type BPNode struct {
	MaxKey	int64
	Nodes 	[]*BPNode
	Items 	[]BPItem
	Next 	*BPNode
}

func (node *BPNode) findItem(key int64) int {
	num := len(node.Items)
	for i:=0; i < num; i++ {
		if node.Items[i].Key > key {
			return -1
		} else if node.Items[i].Key == key {
			return i
		}
	}
	return -1
}

func (node *BPNode) setValue(key int64, value interface{}) {
	item := BPItem{key, value}
	num := len(node.Items)
	if num < 1 {
		node.Items = append(node.Items, item)
		node.MaxKey = item.Key
		return
	} else if key < node.Items[0].Key {
		node.Items = append([]BPItem{item}, node.Items...)
		return
	} else if key > node.Items[num-1].Key {
		node.Items = append(node.Items, item)
		node.MaxKey = item.Key
		return
	}

	for i:=0; i < num; i++ {
		if node.Items[i].Key > key {
			node.Items = append(node.Items, BPItem{})
			copy(node.Items[i+1:], node.Items[i:])
			node.Items[i] = item
			return
		} else if node.Items[i].Key == key {
			node.Items[i] = item
			return
		}
	}
}

func (node *BPNode) addChild(child *BPNode) {
	num := len(node.Nodes)
	if num < 1 {
		node.Nodes = append(node.Nodes, child)
		node.MaxKey = child.MaxKey
		return
	} else if child.MaxKey < node.Nodes[0].MaxKey {
		node.Nodes = append([]*BPNode{child}, node.Nodes...)
		return
	} else if child.MaxKey > node.Nodes[num-1].MaxKey {
		node.Nodes = append(node.Nodes, child)
		node.MaxKey = child.MaxKey
		return
	}

	for i:=0; i < num; i++ {
		if node.Nodes[i].MaxKey > child.MaxKey {
			node.Nodes = append(node.Nodes, nil)
			copy(node.Nodes[i+1:], node.Nodes[i:])
			node.Nodes[i] = child
			return
		}
	}
}

func (node *BPNode) deleteItem(key int64) bool {
	num := len(node.Items)
	for i:=0; i < num; i++ {
		if node.Items[i].Key > key {
			return false
		} else if node.Items[i].Key == key {
			copy(node.Items[i:], node.Items[i+1:])
			node.Items = node.Items[0:len(node.Items)-1]
			node.MaxKey = node.Items[len(node.Items)-1].Key
			return true
		}
	}
	return false
}

func (node *BPNode) deleteChild(child *BPNode) bool {
	num := len(node.Nodes)
	for i:=0; i < num; i++ {
		if node.Nodes[i] == child {
			copy(node.Nodes[i:], node.Nodes[i+1:])
			node.Nodes = node.Nodes[0:len(node.Nodes)-1]
			node.MaxKey = node.Nodes[len(node.Nodes)-1].MaxKey
			return true
		}
	}
	return false
}

type BPTree struct {
	mutex  	sync.RWMutex
	ktype 	int
	root  	*BPNode
	width 	int
	halfw 	int
}

func NewBPTree(width int) *BPTree {
	if width < 3 {
		width = 3
	}

	var bt = &BPTree{}
	bt.root = NewLeafNode(width)
	bt.width = width
	bt.halfw = (bt.width + 1) / 2
	return bt
}

//申请width+1是因为插入时可能暂时出现节点key大于申请width的情况,待后期再分裂处理
func NewLeafNode(width int) *BPNode {
	var node = &BPNode{}
	node.Items = make([]BPItem, width+1)
	node.Items = node.Items[0:0]
	return node
}

//申请width+1是因为插入时可能暂时出现节点key大于申请width的情况,待后期再分裂处理
func NewIndexNode(width int) *BPNode {
	var node = &BPNode{}
	node.Nodes = make([]*BPNode, width+1)
	node.Nodes = node.Nodes[0:0]
	return node
}

func (t *BPTree) Get(key int64) interface{} {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	node := t.root
	for i := 0; i < len(node.Nodes); i++ {
		if key <= node.Nodes[i].MaxKey {
			node = node.Nodes[i]
			i = 0
		}
	}

	//没有到达叶子结点
	if len(node.Nodes) > 0 {
		return nil
	}

	for i := 0; i < len(node.Items); i++ {
		if node.Items[i].Key == key {
			return node.Items[i].Val
		}
	}
	return nil
}

func (t *BPTree) getData(node *BPNode) map[int64]interface{} {
	data := make(map[int64]interface{})
	for {
		if len(node.Nodes) > 0 {
			for i := 0; i < len(node.Nodes); i++ {
				data[node.Nodes[i].MaxKey] = t.getData(node.Nodes[i])
			}
			break
		} else {
			for i := 0; i < len(node.Items); i++ {
				data[node.Items[i].Key] = node.Items[i].Val
			}
			break
		}
	}
	return data
}

func (t *BPTree) GetData() map[int64]interface{} {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	return t.getData(t.root)
}

func (t *BPTree) splitNode(node *BPNode) *BPNode {
	if len(node.Nodes) > t.width {
		//创建新结点
		halfw := t.width / 2 + 1
		node2 := NewIndexNode(t.width)
		node2.Nodes = append(node2.Nodes, node.Nodes[halfw : len(node.Nodes)]...)
		node2.MaxKey = node2.Nodes[len(node2.Nodes)-1].MaxKey

		//修改原结点数据
		node.Nodes = node.Nodes[0:halfw]
		node.MaxKey = node.Nodes[len(node.Nodes)-1].MaxKey

		return node2
	} else if len(node.Items) > t.width {
		//创建新结点
		halfw := t.width / 2 + 1
		node2 := NewLeafNode(t.width)
		node2.Items = append(node2.Items, node.Items[halfw: len(node.Items)]...)
		node2.MaxKey = node2.Items[len(node2.Items)-1].Key

		//修改原结点数据
		node.Next = node2
		node.Items = node.Items[0:halfw]
		node.MaxKey = node.Items[len(node.Items)-1].Key

		return node2
	}

	return nil
}

func (t *BPTree) setValue(parent *BPNode, node *BPNode, key int64, value interface{}) {
	for i:=0; i < len(node.Nodes); i++ {
		if key <= node.Nodes[i].MaxKey || i== len(node.Nodes)-1 {
			t.setValue(node, node.Nodes[i], key, value)
			break
		}
	}

	//叶子结点，添加数据
	if len(node.Nodes) < 1 {
		node.setValue(key, value)
	}

	//结点分裂
	node_new := t.splitNode(node)
	if node_new != nil {
		//若父结点不存在，则创建一个父节点
		if parent == nil {
			parent = NewIndexNode(t.width)
			parent.addChild(node)
			t.root = parent
		}
		//添加结点到父亲结点
		parent.addChild(node_new)
	}
}

func (t *BPTree) Set(key int64, value interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.setValue(nil, t.root, key, value)
}

func (t *BPTree) itemMoveOrMerge(parent *BPNode, node *BPNode) {
	//获取兄弟结点
	var node1 *BPNode = nil
	var node2 *BPNode = nil
	for i:=0; i < len(parent.Nodes); i++ {
		if parent.Nodes[i] == node {
			if i < len(parent.Nodes)-1 {
				node2 = parent.Nodes[i+1]
			} else if i > 0 {
				node1 = parent.Nodes[i-1]
			}
			break
		}
	}

	//将左侧结点的记录移动到删除结点
	if node1 != nil && len(node1.Items) > t.halfw {
		item := node1.Items[len(node1.Items)-1]
		node1.Items = node1.Items[0:len(node1.Items)-1]
		node1.MaxKey = node1.Items[len(node1.Items)-1].Key
		node.Items = append([]BPItem{item}, node.Items...)
		return
	}

	//将右侧结点的记录移动到删除结点
	if node2 != nil && len(node2.Items) > t.halfw {
		item := node2.Items[0]
		node2.Items = node1.Items[1:]
		node.Items = append(node.Items, item)
		node.MaxKey = node.Items[len(node.Items)-1].Key
		return
	}

	//与左侧结点进行合并
	if node1 != nil && len(node1.Items) + len(node.Items) <= t.width {
		node1.Items = append(node1.Items, node.Items...)
		node1.Next = node.Next
		node1.MaxKey = node1.Items[len(node1.Items)-1].Key
		parent.deleteChild(node)
		return
	}

	//与右侧结点进行合并
	if node2 != nil && len(node2.Items) + len(node.Items) <= t.width {
		node.Items = append(node.Items, node2.Items...)
		node.Next = node2.Next
		node.MaxKey = node.Items[len(node.Items)-1].Key
		parent.deleteChild(node2)
		return
	}
}

func (t *BPTree) childMoveOrMerge(parent *BPNode, node *BPNode) {
	if parent == nil {
		return
	}

	//获取兄弟结点
	var node1 *BPNode = nil
	var node2 *BPNode = nil
	for i:=0; i < len(parent.Nodes); i++ {
		if parent.Nodes[i] == node {
			if i < len(parent.Nodes)-1 {
				node2 = parent.Nodes[i+1]
			} else if i > 0 {
				node1 = parent.Nodes[i-1]
			}
			break
		}
	}

	//将左侧结点的子结点移动到删除结点
	if node1 != nil && len(node1.Nodes) > t.halfw {
		item := node1.Nodes[len(node1.Nodes)-1]
		node1.Nodes = node1.Nodes[0:len(node1.Nodes)-1]
		node.Nodes = append([]*BPNode{item}, node.Nodes...)
		return
	}

	//将右侧结点的子结点移动到删除结点
	if node2 != nil && len(node2.Nodes) > t.halfw {
		item := node2.Nodes[0]
		node2.Nodes = node1.Nodes[1:]
		node.Nodes = append(node.Nodes, item)
		return
	}

	if node1 != nil && len(node1.Nodes) + len(node.Nodes) <= t.width {
		node1.Nodes = append(node1.Nodes, node.Nodes...)
		parent.deleteChild(node)
		return
	}

	if node2 != nil && len(node2.Nodes) + len(node.Nodes) <= t.width {
		node.Nodes = append(node.Nodes, node2.Nodes...)
		parent.deleteChild(node2)
		return
	}
}

func (t *BPTree) deleteItem(parent *BPNode, node *BPNode, key int64) {
	for i:=0; i < len(node.Nodes); i++ {
		if key <= node.Nodes[i].MaxKey {
			t.deleteItem(node, node.Nodes[i], key)
			break
		}
	}

	if  len(node.Nodes) < 1 {
		//删除记录后若结点的子项<m/2，则从兄弟结点移动记录，或者合并结点
		node.deleteItem(key)
		if len(node.Items) < t.halfw {
			t.itemMoveOrMerge(parent, node)
		}
	} else {
		//若结点的子项<m/2，则从兄弟结点移动记录，或者合并结点
		node.MaxKey = node.Nodes[len(node.Nodes)-1].MaxKey
		if len(node.Nodes) < t.halfw {
			t.childMoveOrMerge(parent, node)
		}
	}
}

func (t *BPTree) Remove(key int64) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.deleteItem(nil, t.root, key)
}
