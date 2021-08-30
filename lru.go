package lru_go

type LRU struct {
	entries []entry
	head    uint
	tail    uint
}

/// TODO: generics
type entry struct {
	val  interface{}
	prev uint
	next uint
}

// New a fixed size of LRU Cache by cap
func New(cap uint) *LRU {
	return &LRU{
		make([]entry, 0, cap),
		0,
		0,
	}
}

// Insert a value into the LRU Cache; Return the dropped one or nil.
func (lru *LRU) Insert(val interface{}) (ret interface{}) {
	var idx uint
	if lru.IsFull() {
		idx = lru.popBack()
		lru.entries[idx].val, ret = val, lru.entries[idx].val
	} else {
		lru.entries = append(lru.entries, entry{
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

// Find the specific value by the user's supplied Judge function
func (lru *LRU) Find(judge Judge) (ret interface{}) {
	if lru.touch(judge) {
		return lru.entries[lru.head].val
	}
	return nil
}

// IsEmpty return if the LRU Cache is empty
func (lru *LRU) IsEmpty() bool {
	return len(lru.entries) == 0
}

// IsFull return if the LRU Cache is full
func (lru *LRU) IsFull() bool {
	return len(lru.entries) == cap(lru.entries)
}

// Len return counts of elements in this LRU Cache
func (lru *LRU) Len() int {
	return len(lru.entries)
}

// Cap return capacity of this LRU Cache
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

// touch the first element int the cache that match the predicate and mark it as most-recently-used
func (lru *LRU) touch(judge Judge) bool {
	iter := lru.Iterator().(*iterator)
	for iter.HasNext() {
		i := iter.getNext()
		if judge(lru.entries[i].val) {
			lru.indexTouch(i)
			return true
		}
	}
	return false
}

// indexTouch update the given index to the first element
func (lru *LRU) indexTouch(idx uint) {
	if idx != lru.head {
		lru.indexRemove(idx)
		lru.indexFront(idx)
	}
}

// indexRemove remove the given index from the double-order list
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

// indexFront set the given index to the first element
func (lru *LRU) indexFront(idx uint) {
	if len(lru.entries) == 1 {
		lru.tail = idx
	} else {
		lru.entries[idx].next = lru.head
		lru.entries[lru.head].prev = idx
	}
	lru.head = idx
}

type (
	Iterator interface {
		HasNext() bool
		GetNext() interface{}
	}
	iterator struct {
		cur   uint
		isEnd bool
		lru   *LRU
	}
)

// Iterator return an iterator on this LRU Cache
func (lru *LRU) Iterator() Iterator {
	return &iterator{
		lru.head,
		false,
		lru,
	}
}

// HasNext return if this iterator has touch the tail
func (iter *iterator) HasNext() bool {
	return !iter.isEnd
}

// getNext return the next element's index
func (iter *iterator) getNext() uint {
	cur := iter.cur
	if cur == iter.lru.tail {
		iter.isEnd = true
	}
	iter.cur = iter.lru.entries[cur].next
	return cur
}

// GetNext return next element's value in the LRU
func (iter *iterator) GetNext() interface{} {
	return iter.lru.entries[iter.getNext()].val
}

// get specific index's entry and if it exists
func (lru *LRU) get(idx uint) (*entry, bool) {
	iter := lru.Iterator().(*iterator)
	for iter.HasNext() {
		if idx == 0 {
			return &lru.entries[iter.getNext()], true
		}
		idx--
	}
	return nil, false
}

// Get return the element's value of specific index
func (lru *LRU) Get(idx uint) interface{} {
	entry, ok := lru.get(idx)
	if !ok {
		return nil
	}
	return entry.val
}

// ToArray return the array of elements in order
func (lru *LRU) ToArray() (ret []interface{}) {
	ret = make([]interface{}, 0, lru.Len())
	iter := lru.Iterator().(*iterator)
	for iter.HasNext() {
		ret = append(ret, lru.entries[iter.getNext()].val)
	}
	return ret
}
