package ints

//Equal returns true if a and b are the same and false otherwise.
func Equal(a, b []int) bool {
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

//HasPrefix returns true if s begins with prefix and false otherwise.
func HasPrefix(s, prefix []int) bool {
	return len(s) >= len(prefix) && Equal(s[0:len(prefix)], prefix)
}

//Compare returns 1 if a is greater than b in lexicographic order, 0 if they are equal and -1 if b is greater than a.
func Compare(a, b []int) int {
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

//Max returns the largest int in a.
func Max(a []int) int {
	max := a[0]
	for _, v := range a {
		if v > max {
			max = v
		}
	}
	return max
}

//Min returns the smallest int in a.
func Min(a []int) int {
	min := a[0]
	for _, v := range a {
		if v < min {
			min = v
		}
	}
	return min
}

//Sum returns the sum of the ints in a.
func Sum(a []int) (sum int) {
	for _, v := range a {
		sum += v
	}
	return sum
}

//Add adds b to a. This modifies a.
func Add(a, b []int) {
	if len(a) != len(b) {
		panic("Cannot add two slices of different lengths.")
	}
	for i := range a {
		a[i] += b[i]
	}
}

//Reverse reverses a. This modifies a.
func Reverse(a []int) []int {
	for i := 0; i < len(a)/2; i++ {
		j := len(a) - i - 1
		a[i], a[j] = a[j], a[i]
	}
	return a
}
