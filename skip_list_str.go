package gdata

import (
	"math/rand"
	"sync"
	"time"
)

type SkipNodeString struct {
	key		string
	value	interface{}
	next 	[]*SkipNodeString
}

type SkipListString struct {
	SkipNodeString
	mutex  	sync.RWMutex
	update 	[]*SkipNodeString
	rand    *rand.Rand
	maxl  	int
	skip  	int
	level  	int
	length 	int32
}

func NewSkipListString(skip ...int) *SkipListString {
	list := &SkipListString{}
	list.maxl = 32
	list.skip = 4
	list.level = 0
	list.length = 0
	list.SkipNodeString.next = make([]*SkipNodeString, list.maxl)
	list.update = make([]*SkipNodeString, list.maxl)
	list.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	if len(skip) == 1 && skip[0] > 1 {
		list.skip = skip[0]
	}
	return list
}

func (list *SkipListString) Get(key string) interface{} {
	list.mutex.Lock()
	defer list.mutex.Unlock()

	var prev = &list.SkipNodeString
	var next *SkipNodeString
	for i := list.level-1; i >= 0; i-- {
		next = prev.next[i]
		for next != nil && next.key < key {
			prev = next
			next = prev.next[i]
		}
	}

	if next != nil && next.key == key {
		return next.value
	} else {
		return nil
	}
}

func (list *SkipListString) Set(key string, val interface{}) {
	list.mutex.Lock()
	defer list.mutex.Unlock()

	var prev = &list.SkipNodeString
	var next *SkipNodeString
	for i := list.level-1; i >= 0; i-- {
		next = prev.next[i]
		for next != nil && next.key < key {
			prev = next
			next = prev.next[i]
		}
		list.update[i] = prev
	}

	//如果key已经存在
	if next != nil && next.key == key {
		next.value = val
		return
	}

	//随机生成新结点的层数
	level := list.randomLevel();
	if level > list.level {
		level = list.level + 1;
		list.level = level
		list.update[list.level-1] = &list.SkipNodeString
	}

	//申请新的结点
	node:= &SkipNodeString{}
	node.key = key
	node.value = val
	node.next = make([]*SkipNodeString, level)

	//调整next指向
	for i := 0; i < level; i++ {
		node.next[i] = list.update[i].next[i]
		list.update[i].next[i] = node
	}

	list.length++
}

func (list *SkipListString) Remove(key string) interface{} {
	list.mutex.Lock()
	defer list.mutex.Unlock()

	var prev = &list.SkipNodeString
	var next *SkipNodeString
	for i := list.level-1; i >= 0; i-- {
		next = prev.next[i]
		for next != nil && next.key < key {
			prev = next
			next = prev.next[i]
		}
		list.update[i] = prev
	}

	//结点不存在
	node := next
	if next == nil || next.key != key {
		return nil
	}

	//调整next指向
	for i, v := range node.next {
		if list.update[i].next[i] == node {
			list.update[i].next[i] = v
			if list.SkipNodeString.next[i] == nil {
				list.level -= 1
			}
		}
		list.update[i] = nil
	}

	list.length--
	return node.value
}

func (list *SkipListString) GetLength() int32 {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	return list.length
}

func (list *SkipListString) randomLevel() int {
	i := 1
	for ; i < list.maxl; i++ {
		if list.rand.Int31() % int32(list.skip) != 0 {
			break
		}
	}
	return i
}

