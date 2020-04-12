package sortints

import (
	"sort"
)

//SortedInts is used to hold a collection of sorted ints without repeats.
//This is generally used as a way of storing a set of ints.
type SortedInts []int

//NewSortedInts creates a SortedInts from the set input x.
func NewSortedInts(x ...int) SortedInts {
	tmp := make([]int, len(x))
	copy(tmp, x)
	sort.Ints(tmp)
	numberOfRepeats := 0
	for i := 1; i < len(tmp); i++ {
		if tmp[i-1] == tmp[i] {
			numberOfRepeats++
		} else {
			tmp[i-numberOfRepeats] = tmp[i]
		}
	}
	return tmp[:len(tmp)-numberOfRepeats]
}

//Range creates a SortedInts by adding all the elements of the form start + i*step which lie in [start, end).
func Range(start, end, step int) SortedInts {
	if (end < start && step > 0) || (end > start && step < 0) || (end != start && step == 0) {
		panic("Infinite set")
	}
	if end == start {
		return []int{}
	}

	if end < start {
		start, end = end, start
		step = -step
	}

	tmp := make([]int, 0, (end-start+step-1)/step)
	for i := start; i < end; i += step {
		tmp = append(tmp, i)
	}
	return tmp
}

//Remove modifies s by removing the element x if it is present.
func (s *SortedInts) Remove(x int) {
	index := sort.SearchInts(*s, x)
	if index < len(*s) && (*s)[index] == x {
		*s = (*s)[:index+copy((*s)[index:], (*s)[index+1:])]
	}
}

//Add modifies s by adding the element x... if they are not already present.
//The arguments do not need to be sorted.
func (s *SortedInts) Add(x ...int) {
	tmp := make([]int, len(x))
	copy(tmp, x)
	x = tmp
	sort.Ints(x)
	indices := make([]int, len(x)+1)
	numberAlreadySeen := 0
	for i := 0; i < len(x); i++ {
		index := sort.SearchInts(*s, x[i])
		if index < len(*s) && (*s)[index] == x[i] {
			numberAlreadySeen++
			indices[i] = -1
		} else {
			indices[i] = index
		}
	}
	//Check for duplicates
	for i := 0; i < len(x)-1; i++ {
		if x[i] == x[i+1] {
			indices[i+1] = -1
			numberAlreadySeen++
		}
	}
	indices[len(x)] = len(*s)
	tmp = make([]int, len(*s)+len(x)-numberAlreadySeen)
	numberNowSeen := 0
	for i := len(indices) - 2; i >= 0; i-- {
		if indices[i] != -1 {
			copy(tmp[indices[i]+len(x)-numberAlreadySeen-numberNowSeen:indices[i+1]+len(x)-numberAlreadySeen-numberNowSeen], (*s)[indices[i]:indices[i+1]])
			tmp[indices[i]+len(x)-numberAlreadySeen-numberNowSeen-1] = x[i]
			numberNowSeen++
		} else {
			indices[i] = indices[i+1]
		}
	}
	copy(tmp[:indices[0]], (*s)[:indices[0]])
	*s = tmp
}

//IntersectionSize returns the number of elements in both a and b.
func IntersectionSize(a, b SortedInts) int {
	intersection := 0
	i := 0 //Point in a
	j := 0 //Point in b
	for i < len(a) && j < len(b) {
		if a[i] == b[j] {
			intersection++
			i++
			j++
		} else if a[i] > b[j] {
			j++
		} else {
			i++
		}
	}
	return intersection
}

//Union returns are new SortedInts which is the union of a and b.
//a and b are not modified.
func Union(a, b SortedInts) SortedInts {
	r := make([]int, 0, len(a)+len(b)-IntersectionSize(a, b))
	i := 0 //Point in a
	j := 0 //Point in b
	for i < len(a) && j < len(b) {
		if a[i] == b[j] {
			r = append(r, a[i])
			i++
			j++
		} else if a[i] > b[j] {
			r = append(r, b[j])
			j++
		} else {
			r = append(r, a[i])
			i++
		}
	}
	if i < len(a) {
		r = append(r, a[i:]...)
	} else if j < len(b) {
		r = append(r, b[j:]...)
	}
	return r
}

//SetMinus returns a new SortedInts containing the elements in a but not b.
//a and b are not modified.
func SetMinus(a, b SortedInts) SortedInts {
	r := make([]int, 0, len(a)-IntersectionSize(a, b))
	i := 0 //Point in a
	j := 0 //Point in b
	for i < len(a) && j < len(b) {
		if a[i] == b[j] {
			i++
			j++
		} else if a[i] > b[j] {
			j++
		} else {
			r = append(r, a[i])
			i++
		}
	}
	r = append(r, a[i:]...)
	return r
}

//Intersection returns a new SortedInts containing the elements in a and b.
//a and b are not modified.
func Intersection(a, b SortedInts) SortedInts {
	r := make([]int, 0, IntersectionSize(a, b))
	i := 0 //Point in a
	j := 0 //Point in b
	for i < len(a) && j < len(b) {
		if a[i] == b[j] {
			r = append(r, a[i])
			i++
			j++
		} else if a[i] > b[j] {
			j++
		} else {
			i++
		}
	}
	return r
}

//XOR returns a new SortedInts containing the elements in a or b but not both.
//a and b are not modified.
func XOR(a, b SortedInts) SortedInts {
	xor := make([]int, 0, len(a)+len(b)-IntersectionSize(a, b))
	i := 0 //Point in a
	j := 0 //Point in b
	for i < len(a) && j < len(b) {
		if a[i] == b[j] {
			i++
			j++
		} else if a[i] > b[j] {
			xor = append(xor, b[j])
			j++
		} else {
			xor = append(xor, a[i])
			i++
		}
	}
	if i < len(a) {
		xor = append(xor, a[i:]...)
	} else if j < len(b) {
		xor = append(xor, b[j:]...)
	}
	return xor
}

//Complement returns a new SortedInts containing the elements in {0,..., n-1} but not a.
//a is not modified.
func Complement(n int, a SortedInts) SortedInts {
	b := make([]int, 0, n-len(a))
	aIndex := 0
	i := 0
	for i < n && aIndex < len(a) {
		if i > a[aIndex] {
			aIndex++
		} else if i == a[aIndex] {
			i++
			aIndex++
		} else {
			b = append(b, i)
			i++
		}
	}

	for ; i < n; i++ {
		b = append(b, i)
	}
	return b
}

//ContainsSingle returns if a contains x.
func ContainsSingle(a SortedInts, x int) bool {
	index := sort.SearchInts(a, x)
	return index < len(a) && a[index] == x
}

//ContainsSorted returns if a contains the SortedInts b.
func ContainsSorted(a, b SortedInts) bool {
	i := 0
	j := 0
	for i < len(a) && j < len(b) {
		if a[i] == b[j] {
			i++
			j++
		} else if a[i] > b[j] {
			return false
		} else {
			i++
		}
	}

	return j >= len(b)
}
