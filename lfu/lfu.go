package lfu

import (
	"container/heap"
)

type LFU[K comparable, V any] struct {
	m        map[K]*entry[K, V]
	h        minHeap[K, V]
	ts       int // monotonically increasing counter
	capacity int
}

func New[K comparable, V any](capacity int) *LFU[K, V] {
	if capacity <= 0 {
		panic("lfu: capacity must be > 0")
	}
	h := make(minHeap[K, V], 0, capacity)
	heap.Init(&h)
	return &LFU[K, V]{
		m:        make(map[K]*entry[K, V]),
		h:        h,
		capacity: capacity,
	}
}

func (l *LFU[K, V]) Get(key K) (V, bool) {
	if e, ok := l.m[key]; ok {
		e.freq++
		e.ts = l.nextTs()
		heap.Fix(&l.h, e.index)
		return e.val, true
	}
	return *(new(V)), false
}

func (l *LFU[K, V]) Put(key K, val V) {
	if e, ok := l.m[key]; ok {
		e.val = val
		e.freq++
		e.ts = l.nextTs()
		heap.Fix(&l.h, e.index)
		return
	}
	if len(l.m) >= l.capacity {
		evicted := heap.Pop(&l.h).(*entry[K, V])
		delete(l.m, evicted.key)
	}
	e := &entry[K, V]{key: key, val: val, freq: 1, ts: l.nextTs()}
	l.m[key] = e
	heap.Push(&l.h, e)
}

func (l *LFU[K, V]) nextTs() int {
	l.ts++
	return l.ts
}

type entry[K comparable, V any] struct {
	key   K
	val   V
	freq  int
	ts    int // last access counter (for LRU tie-breaking)
	index int
}

type minHeap[K comparable, V any] []*entry[K, V]

func (h minHeap[K, V]) Len() int { return len(h) }

func (h minHeap[K, V]) Less(i, j int) bool {
	if h[i].freq == h[j].freq {
		return h[i].ts < h[j].ts
	}
	return h[i].freq < h[j].freq
}

func (h minHeap[K, V]) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *minHeap[K, V]) Push(x any) {
	e := x.(*entry[K, V])
	e.index = len(*h)
	*h = append(*h, e)
}

func (h *minHeap[K, V]) Pop() any {
	n := len(*h) - 1
	e := (*h)[n]
	*h = (*h)[:n]
	return e
}
