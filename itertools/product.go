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

//RestrictedPrefixProductIterator iterates over all elements (a_0, a_1, ..., a_{len(n)-1}) of {0, ..., n[0] - 1} x {0, ..., n[1] - 1} x ... x {0, ..., n[len(n) - 1] - 1} which pass the tests f(a_0), f(a_0, a_1), ...
//If len(n) == 0, then the value []int{} is considered to pass.
type RestrictedPrefixProductIterator struct {
	state []int
	n     []int
	t     func([]int) bool
	empty bool
}

//RestrictedPrefixProduct returns a *RestrictedPrefixProductIterator which iterates over all elements (a_0, a_1, ..., a_{len(n)-1}) of {0, ..., n[0] - 1} x {0, ..., n[1] - 1} x ... x {0, ..., n[len(n) - 1] - 1} which pass the test f(a_0), f(a_0, a_1), ...
func RestrictedPrefixProduct(t func([]int) bool, n ...int) *RestrictedPrefixProductIterator {
	empty := false
	for _, v := range n {
		if v < 1 {
			empty = true
		}
	}
	m := len(n)
	state := make([]int, 0, m)
	//Create a deep copy of n in case it changes.
	tmpN := make([]int, m)
	copy(tmpN, n)
	return &RestrictedPrefixProductIterator{state: state, n: tmpN, t: t, empty: empty}
}

//Value returns the current value of the iterator. You must not modify the value.
func (p RestrictedPrefixProductIterator) Value() []int {
	return p.state
}

//Next attempts to move the iterator to the next state, returning true if there is one and false if there isn't.
func (p *RestrictedPrefixProductIterator) Next() bool {
	//This could be rewritten without goto statements if we wanted.
	m := len(p.n)

	if p.empty {
		return false
	}

	if len(p.state) == 0 {
		if m == 0 {
			p.empty = true
			return true
		}
		goto x1
	} else {
		goto x2
	}

x1:
	//Increase the size
	p.state = append(p.state, 0)
	goto x3

x2:
	//Find the next state
	if p.state[len(p.state)-1] < p.n[len(p.state)-1]-1 {
		p.state[len(p.state)-1]++
	} else {
		if len(p.state) == 1 {
			return false
		}
		p.state = p.state[:len(p.state)-1]
		goto x2
	}

x3:
	//Test the current state
	if !p.t(p.state) {
		//Fail - find the next state.
		goto x2
	}
	if len(p.state) < m {
		//We pass the tests but we need to extend the product
		goto x1
	}
	//We have a state!
	return true
}
