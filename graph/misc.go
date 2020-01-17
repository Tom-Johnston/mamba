package graph

func addHasOverflowed(a, b int) (sum int, overflow bool) {
	sum = a + b
	if (sum^a)&(sum^b) < 0 {
		return sum, true
	}
	return sum, false
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
