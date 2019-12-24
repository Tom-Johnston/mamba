package dawg

import (
	"bytes"
	"errors"
	"io"
	"math/bits"
	"sort"
)

//Dawg is a directed acyclic word graph (also known as a deterministic acyclic finite state automaton) is a data structure for storing a set of []byte in efficiently while still being easy to query.
type Dawg struct {
	id         uint64 //This is useful for operations which want to visit each node once and don't care about the path used to reach a node. Every node should have a unique ID but the code doesn't assume they are numbered 0, 1, 2... even though this currently happens in practice.
	numWords   int
	final      bool
	linkLabels []byte
	links      []*Dawg
}

//Builder is a structure to incrementally build a Dawg by repeatedly adding words (using Add) until all words have been added and then calling Finish().
//Words but be added in strictly increasing lexicographic order and duplicate words are not allowed. It is safe to add words directly to zero value without calling Initialise().
type Builder struct {
	d        *Dawg
	lastWord []byte
	register []*Dawg
	lastID   uint64
	done     bool
}

//Initialise sets up the internal state ready for use.
//This will be called automatically when calling Add() or Finish().
func (db *Builder) Initialise() {
	db.d = &Dawg{id: 0}
	db.lastID = 0
	db.register = []*Dawg{}
	db.lastWord = nil
	db.done = false
}

//Add adds b to the dawg.
//Words but be added in lexicographic order with no duplicates and words cannot be added to a builder that has already finished.
func (db *Builder) Add(b []byte) error {
	if db.d == nil {
		db.Initialise()
	}

	if db.done {
		return errors.New("DawgBuilder has already finished")
	}

	if db.lastWord != nil && bytes.Compare(db.lastWord, b) != -1 {
		return errors.New("byte slices must be added in lexicographical order")
	}
	db.lastWord = b
	_, suffix, lastNode := db.d.commonPrefix(b)
	if len(lastNode.links) != 0 {
		db.register = replaceOrRegister(lastNode, db.register)
	}
	db.lastID = lastNode.addSuffix(suffix, db.lastID)
	return nil
}

//Finish returns the finished dawg.
//No further modifications may be made to the dawg and finish can only be called once.
func (db *Builder) Finish() (*Dawg, error) {
	if db.d == nil {
		db.Initialise()
	}
	if db.done {
		return nil, errors.New("DawgBuilder has already finished")
	}
	replaceOrRegister(db.d, db.register)
	return db.d, nil
}

//New is a helper function which creates a dawg with the words in s.
//The words in s must be in lexicographic order with no duplicates.
func New(s [][]byte) (*Dawg, error) {
	db := new(Builder)
	db.Initialise()
	for _, b := range s {
		err := db.Add(b)
		if err != nil {
			return nil, err
		}
	}
	return db.Finish()
}

func (t *Dawg) numberOfNodes() int {
	nodes, _ := t.listNodesCountEdges()
	return len(nodes)
}

func (t *Dawg) listNodesCountEdges() (nodes []uint64, numEdges int) {
	nodes = append(nodes, t.id)
	numEdges = len(t.links)
	currDecisions := make([]int, 1, 1)
	currDecisions[0] = -1
	currDawgs := make([]*Dawg, 1, 1)
	currDawgs[0] = t
toCheckLoop:
	for true {
		currDawg := currDawgs[len(currDawgs)-1]
		for j := currDecisions[len(currDecisions)-1] + 1; j < len(currDawg.linkLabels); j++ {
			currDawgs = append(currDawgs, currDawg.links[j])
			currDecisions[len(currDecisions)-1] = j
			currDecisions = append(currDecisions, -1)
			linkID := currDawg.links[j].id
			index := sort.Search(len(nodes), func(i int) bool { return nodes[i] >= linkID })
			if index < len(nodes) && nodes[index] == linkID {
				//We've already seen the node so we don't need to visit it again.
				continue
			} else {
				nodes = append(nodes, 0)
				copy(nodes[index+1:], nodes[index:])
				nodes[index] = linkID
				numEdges += len(currDawg.links[j].links)
			}
			continue toCheckLoop
		}
		currDecisions = currDecisions[:len(currDecisions)-1]
		currDawgs = currDawgs[:len(currDawgs)-1]
		if len(currDawgs) == 0 {
			return
		}
	}
	//We can never return here.
	return nil, -1
}

type letterCount struct {
	letter byte
	count  int
}

//NumberOfWords returns the number of words stored in the dawg.
func (t *Dawg) NumberOfWords() int {
	return t.numWords
}

