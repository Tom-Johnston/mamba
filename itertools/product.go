package itertools

//ProductIterator iterates over {0, ..., n[0] - 1} x {0, ..., n[1] - 1} x ... x {0, ..., n[len(n) - 1] - 1}.
//It should be initliased using Product.
type ProductIterator struct {
	state []int
	n     []int
	empty bool
}

//Product returns a *ProductIterator to iterate over {0, ..., n[0] - 1} x {0, ..., n[1] - 1} x ... x {0, ..., n[len(n) - 1] - 1}.
func Product(n ...int) *ProductIterator {
	empty := false
	for _, v := range n {
		if v < 1 {
			empty = true
		}
	}

	m := len(n)
	//Initialise the state to be (0, 0, ..., 0, -1). The first call to next will increase the last coordinate and the first value will be (0, 0, ..., 0, 0).
	state := make([]int, m)
	if m > 0 {
		state[m-1] = -1
	}
	//Create a deep copy of n in case it changes.
	tmpN := make([]int, m)
	copy(tmpN, n)
	return &ProductIterator{state: state, n: tmpN, empty: empty}
}

//Value returns the current state of the iterator. It isn't safe to modify the state.
func (p ProductIterator) Value() []int {
	return p.state
}

//Next attempts to move the iterator to the next state, returning true if there is one and false if there isn't.
func (p *ProductIterator) Next() bool {
	n := len(p.state)
	for j := n - 1; j >= 0; j-- {
		if p.state[j] < p.n[j]-1 {
			p.state[j]++
			for k := j + 1; k < n; k++ {
				p.state[k] = 0
			}
			return !p.empty
		}
	}
	if n == 0 && !p.empty {
		p.empty = true
		return true
	}
	return false
}
