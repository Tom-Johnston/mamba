package comb

var maxSizes = []uint64{0, ^uint64(0), 4294967296, 33290221, 102570, 13467, 3612, 1449, 746, 453, 308, 227, 178, 147, 125, 110, 99, 90, 84, 80, 75, 72, 69, 68, 66, 65, 64, 63, 63, 62, 62, 62}

const largestK = 31
const maxInt = uint64(^uint(0) >> 1)

func addHasOverflowed(a, b int) (sum int, overflow bool) {
	sum = a + b
	if (sum^a)&(sum^b) < 0 {
		return sum, true
	}
	return sum, false
}

//CoeffUint64 returns the binomial coefficient n choose k as a uint64.
//CoeffUint64 panics if the calculation would overflow the uint64.
func CoeffUint64(n, k uint64) uint64 {
	if k > n {
		return 0
	}

	if k > n/2 {
		k = n - k
	}

	if k == 0 {
		return 1
	}

	if k > largestK || n > maxSizes[k] {
		panic("calculation overflows uint64")
	}

	var comb uint64 = 1
	var i uint64
	for i = 1; i <= k; i++ {
		comb *= (n - k + i)
		comb /= i
	}
	return comb
}

//Coeff calculates the binomial coefficient n choose k and returns it as an int.
//Coeff panics if the calculation would overflow.
func Coeff(n, k int) int {
	if n < 0 {
		panic("n must be non-negative")
	}
	if k < 0 {
		return 0
	}

	comb := CoeffUint64(uint64(n), uint64(k))

	if comb > maxInt {
		panic("coeff does not fit in an int")
	}
	return int(comb)
}

//Ranking and Unranking n choose k.

//Rank converts a sorted set of distinct ints to an int such that the set of ints can be recovered by Unrank.
//Rank panics if calculating the rank would cause an overflow.
func Rank(comb []int) int {
	rank := 0
	var overflow bool
	for i, v := range comb {
		c := Coeff(v, i+1)
		rank, overflow = addHasOverflowed(rank, c)
		if overflow {
			panic("rank has overflowed int")
		}
	}
	return rank
}

//Unrank converts an integer rank to a sorted set of k distinct ints. It is the inverse of rank.
func Unrank(rank, k int) []int {
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
