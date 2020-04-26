package itertools

func partitionFromRestrictedGrowthString(rgs []int) [][]int {
	n := len(rgs)
	max := 0
	sizes := make([]int, n)
	for _, v := range rgs {
		sizes[v]++
		if v > max {
			max = v
		}
	}

	p := make([][]int, max+1)
	for i := range p {
		p[i] = make([]int, 0, sizes[i])
	}
	for i, v := range rgs {
		p[v] = append(p[v], i)
	}
	return p
}

//PartitionIterator iterates over all partitions of the set {0, ..., n-1}.
//The partitions are generated in lexicographic order of their restricted growth strings. It is safe to modify the output of .Value().
type PartitionIterator struct {
	n int
	m int
	a []int
	b []int
}

//Partitions returns a *PartitionIterator which iterates over all partitions of the set {0, ..., n-1}.
//The partitions are generated in lexicographic order of their restricted growth sequences.
//TODO Handle n = 0
func Partitions(n int) *PartitionIterator {
	if n < 1 {
		panic("Cannot handle n < 1")
	}
	a := make([]int, n)
	a[n-1] = -1
	b := make([]int, n)
	for i := range b {
		b[i] = 1
	}
	return &PartitionIterator{n: n, m: 1, a: a, b: b}
}

//Next tries to advance pi to the next partition, returning true if there is one and false if there isn't.
func (pi *PartitionIterator) Next() bool {
	if pi.a[pi.n-1] == pi.m {
		for j := pi.n - 2; j >= 1; j-- {
			if pi.a[j] != pi.b[j] {
				pi.a[j]++
				pi.m = pi.b[j]
				if pi.a[j] == pi.b[j] {
					pi.m++
				}
				for k := j + 1; k < pi.n-1; k++ {
					pi.a[k] = 0
					pi.b[k] = pi.m
				}
				pi.a[pi.n-1] = 0
				return true
			}
		}
		return false
	}
	pi.a[pi.n-1]++
	return true
}

//Value returns the current partition.
//It is safe to modify the output, although one shouldn't depend on this.
func (pi PartitionIterator) Value() [][]int {
	return partitionFromRestrictedGrowthString(pi.a)
}

//IntegerPartitionIterator iterates over all ways of writing n as the sum of positive integers in descending order.
//The integer partitons are generated in reverse lexicographic order.
type IntegerPartitionIterator struct {
	a []int
	m int
	q int //The last partition element that is greater than 1.
}

//IntegerPartitions returns an *IntegerPartitionIterator for iterating over all ways of writing n as the sum of positive integers in descending order.
//The integer partitons are generated in reverse lexicographic order.
func IntegerPartitions(n int) *IntegerPartitionIterator {
	if n == 0 {
		return &IntegerPartitionIterator{a: nil, m: 0, q: -1}
	}
	a := make([]int, n)
	for i := 1; i < n; i++ {
		a[i] = 1
	}
	a[0] = n
	return &IntegerPartitionIterator{a: a, m: 1, q: -2}
}

//Next tries to advance iter to the next integer partition, returning true if there is one and false if there isn't.
func (iter *IntegerPartitionIterator) Next() bool {
	if iter.q == -2 {
		if iter.a[0] == 1 {
			iter.q = -1
		} else {
			iter.q = 0
		}
		return true
	}
	if iter.q == -1 {
		return false
	}
	if iter.a[iter.q] == 2 {
		iter.a[iter.q] = 1
		iter.q--
		iter.m++
		return true
	}
	iter.a[iter.q]--
	x := iter.a[iter.q]
	tailSum := iter.m - iter.q
	for tailSum > x {
		iter.q++
		tailSum -= x
		iter.a[iter.q] = x
	}
	iter.m = iter.q + 1
	iter.a[iter.m] = tailSum
	iter.m++
	if tailSum > 1 {
		iter.q++
	}
	return true
}

//Value returns the current integer partition.
//It is not safe to modify the output
func (iter IntegerPartitionIterator) Value() []int {
	return iter.a[:iter.m]
}
