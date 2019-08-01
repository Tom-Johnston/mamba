package disjoint

import "fmt"

//Set is a simple implementation of the disjoint set structure as in https://en.wikipedia.org/wiki/Disjoint-set_data_structure.
//This is a efficient way of storing disjoint subsets of {0,...,n-1} when the only operations are checking if two elements are in same subet and taking the union of subsets.
type Set []int

//New returns a new disjoint where each element x is in the set {x}.
func New(n int) Set {
	ds := make([]int, n)
	for i := range ds {
		ds[i] = -1
	}
	return Set(ds)
}

//Find returns the current representation of the set which contains x. This only (maybe) changes when the set containing x is unioned with another set.
//This also flattens the tree.
func (dsPtr *Set) Find(x int) int {
	ds := *dsPtr
	if ds[x] < 0 {
		return x
	}
	currentPlace := x
	seenNumbers := []int{x}
	for true {
		if currentPlace = ds[currentPlace]; currentPlace < 0 {
			tmp := seenNumbers[len(seenNumbers)-1]
			for i := 0; i < len(seenNumbers)-2; i++ {
				ds[seenNumbers[i]] = tmp
			}
			return tmp
		}
		seenNumbers = append(seenNumbers, currentPlace)
	}
	return -1
}

//Union unions the sets containing x and y in ds.
//This implements union by rank.
func (dsPtr *Set) Union(x, y int) {
	ds := *dsPtr
	parentX := ds.Find(x)
	parentY := ds.Find(y)
	if parentX == parentY {
		return
	}
	if ds[parentX] < ds[parentY] {
		ds[parentY] = parentX
	} else if ds[parentY] < ds[parentX] {
		ds[parentX] = parentY
	} else {
		ds[parentX] = parentY
		ds[parentY]--
	}
}

//Sets returns the sets of ds. Each set is sorted from smallest to largest and the set of sets is ascending in the smallest element of the sets.
func (dsPtr *Set) Sets() [][]int {
	ds := *dsPtr
	sets := make([][]int, 0, 1)
outer:
	for i := range ds {
		for j := range sets {
			if ds.Find(i) == ds.Find(sets[j][0]) {
				sets[j] = append(sets[j], i)
				continue outer
			}
		}
		sets = append(sets, []int{i})
	}
	return sets
}

//String returns a human-readable string of ds.
func (dsPtr *Set) String() string {
	return fmt.Sprintf("%v", dsPtr.Sets())
}

//SmallestRep returns a []int where the the ith entry contains the smallest member of the set containing i.
func (dsPtr *Set) SmallestRep() []int {
	ds := *dsPtr
	sr := make([]int, len(ds))
position:
	for i := 0; i < len(ds); i++ {
		for j := 0; j < i; j++ {
			if ds.Find(i) == ds.Find(j) {
				sr[i] = sr[j]
				continue position
			}
		}
		sr[i] = i
	}
	return sr
}

//Roots returns a []int which contains the roots of the trees which make up the disjoint set. In particular, it returns one element from each disjoint set.
func (dsPtr *Set) Roots() []int {
	ds := *dsPtr
	roots := make([]int, 0, 1)
	for i, v := range ds {
		if v < 0 {
			roots = append(roots, i)
		}
	}
	return roots
}
