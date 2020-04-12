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
