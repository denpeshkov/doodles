package tagcloud

import (
	"cmp"
	"container/heap"
)

// TagCloudMapPQ aggregates statistics about used tags.
type TagCloudMapPQ[T cmp.Ordered] struct {
	mi map[T]int
	pq []*TagStat[T]
}

// NewMap creates a new TagCloudMapPq instance.
func NewMapPQ[T cmp.Ordered]() *TagCloudMapPQ[T] {
	tc := &TagCloudMapPQ[T]{
		mi: map[T]int{},
		pq: []*TagStat[T]{},
	}
	heap.Init(tc)
	return tc
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count.
func (t *TagCloudMapPQ[T]) AddTag(tag T) {
	if i, ok := t.mi[tag]; !ok {
		heap.Push(t, &TagStat[T]{tag, 1})
	} else {
		count := t.pq[i].Count
		t.update(tag, count+1)
	}
}

// TopN returns n most frequent tags ordered in descending order by occurrence count.
//   - If there are multiple tags with the same occurrence count then the order is undefined;
//   - If n is greater than the TagCloud size then all elements are returned.
func (t *TagCloudMapPQ[T]) TopN(n int) []TagStat[T] {
	l := t.Len()
	n = min(l, n)
	tss := make([]TagStat[T], n)
	for i := range n {
		tss[i] = *heap.Pop(t).(*TagStat[T])
	}
	// restore heap
	t.pq = t.pq[:l]
	heap.Init(t)
	return tss
}

// Len returns the number of tags in the tags cloud.
func (t *TagCloudMapPQ[T]) Len() int {
	return len(t.pq)
}

func (t *TagCloudMapPQ[T]) Less(i, j int) bool {
	return t.pq[j].Count < t.pq[i].Count
}

func (tc *TagCloudMapPQ[T]) Swap(i, j int) {
	pq, mi := tc.pq, tc.mi
	pq[i], pq[j] = pq[j], pq[i]
	mi[pq[i].Tag] = i
	mi[pq[j].Tag] = j
}

func (t *TagCloudMapPQ[T]) Push(x any) {
	ts := x.(*TagStat[T])
	t.mi[ts.Tag] = len(t.pq)
	t.pq = append(t.pq, ts)
}

func (t *TagCloudMapPQ[T]) Pop() any {
	n := len(t.pq)
	tag := t.pq[n-1]
	t.pq = t.pq[:n-1]
	return tag
}

func (t *TagCloudMapPQ[T]) update(tag T, count int) {
	i := t.mi[tag]
	t.pq[i].Count = count
	heap.Fix(t, i)
}
