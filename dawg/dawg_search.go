package dawg

import "sort"

//Searcher is the interface for searching a dawg for words which match a condition.
type Searcher interface {
	AllowStep(b byte) bool
	Step(b byte)
	Backstep()

	AllowWord() bool
}

//Search is used to find words which are accepted by all searchers simultaneously.
func (t *Dawg) Search(searchers ...Searcher) [][]byte {
	solns := make([][]byte, 0)
	currDecisions := make([]int, 1)
	currDecisions[0] = -1
	currDawgs := make([]*Dawg, 1)
	currDawgs[0] = t
	currWord := make([]byte, 0)
toCheckLoop:
	for true {
		//Do we allow this word?
		allow := true
		for _, srch := range searchers {
			if srch.AllowWord() != true {
				allow = false
				break
			}
		}
		if allow {
			tmpWord := make([]byte, len(currWord))
			copy(tmpWord, currWord)
			solns = append(solns, tmpWord)
		}
		//Let's check if we can move from here.
		currDawg := currDawgs[len(currDawgs)-1]
		for j := currDecisions[len(currDecisions)-1] + 1; j < len(currDawg.linkLabels); j++ {
			l := currDawg.linkLabels[j]
			allowStep := true
			for i := range searchers {
				if searchers[i].AllowStep(l) != true {
					allowStep = false
					break
				}
			}
			if !allowStep {
				continue
			}
			for i := range searchers {
				searchers[i].Step(l)
			}
			currWord = append(currWord, l)
			currDawgs = append(currDawgs, currDawg.links[j])
			currDecisions[len(currDecisions)-1] = j
			currDecisions = append(currDecisions, -1)
			continue toCheckLoop
		}
		//Now we go back
		//Can we go back?
		if len(currWord) == 0 {
			return solns
		}
		currWord = currWord[:len(currWord)-1]
		currDecisions = currDecisions[:len(currDecisions)-1]
		currDawgs = currDawgs[:len(currDawgs)-1]
		//Tell the searchers we are stepping back.
		for i := range searchers {
			searchers[i].Backstep()
		}
	}
	//We can never return here
	return solns
}

//PatternSearcher searches the dawg for words where each letter matches the letter in the pattern except in the positions of  the pattern which contain a blank.
type PatternSearcher struct {
	pattern []byte
	blank   byte
	index   int
}

//AllowStep checks if taking the step b is valid from the current position.
func (p PatternSearcher) AllowStep(b byte) bool {
	if len(p.pattern) <= p.index {
		return false
	}
	if p.pattern[p.index] == p.blank || p.pattern[p.index] == b {
		return true
	}
	return false
}

//Step takes the step b from the current position. This modifies the searcher.
func (p *PatternSearcher) Step(b byte) {
	p.index++
}

//Backstep undoes the last step from the searcher.This modifies the searcher.
func (p *PatternSearcher) Backstep() {
	p.index--
}

//AllowWord checks if the current position is allowed as a matching word.
func (p PatternSearcher) AllowWord() bool {
	if p.index == len(p.pattern) {
		return true
	}
	return false
}

//NewPatternSearcher returns a pattern searcher for searching with the given pattern and blank.
func NewPatternSearcher(pattern []byte, blank byte) *PatternSearcher {
	return &PatternSearcher{pattern: pattern, blank: blank, index: 0}
}

//AnagramSearcher searches a dawg for words which have the same multiset of letters as the target. Any letters with the value blank are assumed to be wildcards and match any letter.
type AnagramSearcher struct {
	counts       []letterCount
	blanks       int
	blank        byte
	targetLength int
	currPath     []byte
}

//AllowStep checks if taking the step b is valid from the current position.
func (p AnagramSearcher) AllowStep(b byte) bool {
	if p.targetLength <= len(p.currPath) {
		return false
	}
	if p.blanks > 0 {
		return true
	}
	for i := range p.counts {
		if p.counts[i].letter == b && p.counts[i].count > 0 {
			return true
		}
	}
	return false
}

//Step takes the step b from the current position. This modifies the searcher.
func (p *AnagramSearcher) Step(b byte) {
	for i := range p.counts {
		if p.counts[i].letter == b && p.counts[i].count > 0 {
			p.counts[i].count--
			p.currPath = append(p.currPath, b)
			return
		}
	}
	p.blanks--
	p.currPath = append(p.currPath, p.blank)
}

//Backstep undoes the last step from the searcher.This modifies the searcher.
func (p *AnagramSearcher) Backstep() {
	pathElem := p.currPath[len(p.currPath)-1]
	p.currPath = p.currPath[:len(p.currPath)-1]
	if pathElem == p.blank {
		p.blanks++
		return
	}
	for i := range p.counts {
		if p.counts[i].letter == pathElem {
			p.counts[i].count++
			return
		}
	}
}

//AllowWord checks if the current position is allowed as a matching word.
func (p AnagramSearcher) AllowWord() bool {
	if p.targetLength == len(p.currPath) {
		return true
	}
	return false
}

//NewAnagramSearcher returns an anagram searcher to find anagrams of a word in a dawg with blanks.
func NewAnagramSearcher(anagram []byte, blank byte) *AnagramSearcher {
	tmp := make([]byte, len(anagram))
	copy(tmp, anagram)
	sort.Slice(tmp, func(i, j int) bool { return anagram[i] < anagram[j] })
	counts := make([]letterCount, 0)
	blanks := 0
	for i, l := range tmp {
		if l == blank {
			blanks++
			continue
		}
		if i > 1 && tmp[i-1] == tmp[i] {
			counts[len(counts)-1].count++
		} else {
			counts = append(counts, letterCount{letter: l, count: 1})
		}
	}

	currPath := make([]byte, 0, len(anagram))

	return &AnagramSearcher{counts: counts, blanks: blanks, blank: blank, targetLength: len(anagram), currPath: currPath}
}
