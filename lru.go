package lru

type Iterator interface {
	HasNext() bool
	GetNext() uint
}

type LRU struct {
	entries []Entry
	head    uint
	tail    uint
}

type Entry struct {
	val  interface{}
	prev uint
	next uint
}

func New(cap uint) *LRU {
	return &LRU{
		make([]Entry, 0, cap),
		0,
		0,
	}
}

func (lru *LRU) Insert(val interface{}) (ret interface{}) {
	var idx uint
	if lru.IsFull() {
		idx = lru.popBack()
		lru.entries[idx].val, ret = val, lru.entries[idx].val
	} else {
		lru.entries = append(lru.entries, Entry{
			val,
			0,
			0,
		})
		idx = uint(len(lru.entries)) - 1
	}
	lru.indexFront(idx)
	return
}

type Judge func(interface{}) bool

func (lru *LRU) Find(judge Judge) (ret interface{}) {
	if lru.touch(judge) {
		return lru.entries[lru.head].val
	}
	return nil
}

func (lru *LRU) IsEmpty() bool {
	return len(lru.entries) == 0
}

func (lru *LRU) IsFull() bool {
	return len(lru.entries) == cap(lru.entries)
}

func (lru *LRU) Len() int {
	return len(lru.entries)
}

func (lru *LRU) Cap() int {
	return cap(lru.entries)
}

// popBack pop the last entry from the LRU cache
func (lru *LRU) popBack() uint {
	oldTail := lru.tail
	newTail := lru.entries[oldTail].prev
	lru.tail = newTail
	return oldTail
}

// touch touch the first elements int the cache that match the predicate and mark it as most-recently-used
func (lru *LRU) touch(judge Judge) bool {
	iter := lru.Itor()
	for iter.HasNext() {
		i := iter.GetNext()
		if judge(lru.entries[i].val) {
			lru.indexTouch(i)
			return true
		}
	}
	return false
}

// indexTouch
func (lru *LRU) indexTouch(idx uint) {
	if idx != lru.head {
		lru.indexRemove(idx)
		lru.indexFront(idx)
	}
}

// indexRemove
func (lru *LRU) indexRemove(idx uint) {
	prev := lru.entries[idx].prev
	next := lru.entries[idx].next

	if idx == lru.head {
		lru.head = next
	} else {
		lru.entries[prev].next = next
	}

	if idx == lru.tail {
		lru.tail = prev
	} else {
		lru.entries[next].prev = prev
	}
}

// indexFront update the head of the LRU to the given index
func (lru *LRU) indexFront(idx uint) {
	if len(lru.entries) == 1 {
		lru.tail = idx
	} else {
		lru.entries[idx].next = lru.head
		lru.entries[lru.head].prev = idx
	}
	lru.head = idx
}

type LRUIterator struct {
	cur   uint
	isEnd bool
	lru   *LRU
}

func (lru *LRU) Itor() *LRUIterator {
	return &LRUIterator{
		lru.head,
		false,
		lru,
	}
}

func (iter *LRUIterator) HasNext() bool {
	return !iter.isEnd
}

func (iter *LRUIterator) GetNext() uint {
	cur := iter.cur
	if cur == iter.lru.tail {
		iter.isEnd = true
	}
	iter.cur = iter.lru.entries[cur].next
	return cur
}
