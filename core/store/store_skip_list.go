package store

import (
	"math/rand"
	"sync"
)

type SkipListStore[T Element] struct {
	skipList *SkipList[T]
}

func NewSkipListStore[T Element]() *SkipListStore[T] {
	sls := SkipListStore[T]{
		skipList: New[T](),
	}
	return &sls
}

func (p *SkipListStore[T]) Add(t T) {
	p.skipList.AddOrUpdate(t)
}

func (p *SkipListStore[T]) GetSort(key int64, rangeNum int) ([]T, int64) {
	return p.skipList.FindRangeNum(key, rangeNum)
}

// 跳表的定义
const MaxLevel = 32
const P = 0.25

type SkipList[T Element] struct {
	sync.RWMutex
	header *Node[T]
	tail   *Node[T]
	level  int
	dict   map[int64]*Node[T]
}

type Level[T Element] struct {
	forward *Node[T]
	span    int64 // 在当前层前一个节点到当前节点的跨度
}

type Node[T Element] struct {
	Value    T
	backward *Node[T]
	level    []Level[T]
}

func New[T Element]() *SkipList[T] {
	sortedSet := SkipList[T]{
		level: 1,
		dict:  make(map[int64]*Node[T]),
	}

	var zeroT T
	sortedSet.header = &Node[T]{
		Value: zeroT,
		level: make([]Level[T], MaxLevel),
	}
	return &sortedSet
}

func (sl *SkipList[T]) createNode(level int, psc T) *Node[T] {
	node := Node[T]{
		Value: psc,
		level: make([]Level[T], level),
	}
	return &node
}

func randomLevel() int {
	level := 1
	for rand.Float64() < P && level < MaxLevel {
		level++
	}
	return level
}

//Level 2     [H] ------------------------------------------> [E] ---------------> ...
//|                                                            |
//Level 1     [H] -------------> [B] ---------> [D] --------> [E] ---------------> ...
//|            |                  |              |             ｜
//Level 0     [H] --->  [A] ---->[B] -> [C] --> [D] --------> [E] -> [F] -> [H] -> ...

func (sl *SkipList[T]) insert(value T) *Node[T] {
	var update [MaxLevel]*Node[T]
	var rank [MaxLevel]int64

	current := sl.header
	// 从上而下依次找到每层的插入位置
	for i := sl.level - 1; i >= 0; i-- {
		if sl.level-1 == i {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		// 确定每层的插入位置
		for current.level[i].forward != nil &&
			(current.level[i].forward.Value.Compare(value) < 0 ||
				(current.level[i].forward.Value.Compare(value) == 0 && current.level[i].forward.Value.GetKey() < value.GetKey())) {
			rank[i] += current.level[i].span
			current = current.level[i].forward
		}
		update[i] = current
	}

	// 随机新节点的level
	level := randomLevel()

	// 如果新节点的level大于当前的level，更新header的level
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			rank[i] = 0
			update[i] = sl.header
			update[i].level[i].span = 0
		}
		sl.level = level
	}

	current = sl.createNode(level, value)

	// 更新每层的指针
	for i := 0; i < level; i++ {
		// 更新当前节点的forward指针
		current.level[i].forward = update[i].level[i].forward
		// 更新前一个节点的forward指针
		update[i].level[i].forward = current
		// 更新span (rank[0] - rank[i])表示0层比i层多跨越了几个节点 所以在原来的基础上去减
		current.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	// 更新未触及到的level的span
	for i := level; i < sl.level; i++ {
		update[i].level[i].span++
	}

	// 维护0层的backward指针
	if update[0] == sl.header {
		current.backward = nil
	} else {
		current.backward = update[0]
	}
	if current.level[0].forward != nil {
		current.level[0].forward.backward = current
	} else {
		sl.tail = current
	}
	return current
}

func (sl *SkipList[T]) deleteNode(x *Node[T]) {
	var update [MaxLevel]*Node[T]

	// 从上而下依次找到每层的插入位置
	current := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for current.level[i].forward != nil &&
			(current.level[i].forward.Value.Compare(x.Value) < 0 || (current.level[i].forward.Value.Compare(x.Value) == 0 && current.level[i].forward.Value.GetKey() < x.Value.GetKey())) {
			current = current.level[i].forward
		}
		update[i] = current
	}

	// 更新每一层的span和指针
	for i := 0; i < sl.level; i++ {
		if update[i].level[i].forward == x {
			update[i].level[i].span += x.level[i].span - 1
			update[i].level[i].forward = x.level[i].forward
		} else if update[i].level[i].span > 0 {
			update[i].level[i].span--
		}
	}

	// 更新0层的前向指针
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x.backward
	} else {
		sl.tail = x.backward
	}

	// 降低level
	for sl.level > 1 && sl.header.level[sl.level-1].forward == nil {
		sl.level--
	}
}

func (sl *SkipList[T]) delete(value T) bool {
	oldNode, found := sl.dict[value.GetKey()]
	if !found {
		return false
	}
	sl.deleteNode(oldNode)
	return true
}

// AddOrUpdate 添加或更新 如果已经存在则更新
func (sl *SkipList[T]) AddOrUpdate(value T) bool {
	sl.Lock()
	defer sl.Unlock()
	oldNode, found := sl.dict[value.GetKey()]
	if found {
		if oldNode.Value.Compare(value) == 0 {
			return false
		}
		sl.deleteNode(oldNode)
	}
	newNode := sl.insert(value)
	sl.dict[value.GetKey()] = newNode
	return true
}

func (sl *SkipList[T]) FindRangeNum(key int64, distance int) ([]T, int64) {
	sl.RLock()
	defer sl.RUnlock()
	res := make([]T, 0)
	node, found := sl.dict[key]
	if !found {
		return res, 0
	}
	rank := make([]int64, sl.level) // rank 为每一层的排名数组

	current := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		if sl.level-1 == i {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}
		for current.level[i].forward != nil &&
			(current.level[i].forward.Value.Compare(node.Value) < 0 || (current.level[i].forward.Value.Compare(node.Value) == 0 && current.level[i].forward.Value.GetKey() < node.Value.GetKey())) {
			rank[i] += current.level[i].span
			current = current.level[i].forward
		}
	}
	finalRank := rank[0]
	//println(key, finalRank, distance)

	// 向前看
	currentNode := node.backward
	for i := 0; i < distance; i++ {
		if currentNode == nil {
			break
		}
		res = append([]T{currentNode.Value}, res...)
		currentNode = currentNode.backward
		finalRank--
	}
	res = append(res, node.Value)

	// 向后看
	currentNode = node.level[0].forward
	for i := 0; i < distance; i++ {
		if currentNode == nil {
			break
		}
		res = append(res, currentNode.Value)
		currentNode = currentNode.level[0].forward
	}

	return res, finalRank + 1
}

func (sl *SkipList[T]) Look() {
	sl.RLock()
	defer sl.RUnlock()
	start := sl.header.level[0].forward
	for start != nil {
		println(start.Value.GetKey())
		start = start.level[0].forward
	}
}

func (sl *SkipList[T]) LookTail() {
	sl.RLock()
	defer sl.RUnlock()
	start := sl.tail
	for start != nil {
		println(start.Value.GetKey())
		start = start.backward
	}
}
