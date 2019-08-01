package graph

//Binomial Coefficients

//BinomialCoeffSingle calculates the binomial coefficient n choose k.
//Binomial Coefficients can grow very quickly so this is only suitable for small n and k.
func BinomialCoeffSingle(n, k int) int {
	if k > n || k < 0 {
		return 0
	}
	comb := 1
	for i := 1; i <= k; i++ {
		comb *= (n - k + i)
		comb /= i
	}
	return comb
}

//Ranking and Unranking n choose k.

//RankCombination converts a sorted set of distinct ints to an int. This a bijection from the subsets of 0,..,n-1 with k elements to 0,...,nChoosek-1. It is the inverse of UnrankCombination.
func RankCombination(comb []int) int {
	rank := 0
	for i, v := range comb {
		rank += BinomialCoeffSingle(v, i+1)
	}
	return rank
}

//UnrankCombination converts an integer rank to a sorted set of k distinct ints. It is the inverse of RankCombination.
func UnrankCombination(rank, k int) []int {
	comb := make([]int, k)
	m := rank
	for i := k - 1; i >= 0; i-- {
		l := i + 1
		b := 1
		for b <= m {
			b *= (l + 1)
			b /= (l - i)
			l++
		}
		comb[i] = l - 1
		b *= (l - 1 - i)
		b /= l
		m -= b
	}
	return comb
}

//Permutations

//This

//PermutationIterator is a struct containing the state of the iterator.
type PermutationIterator struct {
	n int
	i int
	c []int
	p []int
}

//NewPermutation returns a new permuation iterator which iterates over all permutations of the elements in a.
//This doesn't modify a.
func NewPermutation(a []int) *PermutationIterator {
	copyOfA := make([]int, len(a))
	copy(copyOfA, a)
	p := PermutationIterator{i: -1, n: len(a), c: make([]int, len(a)), p: copyOfA}
	return &p
}

//Next returns the next permutation if there is one and nil if there are no further permuations.
func (p *PermutationIterator) Next() []int {
	if p.i == p.n {
		return nil
	}

	if p.i == -1 {
		p.i++
		return p.p
	}

	for p.i < p.n {
		if p.c[p.i] < p.i {
			if p.i%2 == 0 {
				p.p[0], p.p[p.i] = p.p[p.i], p.p[0]
			} else {
				p.p[p.c[p.i]], p.p[p.i] = p.p[p.i], p.p[p.c[p.i]]
			}
			p.c[p.i]++
			p.i = 0
			return p.p
		}
		p.c[p.i] = 0
		p.i++
	}
	return nil
}

//Helper functions on []int

func intsSum(a []int) (sum int) {
	for _, v := range a {
		sum += v
	}
	return sum
}

//IntsEqual returns true if a and b are the same and false otherwise.
func IntsEqual(a, b []int) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

//IntsHasPrefix returns true if s begins with prefix and false otherwise.
func IntsHasPrefix(s, prefix []int) bool {
	return len(s) >= len(prefix) && IntsEqual(s[0:len(prefix)], prefix)
}

//IntsCompare returns 1 if a is greater than b in lexicographic order, 0 if they are equal and -1 if b is greater than a.
func IntsCompare(a, b []int) int {
	for i := 0; ; i++ {
		if i >= len(a) && i >= len(b) {
			return 0
		}
		if i >= len(a) {
			return -1
		}
		if i >= len(b) {
			return 1
		}

		if a[i] > b[i] {
			return 1
		}
		if a[i] < b[i] {
			return -1
		}
	}
}

//IntsMax returns the largest int in a.
func IntsMax(a []int) int {
	max := a[0]
	for _, v := range a {
		if v > max {
			max = v
		}
	}
	return max
}

//IntsMin returns the smallest int in a.
func IntsMin(a []int) int {
	min := a[0]
	for _, v := range a {
		if v < min {
			min = v
		}
	}
	return min
}

//IntsAdd adds b to a.
func IntsAdd(a, b []int) {
	if len(a) != len(b) {
		panic("Cannot add two slices of different lengths.")
	}
	for i := range a {
		a[i] += b[i]
	}
}

//IntsReverse reverses a.
func IntsReverse(a []int) []int {
	for i := 0; i < len(a)/2; i++ {
		j := len(a) - i - 1
		a[i], a[j] = a[j], a[i]
	}
	return a
}
