package main

import (
	"math/rand"
	"root/BPTree"
	"testing"
)

func BenchmarkInsertBPTree(b *testing.B) {
	b.StopTimer()

	tree := BPTree.New(99)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		tree.Insert(rand.Intn(100000), 5)
	}
}