func (t *Dawg) commonPrefix(letters []byte) (prefix []byte, suffix []byte, dawg *Dawg) {
	dawg = t
	dawg.numWords++
	i := 0
	//TODO binary search here
letter:
	for ; i < len(letters); i++ {
		for j, link := range dawg.linkLabels {
			if link == letters[i] {
				dawg = dawg.links[j]
				dawg.numWords++
				continue letter
			}
		}
		break letter
	}
	prefix = make([]byte, i)
	copy(prefix, letters[:i])
	suffix = make([]byte, len(letters)-i)
	copy(suffix, letters[i:])
	return prefix, suffix, dawg
}

//Lookup returns the id of the node where the given word ends if it is present and a boolean representing indicating if the word is in the dawg.
func (t *Dawg) Lookup(word []byte) (int, bool) {
	dawg := t
	index := -1 //To find the index we need to count the number of words we see up to and including this word.
	if dawg.final {
		index++
	}
	//TODO Maybe binary search here
letter:
	for _, l := range word {
		for j, link := range dawg.linkLabels {
			if link == l {
				dawg = dawg.links[j]
				if dawg.final {
					index++
				}
				continue letter
			} else {
				index += dawg.links[j].numWords
			}
		}
		return 0, false
	}
	if dawg.final {
		return index, true
	}
	return 0, false
}

func replaceOrRegister(t *Dawg, register []*Dawg) []*Dawg {
	lastChild := t.links[len(t.links)-1]
	if len(lastChild.links) != 0 {
		register = replaceOrRegister(lastChild, register)
	}
	for _, u := range register {
		if areEquivalent(lastChild, u) {
			t.links[len(t.links)-1] = u
			return register
		}
	}
	register = append(register, lastChild)
	return register
}

func areEquivalent(t, u *Dawg) bool {
	if t.final != u.final {
		return false
	}

	if len(t.links) != len(u.links) {
		return false
	}

	for i := 0; i < len(t.links); i++ {
		if t.linkLabels[i] != u.linkLabels[i] {
			return false
		}
	}

	for i := 0; i < len(t.links); i++ {
		if t.links[i] != u.links[i] {
			return false
		}
	}
	return true
}

func (t *Dawg) addSuffix(suffix []byte, lastID uint64) uint64 {
	currDawg := t
	for _, b := range suffix {
		lastID++
		currDawg.linkLabels = append(currDawg.linkLabels, b)
		tmpLinks := make([]*Dawg, 0)
		tmpLinkLabels := make([]byte, 0)
		tmpDawg := &Dawg{id: lastID, numWords: 1, final: false, linkLabels: tmpLinkLabels, links: tmpLinks}
		currDawg.links = append(currDawg.links, tmpDawg)
		currDawg = tmpDawg
	}
	currDawg.final = true
	return lastID
}

//GobEncode encodes the dawg as a []byte suitable for sending or storage. It preserves the id of the nodes in case they are relevant for other data.
func (t *Dawg) GobEncode() ([]byte, error) {

	sortedNodes, numEdges := t.listNodesCountEdges()
	numNodes := len(sortedNodes)

	nodes := make([]uint64, numNodes)

	convertID := func(id uint64) uint64 {
		return uint64(sort.Search(numNodes, func(i int) bool { return sortedNodes[i] >= id }))
	}

	//Need to recalculate the maximum size.
	b := make([]byte, 0, 8+9*numNodes+9*numEdges)

	buf := make([]byte, 9) //This is the buffer for encoding uints.

	//Encode the numberOfNodes first.
	buf = encodeUint64(uint64(numNodes), buf)
	b = append(b, buf...)
	for _, id := range sortedNodes {
		buf = encodeUint64(id, buf)
		b = append(b, buf...)
	}

	currDecisions := make([]int, 1, 1)
	currDecisions[0] = -1
	currDawgs := make([]*Dawg, 1, 1)
	currDawgs[0] = t

	buf = encodeUint64(convertID(t.id), buf)
	b = append(b, buf...)

	buf = encodeUint64(uint64(t.numWords), buf)
	b = append(b, buf...)

	if t.final {
		b = append(b, 1)
	} else {
		b = append(b, 0)
	}

	b = append(b, byte(len(t.linkLabels)))
	for i := range t.linkLabels {
		b = append(b, t.linkLabels[i])
		buf = encodeUint64(convertID(t.links[i].id), buf)
		b = append(b, buf...)
	}

toCheckLoop:
	for true {
		currDawg := currDawgs[len(currDawgs)-1]

		for j := currDecisions[len(currDecisions)-1] + 1; j < len(currDawg.linkLabels); j++ {
			currDawgs = append(currDawgs, currDawg.links[j])
			currDecisions[len(currDecisions)-1] = j
			currDecisions = append(currDecisions, -1)
			linkDawg := currDawg.links[j]
			linkID := linkDawg.id
			index := sort.Search(len(nodes), func(i int) bool { return nodes[i] >= linkID })
			if index < len(nodes) && nodes[index] == linkID {
				//We've already seen the node so we don't need to visit it again.
				continue
			} else {
				nodes = append(nodes, 0)
				copy(nodes[index+1:], nodes[index:])
				nodes[index] = linkID
				numEdges += len(linkDawg.links)

				buf = encodeUint64(convertID(linkID), buf)
				b = append(b, buf...)

				buf = encodeUint64(uint64(linkDawg.numWords), buf)
				b = append(b, buf...)

				if linkDawg.final {
					b = append(b, 1)
				} else {
					b = append(b, 0)
				}

				b = append(b, byte(len(linkDawg.linkLabels)))
				for i := range linkDawg.linkLabels {
					b = append(b, linkDawg.linkLabels[i])
					buf = encodeUint64(convertID(linkDawg.links[i].id), buf)
					b = append(b, buf...)
				}
			}
			continue toCheckLoop
		}
		currDecisions = currDecisions[:len(currDecisions)-1]
		currDawgs = currDawgs[:len(currDawgs)-1]
		if len(currDawgs) == 0 {
			return b, nil
		}
	}
	//We can never return here.
	return nil, nil
}

