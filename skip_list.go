package gdata

import (
	"math/rand"
	"sync"
)

type SkipNode struct {
	key		interface{}
	value	interface{}
	next 	[]*SkipNode
}

type SkipList struct {
	SkipNode
	fn_cmp Comparable
	mutex  	sync.RWMutex
	update 	[]*SkipNode
	maxl  	int
	skip  	int
	randn  	int
	level  	int
	length 	int32
}

func NewSkipList(fn_cmp Comparable, skip ...int) *SkipList {
	list := &SkipList{}
	list.fn_cmp = fn_cmp
	list.maxl = 32
	list.skip = 4
	list.randn = 0
	list.level = 0
	list.length = 0
	list.SkipNode.next = make([]*SkipNode, list.maxl)
	list.update = make([]*SkipNode, list.maxl)
	if len(skip) == 1 && skip[0] > 1 {
		list.skip = skip[0]
	}
	return list
}

func (list *SkipList) Get(key interface{}) interface{} {
	list.mutex.Lock()
	defer list.mutex.Unlock()

	var ppp = &list.SkipNode
	for i := list.level-1; i >= 0; i-- {
		for ppp.next[i] != nil && list.fn_cmp.Compare(ppp.next[i].key, key) < 0 {
			ppp = ppp.next[i]
		}
	}

	cur := ppp.next[0]
	if cur != nil && cur.key == key {
		return cur.value
	} else {
		return nil
	}
}

func (list *SkipList) Set(key interface{}, val interface{}) {
	list.mutex.Lock()
	defer list.mutex.Unlock()

	var ppp = &list.SkipNode
	for i := list.level-1; i >= 0; i-- {
		for ppp.next[i] != nil && list.fn_cmp.Compare(ppp.next[i].key, key) < 0 {
			ppp = ppp.next[i]
		}
		list.update[i] = ppp
	}

	//如果key已经存在
	cur := ppp.next[0]
	if cur != nil && cur.key == key {
		cur.value = val
		return
	}

	//随机生成新结点的层数
	level := list.randomLevel();
	if level > list.level {
		level = list.level + 1;
		list.level = level
		list.update[list.level-1] = &list.SkipNode
	}

	//申请新的结点
	node:= &SkipNode{}
	node.key = key
	node.value = val
	node.next = make([]*SkipNode, level)

	//调整next指向
	for i := 0; i < level; i++ {
		node.next[i] = list.update[i].next[i]
		list.update[i].next[i] = node
		list.update[i] = nil
	}

	list.length++
}

func (list *SkipList) Remove(key interface{}) bool {
	list.mutex.Lock()
	defer list.mutex.Unlock()

	var ppp = &list.SkipNode
	for i := list.level-1; i >= 0; i-- {
		for ppp.next[i] != nil && list.fn_cmp.Compare(ppp.next[i].key, key) < 0 {
			ppp = ppp.next[i]
		}
		list.update[i] = ppp
	}

	//结点不存在
	cur := ppp.next[0]
	if cur == nil || cur.key != key {
		return false
	}

	//调整next指向
	for i, v := range cur.next {
		if list.update[i].next[i] == cur {
			list.update[i].next[i] = v
			if list.SkipNode.next[i] == nil {
				list.level -= 1
			}
		}
		list.update[i] = nil
	}

	list.length--
	return true
}

func (list *SkipList) GetLength() int32 {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	return list.length
}

func (list *SkipList) randomLevel() int {
	i := 1
	for ; i < list.maxl; i++ {
		if rand.Int() % list.skip != 0 {
			break
		}
	}
	return i
}

