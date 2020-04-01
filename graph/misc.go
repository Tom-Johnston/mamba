package graph

func addHasOverflowed(a, b int) (sum int, overflow bool) {
	sum = a + b
	if (sum^a)&(sum^b) < 0 {
		return sum, true
	}
	return sum, false
}
