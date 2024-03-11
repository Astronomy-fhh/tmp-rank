package store

import (
	"sort"
	"sync"
)

type ArraySortStore[T Element] struct {
	sync.RWMutex
	sorted bool
	dict   map[int64]*ArraySortWrapper[T]
	sort   []*ArraySortWrapper[T]
}

func NewArraySortStore[T Element]() *ArraySortStore[T] {
	as := &ArraySortStore[T]{
		dict: make(map[int64]*ArraySortWrapper[T]),
	}
	return as
}

type ArraySortWrapper[T Element] struct {
	value T
	idx   int
}

// Add 活动期间添加/更新玩家分数
func (p *ArraySortStore[T]) Add(psc T) {
	if p.sorted {
		return // 排序后不允许再插入数据，return或返回错误
	}
	p.Lock()
	defer p.Unlock()
	if p.sorted {
		return // 排序后不允许再插入数据，return或返回错误
	}
	p.dict[psc.GetKey()] = &ArraySortWrapper[T]{
		value: psc,
		idx:   0,
	}
}

// Sort 活动结束后 统一排序
func (p *ArraySortStore[T]) Sort() {
	p.Lock()
	defer p.Unlock()
	if p.sort != nil {
		return // 已经排完序
	}

	p.sort = make([]*ArraySortWrapper[T], 0, len(p.dict))
	for _, v := range p.dict {
		p.sort = append(p.sort, v)
	}

	sort.Slice(p.sort, func(i, j int) bool {
		return p.sort[i].value.Compare(p.sort[j].value) < 0
	})

	for i, wrapper := range p.sort {
		wrapper.idx = i
	}
}

// GetSort 活动结束后 获取玩家分数及周围玩家
func (p *ArraySortStore[T]) GetSort(key int64, rangeNum int) ([]T, int64) {
	p.RLock()
	defer p.RUnlock()

	var rank int64
	ret := make([]T, 0)
	if p.sort == nil {
		return ret, rank
	}
	userScoreWrapper, ok := p.dict[key]
	if !ok {
		return ret, rank
	}
	startIdx := userScoreWrapper.idx - rangeNum
	if startIdx < 0 {
		startIdx = 0
	}
	rank = int64(startIdx)
	endIdx := userScoreWrapper.idx + rangeNum + 1
	if endIdx > len(p.sort) {
		endIdx = len(p.sort)
	}
	for i := startIdx; i < endIdx; i++ {
		ret = append(ret, p.sort[i].value)
	}
	return ret, rank + 1
}
