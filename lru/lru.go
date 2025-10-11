package lru

import "container/list"

type kv[K comparable, V any] struct {
	key K
	val V
}

type LRU[K comparable, V any] struct {
	m        map[K]*list.Element // Element's Value is of type kv
	l        *list.List          // List of kv
	capacity int
}

func New[K comparable, V any](capacity int) *LRU[K, V] {
	if capacity <= 0 {
		panic("lru: capacity must be positive")
	}
	return &LRU[K, V]{
		m:        make(map[K]*list.Element, capacity),
		l:        list.New(),
		capacity: capacity,
	}
}

func (l *LRU[K, V]) Get(key K) (V, bool) {
	if e, ok := l.m[key]; ok {
		l.l.MoveToFront(e)
		return e.Value.(kv[K, V]).val, true
	}
	return *(new(V)), false
}

func (l *LRU[K, V]) Put(key K, val V) {
	if e, ok := l.m[key]; ok {
		e.Value = kv[K, V]{key, val}
		l.l.MoveToFront(e)
		return
	}
	if len(l.m) == l.capacity {
		v := l.l.Remove(l.l.Back())
		delete(l.m, v.(kv[K, V]).key)
	}
	e := l.l.PushFront(kv[K, V]{key, val})
	l.m[key] = e
}
