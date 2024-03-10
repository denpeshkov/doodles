package tagcloud

import (
	"cmp"

	"github.com/google/btree"
)

// TagCloudBST aggregates statistics about used tags.
type TagCloudBST[T cmp.Ordered] struct {
	tr *btree.BTreeG[*TagStat[T]]
}

// NewBST creates a new TagCloudBST instance.
func NewBST[T cmp.Ordered]() *TagCloudBST[T] {
	less := func(a, b *TagStat[T]) bool { return cmp.Less(a.Tag, b.Tag) }
	return &TagCloudBST[T]{
		tr: btree.NewG(2, less),
	}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count.
func (t *TagCloudBST[T]) AddTag(tag T) {
	ts := &TagStat[T]{Tag: tag, Count: 1}
	if tsOld, ok := t.tr.Delete(ts); ok {
		ts = tsOld
		ts.Count++
	}
	t.tr.ReplaceOrInsert(ts)
}

// TopN returns n most frequent tags ordered in descending order by occurrence count.
//   - If there are multiple tags with the same occurrence count then the order is undefined;
//   - If n is greater than the TagCloud size then all elements are returned.
func (t *TagCloudBST[T]) TopN(n int) []TagStat[T] {
	l := min(n, t.tr.Len())
	tss := make([]TagStat[T], l)
	i := 0
	iter := func(ts *TagStat[T]) bool {
		if i == l {
			return false
		}
		tss[i] = *ts
		i++
		return true
	}
	t.tr.Descend(iter)
	return tss
}

// Len returns the number of tags in the tags cloud.
func (t *TagCloudBST[T]) Len() int {
	return t.tr.Len()
}
