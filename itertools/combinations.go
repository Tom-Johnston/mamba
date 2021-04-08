package itertools

//CombinationIterator iterates over all ways of choosing k elements from 0, ..., n-1 in lexicographic order.
type CombinationIterator struct {
	n    int
	k    int
	data []int
}

//Next moves the Iterator to the next combination, returning true if there is one and false is there isn't.
func (b *CombinationIterator) Next() bool {
	if b.k == 0 {
		b.k--
		return true
	}

	for i := b.k - 1; i >= 0; i-- {
		if b.data[i] < b.n+i-b.k {
			b.data[i]++
			for j := i + 1; j < b.k; j++ {
				b.data[j] = b.data[j-1] + 1
			}
			return true
		}
	}
	return false
}

//Value returns the current state of the iterator. This must not be modified.
func (b CombinationIterator) Value() []int {
	return b.data
}

//Combinations returns a new CombinationIterator which iterates over all subsets of k distinct elements from 0, ..., n-1 in lexicographic order.
func Combinations(n, k int) *CombinationIterator {
	data := make([]int, k)
	for i := 0; i < k; i++ {
		data[i] = i
	}

	if k > 0 {
		data[k-1]--
	}
	return &CombinationIterator{n: n, k: k, data: data}
}

//CombinationColexIterator iterates over all ways of chooding k elements from 0, ..., n-1 in colexicographic order.
//It can be initialised by calling CombinationsColex.
type CombinationColexIterator struct {
	n    int
	k    int
	j    int //Position to try to increase
	data []int
}

//CombinationsColex returns a new CombinationColexIterator which iterates over all subsets of k distinct elements from 0, ..., n-1 in colexicographic order.
func CombinationsColex(n, k int) *CombinationColexIterator {
	data := make([]int, k)
	for i := 0; i < k; i++ {
		data[i] = i
	}

	if k > 0 {
		data[k-1]--
	}
	return &CombinationColexIterator{n: n, k: k, j: k, data: data}
}

//Value returns the current subset. The output must not be modified.
func (b CombinationColexIterator) Value() []int {
	return b.data
}

//Next attempts to advance the iterator to the next subset, returning true if there is one and false if not.
func (b *CombinationColexIterator) Next() bool {
	if b.k <= 0 {
		b.k--
		return b.k == -1
	}

	if b.j >= b.k-1 {
		if b.data[b.k-1] == b.n-1 {
			return false
		}
		b.data[b.k-1]++
		b.j--
		return true
	}

	if b.j != -1 {
		b.data[b.j]++
		b.j--
		return true
	}

	for j := 0; j < b.k-1; j++ {
		if b.data[j] < b.data[j+1]-1 {
			b.data[j]++
			b.j = j - 1
			return true
		}
		b.data[j] = j
	}

	if b.data[b.k-1] == b.n-1 {
		return false
	}
	b.data[b.k-1]++
	b.j = b.k - 2
	return true
}

//MultisetCombinationIterator iterates over all multisets containing k elements and with a maximum of m[i] copies of i. Value returns the multiset of k elements and FreqValue returns a slice v where v[i] is the number of copies of i in the multiset.
type MultisetCombinationIterator struct {
	state []int
	m     []int
	k     int
	j     int

	//A buffer slice to return the value in as we iterate using FreqValue
	value []int
}

//MultisetCombinations returns an iterator which iterates over all multisets containing k elements and with a maximum of m[i] elements of type i. Value returns the multiset of k items and FreqValue returns a slice v where v[i] is the number of i in the multiset.
func MultisetCombinations(m []int, k int) *MultisetCombinationIterator {
	return &MultisetCombinationIterator{state: nil, m: m, k: k}
}

//Value returns the multiset of k elements.
//You may modify the return value.
func (iter MultisetCombinationIterator) Value() []int {
	c := 0

	for i, v := range iter.state {
		for j := 0; j < v; j++ {
			iter.value[c] = i
			c++
		}
	}

	return iter.value
}

//FreqValue returns a slice v where v[i] is the number of i in the multiset.
//You must not modify the return value.
func (iter MultisetCombinationIterator) FreqValue() []int {
	return iter.state
}

//Next attempts to advance the iterator to the next multiset, returning true if there is one and false if not.
//This is an implementation of Algorithm Q from The Art of Computer Programming Volume 4a section 7.2.1.3.
func (iter *MultisetCombinationIterator) Next() bool {
	if iter.state == nil {
		//Initial call
		iter.value = make([]int, iter.k)
		//Q2
		iter.state = make([]int, len(iter.m))
		x := iter.k
		for j := 0; j < len(iter.m); j++ {
			if x > iter.m[j] {
				iter.state[j] = iter.m[j]
				x -= iter.m[j]
				continue
			}
			iter.state[j] = x
			x = 0
			iter.j = j
			break
		}
		if x > 0 {
			return false
		}

		return true
	}

	//Q4
	x := 0
	j := iter.j
	if j == 0 {
		x = iter.state[0] - 1
		j = 1
	} else if iter.state[0] == 0 {
		x = iter.state[j] - 1
		iter.state[j] = 0
		j++
	} else {
		goto Q7
	}

	//Q5
Q5:
	if j >= len(iter.m) {
		return false
	}

	if iter.state[j] == iter.m[j] {
		x += iter.m[j]
		iter.state[j] = 0
		j++
		goto Q5
	}

	//Q6
	iter.state[j]++
	if x == 0 {
		iter.state[0] = 0
		iter.j = j
		return true
	}

	//Q2 again
	for j = 0; j < len(iter.m); j++ {
		if x > iter.m[j] {
			iter.state[j] = iter.m[j]
			x -= iter.m[j]
			continue
		}
		iter.state[j] = x
		x = 0
		break
	}
	iter.j = j
	return true

	//Q7
Q7:
	for iter.state[j] == iter.m[j] {
		j++
		if j >= len(iter.m) {
			return false
		}
	}

	iter.state[j]++
	j--
	iter.state[j]--
	if iter.state[0] == 0 {
		j = 1
	}
	iter.j = j
	return true
}
