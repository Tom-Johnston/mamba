package comb

var maxSizes = []uint64{0, ^uint64(0), 4294967296, 33290221, 102570, 13467, 3612, 1449, 746, 453, 308, 227, 178, 147, 125, 110, 99, 90, 84, 80, 75, 72, 69, 68, 66, 65, 64, 63, 63, 62, 62, 62}
var smallEntries = [][]uint64{{1}, {1}, {1, 2}, {1, 3}, {1, 4, 6}, {1, 5, 10}, {1, 6, 15, 20}, {1, 7, 21, 35}, {1, 8, 28, 56, 70}, {1, 9, 36, 84, 126}, {1, 10, 45, 120, 210, 252}, {1, 11, 55, 165, 330, 462}, {1, 12, 66, 220, 495, 792, 924}, {1, 13, 78, 286, 715, 1287, 1716}, {1, 14, 91, 364, 1001, 2002, 3003, 3432}, {1, 15, 105, 455, 1365, 3003, 5005, 6435}, {1, 16, 120, 560, 1820, 4368, 8008, 11440, 12870}, {1, 17, 136, 680, 2380, 6188, 12376, 19448, 24310}, {1, 18, 153, 816, 3060, 8568, 18564, 31824, 43758, 48620}, {1, 19, 171, 969, 3876, 11628, 27132, 50388, 75582, 92378}, {1, 20, 190, 1140, 4845, 15504, 38760, 77520, 125970, 167960, 184756}, {1, 21, 210, 1330, 5985, 20349, 54264, 116280, 203490, 293930, 352716}, {1, 22, 231, 1540, 7315, 26334, 74613, 170544, 319770, 497420, 646646, 705432}, {1, 23, 253, 1771, 8855, 33649, 100947, 245157, 490314, 817190, 1144066, 1352078}, {1, 24, 276, 2024, 10626, 42504, 134596, 346104, 735471, 1307504, 1961256, 2496144, 2704156}, {1, 25, 300, 2300, 12650, 53130, 177100, 480700, 1081575, 2042975, 3268760, 4457400, 5200300}, {1, 26, 325, 2600, 14950, 65780, 230230, 657800, 1562275, 3124550, 5311735, 7726160, 9657700, 10400600}, {1, 27, 351, 2925, 17550, 80730, 296010, 888030, 2220075, 4686825, 8436285, 13037895, 17383860, 20058300}, {1, 28, 378, 3276, 20475, 98280, 376740, 1184040, 3108105, 6906900, 13123110, 21474180, 30421755, 37442160, 40116600}, {1, 29, 406, 3654, 23751, 118755, 475020, 1560780, 4292145, 10015005, 20030010, 34597290, 51895935, 67863915, 77558760}, {1, 30, 435, 4060, 27405, 142506, 593775, 2035800, 5852925, 14307150, 30045015, 54627300, 86493225, 119759850, 145422675, 155117520}, {1, 31, 465, 4495, 31465, 169911, 736281, 2629575, 7888725, 20160075, 44352165, 84672315, 141120525, 206253075, 265182525, 300540195}, {1, 32, 496, 4960, 35960, 201376, 906192, 3365856, 10518300, 28048800, 64512240, 129024480, 225792840, 347373600, 471435600, 565722720, 601080390}}

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

	if n <= 32 {
		return smallEntries[n][k]
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

//Coeffs calculates all binomial coefficeints m choose k for 0 <= m <= n and k <= m/2.
func Coeffs(n int) [][]int {
	coeffs := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		tmp := make([]int, i/2+1)
		tmp[0] = 1
		for j := 1; j < i/2+1; j++ {
			if 2*j == i {
				tmp[j] = 2 * coeffs[i-1][j-1]
				continue
			}
			tmp[j] = coeffs[i-1][j-1] + coeffs[i-1][j]
		}
		coeffs[i] = tmp
	}
	return coeffs
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
