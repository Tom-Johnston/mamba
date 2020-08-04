package ints

import (
	"fmt"
	"sort"
	"testing"
)

func BenchmarkIntsSort(b *testing.B) {
	sizes := []int{100, 1000}
	for _, n := range sizes {
		b.Run(fmt.Sprint("IntsSort-", n), func(b *testing.B) {
			b.StopTimer()
			unsorted := make([]int, n)
			for i := range unsorted {
				unsorted[i] = i ^ 0x2cc
			}
			data := make([]int, len(unsorted))
			for i := 0; i < b.N; i++ {
				copy(data, unsorted)
				b.StartTimer()
				Sort(data)
				b.StopTimer()
			}
		})
	}

}

func BenchmarkSortInts(b *testing.B) {
	sizes := []int{100, 1000}
	for _, n := range sizes {
		b.Run(fmt.Sprint("SortInts-", n), func(b *testing.B) {
			b.StopTimer()
			unsorted := make([]int, n)
			for i := range unsorted {
				unsorted[i] = i ^ 0x2cc
			}
			data := make([]int, len(unsorted))
			for i := 0; i < b.N; i++ {
				copy(data, unsorted)
				b.StartTimer()
				sort.Ints(data)
				b.StopTimer()
			}
		})
	}
}