//GobDecode decodes the dawg given in b into t, replacing the current contents of t.
func (t *Dawg) GobDecode(b []byte) error {
	buf := make([]byte, 9) //This is the buffer for decoding uints.
	r := bytes.NewReader(b)
	numNodes, _, err := decodeUint64(r, buf)
	if err != nil {
		return err
	}
	ts := make([]*Dawg, numNodes)
	ts[0] = t
	var i uint64
	for i = 1; i < numNodes; i++ {
		ts[i] = new(Dawg)
	}

	for i = 0; i < numNodes; i++ {
		x, _, err := decodeUint64(r, buf)
		if err != nil {
			return err
		}
		ts[i].id = x
	}

	for i = 0; i < numNodes; i++ {
		indexID, _, err := decodeUint64(r, buf)
		if err != nil {
			return err
		}

		numWords, _, err := decodeUint64(r, buf)
		if err != nil {
			return err
		}

		ts[indexID].numWords = int(numWords)

		final, err := r.ReadByte()
		if err != nil {
			return err
		}

		if final != 0 {
			ts[indexID].final = true
		} else {
			ts[indexID].final = false
		}

		numChild, _, err := decodeUint64(r, buf)
		ts[indexID].linkLabels = make([]byte, 0, numChild)
		ts[indexID].links = make([]*Dawg, 0, numChild)
		var j uint64
		for j = 0; j < numChild; j++ {
			label, err := r.ReadByte()
			if err != nil {
				return err
			}
			target, _, err := decodeUint64(r, buf)
			ts[indexID].linkLabels = append(ts[indexID].linkLabels, label)
			ts[indexID].links = append(ts[indexID].links, ts[target])
		}
	}
	return nil
}

//This encodes a uint64 in a similar but not identical fashion to Gob.
//If the value of x <= 127, we encode it directly in the first byte. Else the first byte contains 128 + bytelen(x) where bytelen(x) is the number of bytes used to encode x, and the remaining bytes encode x in big endian order. This assumes buf contains at least 9 bytes. This is probably less efficient than the Gob version.
func encodeUint64(x uint64, buf []byte) []byte {
	buf = buf[:1]
	if x <= 127 {
		buf[0] = uint8(x)
		return buf
	}
	zeroBytes := (bits.LeadingZeros64(x) >> 3) //We have a zero byte for every 8 bits at the start which are all 0.
	buf = buf[:9-zeroBytes]
	buf[0] = 128 + 8 - byte(zeroBytes)
	for i := 0; i < 8-zeroBytes; i++ {
		buf[1+i] = byte(x >> uint(8*(7-(i+zeroBytes))))
	}
	return buf
}

//This assumes buf has at least 9 bytes.
func decodeUint64(r io.Reader, buf []byte) (x uint64, width int, err error) {
	width = 1
	n, err := io.ReadFull(r, buf[0:width])
	if n == 0 {
		return
	}
	b := buf[0]
	if buf[0] <= 127 {
		return uint64(b), width, nil
	}
	n = int(b) - 128
	if n > 8 {
		err = errors.New("decoding: too many bytes used for a uint64")
		return
	}
	width, err = io.ReadFull(r, buf[0:n])
	width++
	if err != nil {
		return
	}
	for _, b := range buf[0 : width-1] {
		x = x<<8 | uint64(b)
	}
	return
}
