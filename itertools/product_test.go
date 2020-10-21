package itertools

import (
	"testing"

	"github.com/Tom-Johnston/mamba/ints"
)

func TestProduct(t *testing.T) {
	correctOutputs := [][][]int{{{}},
		{},
		{
			{0, 0, 0},
			{0, 0, 1},
			{0, 1, 0},
			{0, 1, 1},
			{0, 2, 0},
			{0, 2, 1},
			{1, 0, 0},
			{1, 0, 1},
			{1, 1, 0},
			{1, 1, 1},
			{1, 2, 0},
			{1, 2, 1},
			{2, 0, 0},
			{2, 0, 1},
			{2, 1, 0},
			{2, 1, 1},
			{2, 2, 0},
			{2, 2, 1}},
	}
	inputs := [][]int{{}, {1, 0, 1}, {3, 3, 2}}

	for i := range inputs {
		iter := Product(inputs[i]...)
		correctOutput := correctOutputs[i]
		index := 0
		for iter.Next() {
			if index >= len(correctOutput) {
				t.Log("Found too many products.")
				t.FailNow()
			}
			found := iter.Value()
			truth := correctOutput[index]
			if len(found) != len(truth) {
				t.Fail()
			}
			for i := range found {
				if found[i] != truth[i] {
					t.Fail()
				}
			}
			index++
		}
		if index != len(correctOutput) {
			t.Error("Found too few products.")
		}
	}

}

func TestRestrictedPrefixProduct(t *testing.T) {
	correctOutputs := [][][]int{{{}},
		{},
		{},
		{
			{1, 2, 0},
			{1, 2, 1},
			{2, 0, 0},
			{2, 0, 1},
			{2, 1, 0},
			{2, 1, 1},
			{2, 2, 0},
			{2, 2, 1}},
	}
	inputs := [][]int{{}, {1, 0, 1}, {1, 1, 4}, {3, 3, 2}}
	inputFunctions := []func([]int) bool{func(a []int) bool { return true },
		func(a []int) bool { return true },
		func(a []int) bool {
			if ints.Equal(a, []int{0}) || ints.Equal(a, []int{1, 0}) || ints.Equal(a, []int{1, 1}) {
				return false
			}
			return true
		},
		func(a []int) bool {
			if ints.Equal(a, []int{0}) || ints.Equal(a, []int{1, 0}) || ints.Equal(a, []int{1, 1}) {
				return false
			}
			return true
		},
	}

	for i := range inputs {
		iter := RestrictedPrefixProduct(inputFunctions[i], inputs[i]...)
		correctOutput := correctOutputs[i]
		index := 0
		for iter.Next() {
			if index >= len(correctOutput) {
				t.Log("Found too many products.")
				t.FailNow()
			}
			if !ints.Equal(iter.Value(), correctOutput[index]) {
				t.Fail()
			}
			index++
		}
		if index != len(correctOutput) {
			t.Error("Found too few products.")
			t.Log(inputs[i])
		}
	}

}
