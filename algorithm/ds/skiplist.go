package ds

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

/*
跳表(skiplist)是一种概率性的数据结构，支持插入和查询复杂度均为O(logn),比起红黑树有着近似的插入及查询复杂度，
实现上比红黑树相对简单。redis的有序集合在满足一定条件后就会由跳表实现。

skiplist是由多层的有序链表组成，每层的节点在上一层出现的概率为p(通常为1/2或1/4)，均匀尺度上每个节点出现的层数
为1/(1-p)= lim (1+q+q*q+q*q*q+...) with n->正无穷 = 1/(1-p),最高的节点(通常是头节点)出现在所有层,跳表包含
log(1/p)n层(why?: n*p^k=1 -> k = log(1/p)n)。总体来说搜索复杂度为(1/p)*log(1/p)n=O(logn)。

head -> 1 -----------------------------------------------> nil
head -> -------------> 4 ------> 6 ----------------------> nil
head -> 1 ------> 3 -> 4 ------> 6 -----------> 9 -------> nil
head -> 1 -> 2 -> 3 -> 4 -> 5 -> 6 -> 7 -> 8 -> 9 -> 10 -> nil

#C语言版本实现：https://www.epaperpress.com/sortsearch/txt/skl.txt
*/

type SkipList struct {
	Head        *SkipListNode
	Probability float32
	MaxLevel    int
	LevelUsed   int
}

type SkipListNode struct {
	Key     int
	Val     interface{}
	Forward []*SkipListNode
	Level   int
}

const (
	DefaultMaxLevel    = 16
	DefaultProbability = 1.0 / 2 //Probability can be used to make a trade-off between storage cost and search cost,usually 0.5 or 0.25
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

func NewNode(k int, v interface{}, crlevel int, maxlevel int) *SkipListNode {
	emptyfwd := make([]*SkipListNode, maxlevel)
	for i := 0; i < maxlevel; i++ {
		emptyfwd[i] = nil
	}

	return &SkipListNode{
		Key:     k,
		Val:     v,
		Forward: emptyfwd,
		Level:   crlevel,
	}
}

func NewSkipList() *SkipList {
	rand.Seed(int64(time.Now().Nanosecond()))

	return &SkipList{
		Head:        NewNode(0, "headnode", 1, DefaultMaxLevel),
		Probability: DefaultProbability,
		MaxLevel:    DefaultMaxLevel,
		LevelUsed:   1,
	}
}

func (sk *SkipList) randomf() float32 {
	return rand.Float32()
}

func (sk *SkipList) RandomLevel() int {
	level := 1
	for sk.randomf() < sk.Probability && level < sk.MaxLevel {
		level++
	}
	return level
}

// Search will try to find data with key==dst and return its val.
func (sk *SkipList) Search(dst int) (interface{}, error) {
	cur := sk.Head

	for i := sk.LevelUsed - 1; i >= 0; i-- {
		for cur.Forward[i] != nil && cur.Forward[i].Key < dst {
			cur = cur.Forward[i]
		}
	}

	//set to bottom level(all data stored in this level)
	cur = cur.Forward[0]

	if cur != nil && cur.Key == dst {
		return cur.Val, nil
	}

	return nil, ErrKeyNotFound
}

// Insert will put key&val into right place, if key already exists then update value only.
func (sk *SkipList) Insert(key int, val interface{}) {
	updatelist := make([]*SkipListNode, sk.MaxLevel)

	cur := sk.Head
	for i := sk.LevelUsed - 1; i >= 0; i-- {
		for cur.Forward[i] != nil && cur.Forward[i].Key < key {
			cur = cur.Forward[i]
		}
		updatelist[i] = cur
	}

	cur = cur.Forward[0]

	if cur != nil && cur.Key == key {
		//key already exists -> update to newer value
		cur.Val = val
	} else {
		newlevel := sk.RandomLevel()
		if newlevel > sk.LevelUsed {
			for exlvl := sk.LevelUsed; exlvl < newlevel; exlvl++ {
				updatelist[exlvl] = sk.Head
			}
			sk.LevelUsed = newlevel
			sk.Head.Level = newlevel //head exists in every level
		}

		newnode := NewNode(key, val, newlevel, sk.MaxLevel)
		for i := 0; i < newlevel; i++ {
			newnode.Forward[i] = updatelist[i].Forward[i]
			updatelist[i].Forward[i] = newnode
		}
	}
}

// Delete will
func (sk *SkipList) Delete(dst int) error {
	updatelist := make([]*SkipListNode, sk.MaxLevel)

	cur := sk.Head
	for i := sk.Head.Level - 1; i >= 0; i-- {
		for cur.Forward[i] != nil && cur.Forward[i].Key < dst {
			cur = cur.Forward[i]
		}
		updatelist[i] = cur
	}

	//move to location which should be deleted
	cur = cur.Forward[0]

	if cur == nil || cur.Key < dst {
		return ErrKeyNotFound
	}

	for i := 0; i < cur.Level; i++ {
		if updatelist[i].Forward[i] != nil && updatelist[i].Forward[i].Key != cur.Key {
			break
		}
		updatelist[i].Forward[i] = cur.Forward[i]
	}

	//need to adjust levelused

	for cur.Level > 1 && sk.Head.Forward[cur.Level] == nil {
		cur.Level--
	}

	cur = nil
	return nil
}

func (sk *SkipList) PrintSkipList() {
	if sk == nil {
		fmt.Println("nil")
	}

	cur := sk.Head

	for cur != nil {
		fmt.Printf("node:[%d], val:%v, level:%d, forward:%v\n", cur.Key, cur.Val, cur.Level, cur.Forward)
		cur = cur.Forward[0]
	}
}
