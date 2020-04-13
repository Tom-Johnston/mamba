package itertools

import (
	"testing"

	"github.com/Tom-Johnston/mamba/ints"
)

func TestPartitions(t *testing.T) {
	//Truth data
	correctOutput := [][][]int{
		{{0, 1, 2, 3}},
		{{0, 1, 2}, {3}},
		{{0, 1, 3}, {2}},
		{{0, 1}, {2, 3}},
		{{0, 1}, {2}, {3}},
		{{0, 2, 3}, {1}},
		{{0, 2}, {1, 3}},
		{{0, 2}, {1}, {3}},
		{{0, 3}, {1, 2}},
		{{0}, {1, 2, 3}},
		{{0}, {1, 2}, {3}},
		{{0, 3}, {1}, {2}},
		{{0}, {1, 3}, {2}},
		{{0}, {1}, {2, 3}},
		{{0}, {1}, {2}, {3}},
	}

	pi := Partitions(4)
	index := 0
	for pi.Next() {
		if index >= len(correctOutput) {
			t.Log("Found too many partitions.")
			t.FailNow()
		}
		found := pi.Value()
		truth := correctOutput[index]
		if len(found) != len(truth) {
			t.Fail()
		}
		for i := range found {
			if !ints.Equal(found[i], truth[i]) {
				t.Error(found[i], truth[i])
			}
		}
		index++
	}
	if index != len(correctOutput) {
		t.Error("Found too few partitions.")
	}
}
