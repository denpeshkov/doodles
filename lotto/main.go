package main

import (
	"fmt"
)

func main() {
	ts := []set{
		new([]int{1, 2, 3}),
		new([]int{1, 4, 5}),
		new([]int{2, 4, 5}),
		new([]int{3, 4, 5}),
	}

	fmt.Println(setCover(ts, 5, 2))
}

// setCover returns all sets not covered by tickets.
func setCover(ts []set, n, l int) []set {
	ss := genKSub(n, l)
	for _, t := range ts {
		i := 0
		for _, s := range ss {
			if !t.containsAll(s) {
				ss[i] = s
				i++
			}
		}
		ss = ss[:i]
	}
	return ss
}

// set represent a set of elements from universum {1..64}.
// set is implemented as a bit-vector.
type set uint64

// new returns a set containing elements ee.
func new(ee []int) set {
	var s set
	for _, e := range ee {
		s.add(e)
	}
	return s
}

// add adds element e to the set.
func (v *set) add(e int) {
	(*v) |= 1 << (64 - e)
}

// delete removes element e from the set.
func (v *set) delete(e int) {
	(*v) &^= 1 << (64 - e)
}

// contains determines whether the set contains an element e.
func (v set) contains(e int) bool {
	return v&(1<<(64-e)) != 0
}

// containsAll determines whether set s is a subset of the current set.
func (v set) containsAll(s set) bool {
	return ^(v | ^s) == 0
}

// elements returns elements of the set
func (v set) elements() []int {
	var s []int
	for i := 1; i <= 64; i++ {
		if v.contains(i) {
			s = append(s, i)
		}
	}
	return s
}

// String returns elements of the set as a string.
func (v set) String() string {
	return fmt.Sprint(v.elements())
}

// Function genKSub generates all subsets of length l from the subset {1..n}.
func genKSub(n, l int) []set {
	var res []set

	genKSubAux(0, n, l, 1, 0, &res)

	return res
}

func genKSubAux(tmp set, n, l, e, depth int, res *[]set) {
	if depth == l {
		*res = append(*res, tmp)
		return
	}

	for i := e; i <= n-l+depth+1; i++ {
		tmp.add(i)
		genKSubAux(tmp, n, l, i+1, depth+1, res)
		tmp.delete(i)
	}
}
