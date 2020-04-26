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

func TestIntegerPartitions(t *testing.T) {
	correctOutput := [][]int{
		{8},
		{7, 1},
		{6, 2},
		{6, 1, 1},
		{5, 3},
		{5, 2, 1},
		{5, 1, 1, 1},
		{4, 4},
		{4, 3, 1},
		{4, 2, 2},
		{4, 2, 1, 1},
		{4, 1, 1, 1, 1},
		{3, 3, 2},
		{3, 3, 1, 1},
		{3, 2, 2, 1},
		{3, 2, 1, 1, 1},
		{3, 1, 1, 1, 1, 1},
		{2, 2, 2, 2},
		{2, 2, 2, 1, 1},
		{2, 2, 1, 1, 1, 1},
		{2, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1},
	}
	iter := IntegerPartitions(8)
	index := 0
	for iter.Next() {
		if index >= len(correctOutput) {
			t.Log("Found too many integer partitions.")
			t.FailNow()
		}
		found := iter.Value()
		truth := correctOutput[index]
		if !ints.Equal(found, truth) {
			t.Fail()
		}
		index++
	}
	if index != len(correctOutput) {
		t.Error("Found too few integer partitions.")
	}
}
