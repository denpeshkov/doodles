package invcount

import (
	"cmp"
)

// Count returns the number of inversions in the input slice.
// Time complexity is N*log(N)
func Count[T ~[]E, E cmp.Ordered](s T) int {
	aux := make(T, len(s))

	sc := append(T(nil), s...)
	return sortCount(sc, aux)
}

func sortCount[T ~[]E, E cmp.Ordered](s, a T) int {
	if len(s) <= 1 {
		return 0
	}

	m := len(s) / 2
	invsL := sortCount(s[:m], a)
	invsR := sortCount(s[m:], a)
	invs := invsL + invsR + mergeCount(s[:m], s[m:], a)

	copy(s, a)

	return invs
}

func mergeCount[T ~[]E, E cmp.Ordered](s1, s2, a T) int {
	invs := 0

	l1, l2 := len(s1), len(s2)
	for i, j, k := 0, 0, 0; k < l1+l2; k++ {
		switch {
		case i >= l1:
			a[k] = s2[j]
			j++
		case j >= l2:
			a[k] = s1[i]
			i++
		case s1[i] <= s2[j]:
			a[k] = s1[i]
			i++
		case s1[i] > s2[j]:
			a[k] = s2[j]
			j++
			// all to the right are inversions
			invs += l1 - i
		}
	}

	return invs
}
