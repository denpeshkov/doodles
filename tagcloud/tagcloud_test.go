package tagcloud

import (
	"cmp"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type TagCloud[T cmp.Ordered] interface {
	AddTag(T)
	TopN(n int) []TagStat[T]
	Len() int
}

var testGen = [...]func() TagCloud[int]{
	func() TagCloud[int] { return NewMap[int]() },
	func() TagCloud[int] { return NewMapPQ[int]() },
	func() TagCloud[int] { return NewBST[int]() },
}

func TestEmpty(t *testing.T) {
	for _, tcGen := range testGen {
		tc := tcGen()
		t.Run(fmt.Sprintf("%T", tc), func(t *testing.T) {
			t.Parallel()

			topN := tc.TopN(1000)
			if len(topN) != 0 {
				t.Error("Expected empty tag cloud")
			}
		})
	}
}

func TestTopNGreaterThanCloudSize(t *testing.T) {
	for _, tcGen := range testGen {
		tc := tcGen()
		t.Run(fmt.Sprintf("%T", tc), func(t *testing.T) {
			t.Parallel()

			tc.AddTag(1)

			requestCount := 10
			topN := tc.TopN(requestCount)
			if len(topN) != 1 {
				t.Errorf("TopN(%d) returned %d items, want 1", requestCount, len(topN))
			}
		})
	}
}

func TestHappyPath(t *testing.T) {
	for _, tcGen := range testGen {
		tc := tcGen()
		t.Run(fmt.Sprintf("%T", tc), func(t *testing.T) {
			t.Parallel()

			tc.AddTag(1)
			tc.AddTag(2)
			tc.AddTag(2)

			top := tc.TopN(1)

			if len(top) != 1 {
				t.Errorf("TopN(1) returned %d items, want 1", len(top))
			}
			if top[0].Tag != 2 {
				t.Fatalf("TopN(1)[0].Tag == %d, want 2", top[0].Tag)
			}
			if top[0].Tag != 2 {
				t.Errorf("TopN(1)[0].Count == %d, want 2", top[0].Count)
			}
		})
	}
}

func TestTopN(t *testing.T) {
	for _, tcGen := range testGen {
		tc := tcGen()
		t.Run(fmt.Sprintf("%T", tc), func(t *testing.T) {
			t.Parallel()

			size := 1000
			for i := 0; i < size; i++ {
				for j := 0; j < i; j++ {
					tc.AddTag(i)
				}
			}

			validateTopN := func(n int) {
				topN := tc.TopN(n)
				if len(topN) != n {
					t.Errorf("TopN(%d) returned %d items, want %d", n, len(topN), n)
				}

				for i, el := range topN {
					tag := size - i - 1
					if tag != el.Tag {
						t.Errorf("TopN(%d)[%d].Tag = %d, want %d", n, i, el.Tag, tag)
					}
					if tag != el.Tag {
						t.Errorf("TopN(%d)[%d].Count = %d, want %d", n, i, el.Count, tag)
					}
				}
			}

			for i := 0; i < size; i++ {
				validateTopN(i)
			}
		})
	}
}

func TestTopNWithRepeatedOccurrence(t *testing.T) {
	for _, tcGen := range testGen {
		tc := tcGen()
		t.Run(fmt.Sprintf("%T", tc), func(t *testing.T) {
			t.Parallel()

			tc.AddTag(1)
			tc.AddTag(2)
			tc.AddTag(3)
			tc.AddTag(4)

			requestCount := 3
			topN := tc.TopN(requestCount)
			if len(topN) != requestCount {
				t.Errorf("TopN(%d) returned %d items, want %d", requestCount, len(topN), requestCount)
			}

			distinctMap := make(map[int]struct{})
			for i, v := range topN {
				if v.Count != 1 {
					t.Errorf("TopN(%d).Count = %d, want 1", i, v.Count)
				}
				distinctMap[v.Tag] = struct{}{}
			}

			if len(distinctMap) != requestCount {
				t.Errorf("TopN(%d) returned array with non-distinct tags", requestCount)
			}
		})
	}
}

func BenchmarkTagCloud(b *testing.B) {
	n := 10_000
	vv := make([]int, n)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range n {
		vv[i] = rnd.Intn(500)
	}

	for _, tcGen := range testGen {
		tc := tcGen()
		b.Run(fmt.Sprintf("%T", tc), func(b *testing.B) {
			for range b.N {
				for _, v := range vv {
					tc.AddTag(v)

					delta := tc.Len()/5 + 1
					for i := 1; i <= tc.Len()+delta; i += delta {
						tc.TopN(i)
					}
				}
			}
		})
	}
}
