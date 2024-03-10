package tagcloud

import (
	"cmp"
	"slices"
)

// TagCloudMap aggregates statistics about used tags.
type TagCloudMap[T cmp.Ordered] struct {
	m map[T]int
}

// TagStat represents statistics regarding a single tag.
type TagStat[T cmp.Ordered] struct {
	Tag   T
	Count int
}

// NewMap creates a new TagCloudMap instance.
func NewMap[T cmp.Ordered]() *TagCloudMap[T] {
	return &TagCloudMap[T]{
		m: map[T]int{},
	}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count.
func (t *TagCloudMap[T]) AddTag(tag T) {
	t.m[tag]++
}

// TopN returns n most frequent tags ordered in descending order by occurrence count.
//   - If there are multiple tags with the same occurrence count then the order is undefined;
//   - If n is greater than the TagCloud size then all elements are returned.
func (t *TagCloudMap[T]) TopN(n int) []TagStat[T] {
	ts := make([]TagStat[T], len(t.m))
	i := 0
	for k, v := range t.m {
		ts[i] = TagStat[T]{k, v}
		i++
	}
	slices.SortFunc(ts, func(a, b TagStat[T]) int { return cmp.Compare(b.Count, a.Count) })
	return ts[:min(len(t.m), n)]
}

// Len returns the number of tags in the tags cloud.
func (t *TagCloudMap[T]) Len() int {
	return len(t.m)
}
